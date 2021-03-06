package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

type SummarizeCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *summarizeUseCase
}

func NewSummarizeCmdHandler(cfg *config.Config, db *store.Store, c *client.Client) *SummarizeCmdHandler {
	return &SummarizeCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *SummarizeCmdHandler) Handle(ctx context.Context) {
	logger.Info(fmt.Sprintf("summarizing indexer use case [handler=cmd]"))

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *SummarizeCmdHandler) getUseCase() *summarizeUseCase {
	if h.useCase == nil {
		h.useCase = NewSummarizeUseCase(h.cfg, h.db)
	}
	return h.useCase
}

