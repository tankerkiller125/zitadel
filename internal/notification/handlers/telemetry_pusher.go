package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/query"

	"github.com/zitadel/zitadel/internal/api/call"

	"github.com/zitadel/zitadel/internal/repository/pseudo"

	"github.com/zitadel/zitadel/internal/errors"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query/projection"
)

const (
	TelemetryProjectionTable = "projections.telemetry"
)

type TelemetryPusherConfig struct {
	Enabled   bool
	Endpoints []string
}

type telemetryPusher struct {
	crdb.StatementHandler
	commands                       *command.Commands
	queries                        *NotificationQueries
	metricSuccessfulDeliveriesJSON string
	metricFailedDeliveriesJSON     string
	endpoints                      []string
}

func NewTelemetryPusher(
	ctx context.Context,
	telemetryCfg TelemetryPusherConfig,
	handlerCfg crdb.StatementHandlerConfig,
	commands *command.Commands,
	queries *NotificationQueries,
	metricSuccessfulDeliveriesJSON,
	metricFailedDeliveriesJSON string,
) *telemetryPusher {
	p := new(telemetryPusher)
	handlerCfg.ProjectionName = TelemetryProjectionTable
	handlerCfg.Reducers = []handler.AggregateReducer{{}}
	if telemetryCfg.Enabled {
		handlerCfg.Reducers = p.reducers()
	}
	p.endpoints = telemetryCfg.Endpoints
	p.StatementHandler = crdb.NewStatementHandler(ctx, handlerCfg)
	p.commands = commands
	p.queries = queries
	p.metricSuccessfulDeliveriesJSON = metricSuccessfulDeliveriesJSON
	p.metricFailedDeliveriesJSON = metricFailedDeliveriesJSON
	projection.TelemetryPusherProjection = p
	return p
}

func (t *telemetryPusher) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{{
		Aggregate: pseudo.AggregateType,
		EventRedusers: []handler.EventReducer{{
			Event:  pseudo.TimestampEventType,
			Reduce: t.pushMilestones,
		}},
	}}
}

func (t *telemetryPusher) pushMilestones(event eventstore.Event) (*handler.Statement, error) {
	ctx := call.WithTimestamp(context.Background())
	timestampEvent, ok := event.(pseudo.TimestampEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-lDTs5", "reduce.wrong.event.type %s", event.Type())
	}

	isReached, err := query.NewNotNullQuery(query.MilestoneReachedDateColID)
	if err != nil {
		return nil, err
	}
	isNotPushed, err := query.NewIsNullQuery(query.MilestonePushedDateColID)
	if err != nil {
		return nil, err
	}
	hasPrimaryDomain, err := query.NewNotNullQuery(query.MilestonePrimaryDomainColID)
	if err != nil {
		return nil, err
	}
	unpushedMilestones, err := t.queries.Queries.SearchMilestones(ctx, timestampEvent.InstanceIDs, &query.MilestonesSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        100,
			SortingColumn: query.MilestoneReachedDateColID,
			Asc:           true,
		},
		Queries: []query.SearchQuery{isReached, isNotPushed, hasPrimaryDomain},
	})
	if err != nil {
		return nil, err
	}
	var errs int
	for _, ms := range unpushedMilestones.Milestones {
		if err = t.pushMilestone(ctx, ms); err != nil {
			errs++
			logging.Warnf("pushing milestone %+v failed: %s", *ms, err.Error())
		}
	}
	if errs > 0 {
		return nil, fmt.Errorf("pushing %d of %d milestones failed", errs, unpushedMilestones.Count)
	}

	return crdb.NewNoOpStatement(timestampEvent), nil
}

func (t *telemetryPusher) pushMilestone(ctx context.Context, ms *query.Milestone) error {
	for _, endpoint := range t.endpoints {
		if err := types.SendJSON(
			ctx,
			webhook.Config{
				CallURL: endpoint,
				Method:  http.MethodPost,
			},
			t.queries.GetFileSystemProvider,
			t.queries.GetLogProvider,
			ms,
			nil,
			t.metricSuccessfulDeliveriesJSON,
			t.metricFailedDeliveriesJSON,
		).WithoutTemplate(); err != nil {
			return err
		}
	}
	return t.commands.MilestonePushed(ctx, ms.InstanceID, ms.MilestoneType, t.endpoints, ms.PrimaryDomain)
}