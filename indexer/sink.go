package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
)

var (
	_ pipeline.Sink = (*sink)(nil)
)

func NewSink(db *store.Store, versionNumber int64) *sink {
	return &sink{
		db:            db,
		versionNumber: versionNumber,
	}
}

type sink struct {
	db            *store.Store
	versionNumber int64

	successCount int64
}

func (s *sink) Consume(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.DebugJSON(payload,
		logger.Field("process", "pipeline"),
		logger.Field("stage", "sink"),
		logger.Field("height", payload.CurrentHeight),
	)

	if err := s.setProcessed(payload); err != nil {
		return err

	}

	if err := s.addMetrics(payload); err != nil {
		return err
	}

	s.successCount += 1

	logger.Info(fmt.Sprintf("processing completed [status=success] [height=%d]", payload.CurrentHeight))

	return nil
}

func (s *sink) setProcessed(payload *payload) error {
	payload.Syncable.MarkProcessed(s.versionNumber)
	if err := s.db.Syncables.Save(payload.Syncable); err != nil {
		return errors.Wrap(err, "failed saving syncable in sink")
	}
	return nil
}

func (s *sink) addMetrics(payload *payload) error {
	res, err := s.db.Database.GetTotalSize()
	if err != nil {
		return err
	}

	indexerHeightSuccess.Inc()
	indexerHeightDuration.Observe(payload.Syncable.Duration.Seconds())
	indexerDbSizeAfterHeight.Observe(res.Size)
	return nil
}
