package eventstore_test

import (
	"context"
	"sync"
	"testing"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestCRDB_Push_OneAggregate(t *testing.T) {
	type args struct {
		ctx      context.Context
		commands []eventstore.Command
		// uniqueConstraints    *eventstore.UniqueConstraint
		uniqueDataType       string
		uniqueDataField      string
		uniqueDataInstanceID string
	}
	type eventsRes struct {
		pushedEventsCount int
		uniqueCount       int
		assetCount        int
		aggType           eventstore.AggregateType
		aggIDs            database.TextArray[string]
	}
	type res struct {
		wantErr   bool
		eventsRes eventsRes
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "push 1 event",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(t, "1"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					aggIDs:            []string{"1"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
		{
			name: "push two events on agg",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(t, "6"),
					generateCommand(t, "6"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 2,
					aggIDs:            []string{"6"},
					aggType:           eventstore.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "failed push because context canceled",
			args: args{
				ctx: canceledCtx(),
				commands: []eventstore.Command{
					generateCommand(t, "9"),
				},
			},
			res: res{
				wantErr: true,
				eventsRes: eventsRes{
					pushedEventsCount: 0,
					aggIDs:            []string{"9"},
					aggType:           eventstore.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "push 1 event and add unique constraint",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(t, "10",
						generateAddUniqueConstraint("usernames", "field"),
					),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					uniqueCount:       1,
					aggIDs:            []string{"10"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
		{
			name: "push 1 event and remove unique constraint",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(t, "11",
						generateRemoveUniqueConstraint("usernames", "testremove"),
					),
				},
				uniqueDataType:  "usernames",
				uniqueDataField: "testremove",
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					uniqueCount:       0,
					aggIDs:            []string{"11"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
		{
			name: "push 1 event and remove instance unique constraints",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(t, "12",
						generateRemoveUniqueConstraint("", ""),
						// generateRemoveInstanceUniqueConstraints("instanceID"),
					),
				},
				uniqueDataType:       "usernames",
				uniqueDataField:      "testremove",
				uniqueDataInstanceID: "instanceID",
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					uniqueCount:       0,
					aggIDs:            []string{"12"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
		{
			name: "push 1 event and add asset",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(t, "13"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					assetCount:        1,
					aggIDs:            []string{"13"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
		{
			name: "push 1 event and remove asset",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(t, "14"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					assetCount:        0,
					aggIDs:            []string{"14"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
	}
	for _, tt := range tests {
		for pusherName, pusher := range pushers {
			t.Run(pusherName+"/"+tt.name, func(t *testing.T) {
				t.Cleanup(cleanupEventstore)
				db := eventstore.NewEventstore(
					&eventstore.Config{
						Querier: queriers["v2"],
						Pusher:  pusher,
					},
				)

				if tt.args.uniqueDataType != "" && tt.args.uniqueDataField != "" {
					err := fillUniqueData(tt.args.uniqueDataType, tt.args.uniqueDataField, tt.args.uniqueDataInstanceID)
					if err != nil {
						t.Error("unable to prefill insert unique data: ", err)
						return
					}
				}
				if _, err := db.Push(tt.args.ctx, tt.args.commands...); (err != nil) != tt.res.wantErr {
					t.Errorf("CRDB.Push() error = %v, wantErr %v", err, tt.res.wantErr)
				}

				assertEventCount(t,
					database.TextArray[eventstore.AggregateType]{tt.res.eventsRes.aggType},
					tt.res.eventsRes.aggIDs,
					tt.res.eventsRes.pushedEventsCount,
				)

				assertUniqueConstraint(t, tt.args.commands, tt.res.eventsRes.uniqueCount)
				// if tt.args.uniqueConstraints != nil {
				// 	countUniqueRow := testCRDBClient.QueryRow("SELECT COUNT(*) FROM eventstore.unique_constraints where unique_type = $1 AND unique_field = $2 AND instance_id = $3", tt.args.uniqueConstraints.UniqueType, tt.args.uniqueConstraints.UniqueField, tt.args.uniqueConstraints.InstanceID)
				// 	var uniqueCount int
				// 	err := countUniqueRow.Scan(&uniqueCount)
				// 	if err != nil {
				// 		t.Error("unable to query inserted rows: ", err)
				// 		return
				// 	}
				// 	if uniqueCount != tt.res.eventsRes.uniqueCount {
				// 		t.Errorf("expected unique count %d got %d", tt.res.eventsRes.uniqueCount, uniqueCount)
				// 	}
				// }
			})
		}
	}
}

func TestCRDB_Push_MultipleAggregate(t *testing.T) {
	type args struct {
		commands []eventstore.Command
	}
	type eventsRes struct {
		pushedEventsCount int
		aggType           database.TextArray[eventstore.AggregateType]
		aggID             database.TextArray[string]
	}
	type res struct {
		wantErr   bool
		eventsRes eventsRes
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "push two aggregates",
			args: args{
				commands: []eventstore.Command{
					generateCommand(t, "100"),
					generateCommand(t, "101"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 2,
					aggID:             []string{"100", "101"},
					aggType:           database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "push two aggregates both multiple events",
			args: args{
				commands: []eventstore.Command{
					generateCommand(t, "102"),
					generateCommand(t, "102"),
					generateCommand(t, "103"),
					generateCommand(t, "103"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 4,
					aggID:             []string{"102", "103"},
					aggType:           database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "push two aggregates mixed multiple events",
			args: args{
				commands: []eventstore.Command{
					generateCommand(t, "106"),
					generateCommand(t, "106"),
					generateCommand(t, "106"),
					generateCommand(t, "106"),
					generateCommand(t, "107"),
					generateCommand(t, "107"),
					generateCommand(t, "107"),
					generateCommand(t, "107"),
					generateCommand(t, "108"),
					generateCommand(t, "108"),
					generateCommand(t, "108"),
					generateCommand(t, "108"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 12,
					aggID:             []string{"106", "107", "108"},
					aggType:           database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
	}
	for _, tt := range tests {
		for pusherName, pusher := range pushers {
			t.Run(pusherName+"/"+tt.name, func(t *testing.T) {
				t.Cleanup(cleanupEventstore)

				db := eventstore.NewEventstore(
					&eventstore.Config{
						Querier: queriers["v2"],
						Pusher:  pusher,
					},
				)
				if _, err := db.Push(context.Background(), tt.args.commands...); (err != nil) != tt.res.wantErr {
					t.Errorf("CRDB.Push() error = %v, wantErr %v", err, tt.res.wantErr)
				}

				countRow := testCRDBClient.QueryRow("SELECT COUNT(*) FROM eventstore.events where aggregate_type = ANY($1) AND aggregate_id = ANY($2)", tt.res.eventsRes.aggType, tt.res.eventsRes.aggID)
				var count int
				err := countRow.Scan(&count)
				if err != nil {
					t.Error("unable to query inserted rows: ", err)
					return
				}
				if count != tt.res.eventsRes.pushedEventsCount {
					t.Errorf("expected push count %d got %d", tt.res.eventsRes.pushedEventsCount, count)
				}
			})
		}
	}
}

func TestCRDB_Push_Parallel(t *testing.T) {
	type args struct {
		commands [][]eventstore.Command
	}
	type eventsRes struct {
		pushedEventsCount int
		aggTypes          database.TextArray[eventstore.AggregateType]
		aggIDs            database.TextArray[string]
	}
	type res struct {
		errCount  int
		eventsRes eventsRes
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "clients push different aggregates",
			args: args{
				commands: [][]eventstore.Command{
					{
						generateCommand(t, "200"),
						generateCommand(t, "200"),
						generateCommand(t, "200"),
						generateCommand(t, "201"),
						generateCommand(t, "201"),
						generateCommand(t, "201"),
					},
					{
						generateCommand(t, "202"),
						generateCommand(t, "203"),
						generateCommand(t, "203"),
					},
				},
			},
			res: res{
				errCount: 0,
				eventsRes: eventsRes{
					aggIDs:            []string{"200", "201", "202", "203"},
					pushedEventsCount: 9,
					aggTypes:          database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "clients push same aggregates",
			args: args{
				commands: [][]eventstore.Command{
					{
						generateCommand(t, "204"),
						generateCommand(t, "204"),
					},
					{
						generateCommand(t, "204"),
						generateCommand(t, "204"),
					},
				},
			},
			res: res{
				errCount: 1,
				eventsRes: eventsRes{
					aggIDs:            []string{"204"},
					pushedEventsCount: 2,
					aggTypes:          database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "clients push different aggregates",
			args: args{
				commands: [][]eventstore.Command{
					{
						generateCommand(t, "207"),
						generateCommand(t, "207"),
						generateCommand(t, "207"),
						generateCommand(t, "207"),
						generateCommand(t, "207"),
						generateCommand(t, "207"),
					},
					{
						generateCommand(t, "208"),
						generateCommand(t, "208"),
						generateCommand(t, "208"),
						generateCommand(t, "208"),
						generateCommand(t, "208"),
					},
				},
			},
			res: res{
				errCount: 0,
				eventsRes: eventsRes{
					aggIDs:            []string{"207", "208"},
					pushedEventsCount: 11,
					aggTypes:          database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
	}
	for _, tt := range tests {
		for pusherName, pusher := range pushers {
			t.Run(pusherName+"/"+tt.name, func(t *testing.T) {
				t.Cleanup(cleanupEventstore)

				db := eventstore.NewEventstore(
					&eventstore.Config{
						Querier: queriers["v2"],
						Pusher:  pusher,
					},
				)

				errs := pushAggregates(db, tt.args.commands)

				if len(errs) != tt.res.errCount {
					t.Errorf("eventstore.Push() error count = %d, wanted err count %d, errs: %v", len(errs), tt.res.errCount, errs)
				}

				assertEventCount(t, tt.res.eventsRes.aggTypes, tt.res.eventsRes.aggIDs, tt.res.eventsRes.pushedEventsCount)
			})
		}
	}
}

func TestCRDB_Push_ResourceOwner(t *testing.T) {
	type args struct {
		commands []eventstore.Command
	}
	type res struct {
		resourceOwners database.TextArray[string]
	}
	type fields struct {
		aggregateIDs  database.TextArray[string]
		aggregateType string
	}
	tests := []struct {
		name   string
		args   args
		res    res
		fields fields
	}{
		{
			name: "two events of same aggregate same resource owner",
			args: args{
				commands: []eventstore.Command{
					generateCommand(t, "500", func(e *testEvent) { e.Agg.ResourceOwner = "caos" }),
					generateCommand(t, "500", func(e *testEvent) { e.Agg.ResourceOwner = "caos" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"500"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos"},
			},
		},
		{
			name: "two events of different aggregate same resource owner",
			args: args{
				commands: []eventstore.Command{
					generateCommand(t, "501", func(e *testEvent) { e.Agg.ResourceOwner = "caos" }),
					generateCommand(t, "502", func(e *testEvent) { e.Agg.ResourceOwner = "caos" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"501", "502"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos"},
			},
		},
		{
			name: "two events of different aggregate different resource owner",
			args: args{
				commands: []eventstore.Command{
					generateCommand(t, "503", func(e *testEvent) { e.Agg.ResourceOwner = "caos" }),
					generateCommand(t, "504", func(e *testEvent) { e.Agg.ResourceOwner = "zitadel" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"503", "504"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "zitadel"},
			},
		},
		{
			name: "events of different aggregate different resource owner",
			args: args{
				commands: []eventstore.Command{
					generateCommand(t, "505", func(e *testEvent) { e.Agg.ResourceOwner = "caos" }),
					generateCommand(t, "505", func(e *testEvent) { e.Agg.ResourceOwner = "caos" }),
					generateCommand(t, "506", func(e *testEvent) { e.Agg.ResourceOwner = "zitadel" }),
					generateCommand(t, "506", func(e *testEvent) { e.Agg.ResourceOwner = "zitadel" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"505", "506"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos", "zitadel", "zitadel"},
			},
		},
		{
			name: "events of different aggregate different resource owner per event",
			args: args{
				commands: []eventstore.Command{
					generateCommand(t, "507", func(e *testEvent) { e.Agg.ResourceOwner = "caos" }),
					generateCommand(t, "507", func(e *testEvent) { e.Agg.ResourceOwner = "ignored" }),
					generateCommand(t, "508", func(e *testEvent) { e.Agg.ResourceOwner = "zitadel" }),
					generateCommand(t, "508", func(e *testEvent) { e.Agg.ResourceOwner = "ignored" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"507", "508"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos", "zitadel", "zitadel"},
			},
		},
		{
			name: "events of one aggregate different resource owner per event",
			args: args{
				commands: []eventstore.Command{
					generateCommand(t, "509", func(e *testEvent) { e.Agg.ResourceOwner = "caos" }),
					generateCommand(t, "509", func(e *testEvent) { e.Agg.ResourceOwner = "ignored" }),
					generateCommand(t, "509", func(e *testEvent) { e.Agg.ResourceOwner = "ignored" }),
					generateCommand(t, "509", func(e *testEvent) { e.Agg.ResourceOwner = "ignored" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"509"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos", "caos", "caos"},
			},
		},
	}
	for _, tt := range tests {
		for pusherName, pusher := range pushers {
			t.Run(pusherName+"/"+tt.name, func(t *testing.T) {
				t.Cleanup(cleanupEventstore)

				db := eventstore.NewEventstore(
					&eventstore.Config{
						Querier: queriers["v2"],
						Pusher:  pusher,
					},
				)

				events, err := db.Push(context.Background(), tt.args.commands...)
				if err != nil {
					t.Errorf("CRDB.Push() error = %v", err)
				}

				if len(events) != len(tt.res.resourceOwners) {
					t.Errorf("length of events (%d) and resource owners (%d) must be equal", len(events), len(tt.res.resourceOwners))
					return
				}

				for i, event := range events {
					if event.Aggregate().ResourceOwner != tt.res.resourceOwners[i] {
						t.Errorf("resource owner not expected want: %q got: %q", tt.res.resourceOwners[i], event.Aggregate().ResourceOwner)
					}
				}

				assertResourceOwners(t, testCRDBClient, tt.res.resourceOwners, tt.fields.aggregateIDs, tt.fields.aggregateType)
			})
		}
	}
}

func pushAggregates(pusher eventstore.Pusher, aggregateCommands [][]eventstore.Command) []error {
	wg := sync.WaitGroup{}
	errs := make([]error, 0)
	errsMu := sync.Mutex{}
	for _, commands := range aggregateCommands {
		wg.Add(1)
		go func(events []eventstore.Command) {
			_, err := pusher.Push(context.Background(), events...)
			if err != nil {
				errsMu.Lock()
				errs = append(errs, err)
				errsMu.Unlock()
			}

			wg.Done()
		}(commands)
	}
	wg.Wait()

	return errs
}

func assertResourceOwners(t *testing.T, db *database.DB, resourceOwners, aggregateIDs database.TextArray[string], aggregateType string) {
	rows, err := testCRDBClient.Query("SELECT resource_owner FROM eventstore.events WHERE aggregate_type = $1 AND aggregate_id = ANY($2) ORDER BY created_at", aggregateType, aggregateIDs)
	if err != nil {
		t.Error("unable to query inserted rows: ", err)
		return
	}
	defer rows.Close()

	eventCount := 0
	for i := 0; rows.Next(); i++ {
		var resourceOwner string
		err = rows.Scan(&resourceOwner)
		if err != nil {
			t.Error("unable to scan row: ", err)
			return
		}
		if resourceOwner != resourceOwners[i] {
			t.Errorf("unexpected resource owner in queried event. want %q, got: %q", resourceOwners[i], resourceOwner)
		}
		eventCount++
	}

	if err := rows.Err(); err != nil {
		t.Errorf("unexpected rows.Err: %v", err)
	}

	if eventCount != len(resourceOwners) {
		t.Errorf("wrong queried event count: want %d, got %d", len(resourceOwners), eventCount)
	}
}

func assertEventCount(t *testing.T, aggTypes database.TextArray[eventstore.AggregateType], aggIDs database.TextArray[string], pushedEventsCount int) {
	row := testCRDBClient.QueryRow("SELECT count(*) FROM eventstore.events where aggregate_type = ANY($1) AND aggregate_id = ANY($2)", aggTypes, aggIDs)

	var count int
	if err := row.Scan(&count); err != nil {
		t.Errorf("unexpected err in row.Scan: %v", err)
		return
	}

	if count != pushedEventsCount {
		t.Errorf("expected push count %d got %d", pushedEventsCount, count)
	}
}

func assertUniqueConstraint(t *testing.T, commands []eventstore.Command, expectedCount int) {
	var uniqueConstraint *eventstore.UniqueConstraint
	for _, command := range commands {
		if e := command.(*testEvent); len(e.uniqueConstraints) > 0 {
			uniqueConstraint = e.uniqueConstraints[0]
			break
		}
	}
	if uniqueConstraint == nil {
		return
	}

	countUniqueRow := testCRDBClient.QueryRow("SELECT COUNT(*) FROM eventstore.unique_constraints where unique_type = $1 AND unique_field = $2", uniqueConstraint.UniqueType, uniqueConstraint.UniqueField)
	var uniqueCount int
	err := countUniqueRow.Scan(&uniqueCount)
	if err != nil {
		t.Error("unable to query inserted rows: ", err)
		return
	}
	if uniqueCount != expectedCount {
		t.Errorf("expected unique count %d got %d", expectedCount, uniqueCount)
	}
}