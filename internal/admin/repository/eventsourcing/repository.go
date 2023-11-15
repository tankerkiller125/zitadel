package eventsourcing

import (
	"context"

	admin_handler "github.com/zitadel/zitadel/v2/internal/admin/repository/eventsourcing/handler"
	admin_view "github.com/zitadel/zitadel/v2/internal/admin/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/v2/internal/database"
	"github.com/zitadel/zitadel/v2/internal/static"
)

type Config struct {
	Spooler admin_handler.Config
}

func Start(ctx context.Context, conf Config, static static.Storage, dbClient *database.DB) error {
	view, err := admin_view.StartView(dbClient)
	if err != nil {
		return err
	}

	admin_handler.Register(ctx, conf.Spooler, view, static)

	return nil
}
