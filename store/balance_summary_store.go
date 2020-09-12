package store

import (
	"fmt"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/jinzhu/gorm"
)

var (
	_ BalanceSummaryStore = (*balanceSummaryStore)(nil)
)

type BalanceSummaryStore interface {
	BaseStore

	Find(*model.BalanceSummary) (*model.BalanceSummary, error)
	FindActivityPeriods(types.SummaryInterval, int64) ([]ActivityPeriodRow, error)
}

func NewBalanceSummaryStore(db *gorm.DB) *balanceSummaryStore {
	return &balanceSummaryStore{scoped(db, model.BalanceSummary{})}
}

type balanceSummaryStore struct {
	baseStore
}

// Find find balance summary by query
func (s balanceSummaryStore) Find(query *model.BalanceSummary) (*model.BalanceSummary, error) {
	var result model.BalanceSummary

	err := s.db.
		Where(query).
		First(&result).
		Error

	return &result, checkErr(err)
}

// FindActivityPeriods Finds activity periods
func (s *balanceSummaryStore) FindActivityPeriods(interval types.SummaryInterval, indexVersion int64) ([]ActivityPeriodRow, error) {
	t := metrics.NewTimer(databaseQueryDuration.WithLabels("BalanceSummaryStore_FindActivityPeriods"))
	defer t.ObserveDuration()

	rows, err := s.db.Raw(balanceSummaryActivityPeriodsQuery, fmt.Sprintf("1%s", interval), interval, indexVersion).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []ActivityPeriodRow
	for rows.Next() {
		var row ActivityPeriodRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}
