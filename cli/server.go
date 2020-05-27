package cli

import (
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/server"
	"github.com/figment-networks/oasishub-indexer/usecase"
)

func startServer(cfg *config.Config) error {
	client, err := initClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	httpHandlers := usecase.NewHttpHandlers(db, client)

	a := server.New(httpHandlers)
	if err := a.Start(cfg.ListenAddr()); err != nil {
		return err
	}
	return nil
}