package validatorseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/validatordomain"
	"github.com/figment-networks/oasishub-indexer/mappers/validatorseqmapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm"
)

var _ DbRepo = (*dbRepo)(nil)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	Count() (*int64, errors.ApplicationError)
	GetByHeight(types.Height) ([]*validatordomain.ValidatorSeq, errors.ApplicationError)

	// Commands
	Save(*validatordomain.ValidatorSeq) errors.ApplicationError
	Create(*validatordomain.ValidatorSeq) errors.ApplicationError
}

type dbRepo struct {
	client *gorm.DB
}

func NewDbRepo(c *gorm.DB) *dbRepo {
	return &dbRepo{
		client: c,
	}
}

func (r *dbRepo) Exists(h types.Height) bool {
	query := heightQuery(h)
	foundSyncableValidator := orm.ValidatorSeqModel{}

	if err := r.client.Where(&query).Take(&foundSyncableValidator).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count() (*int64, errors.ApplicationError) {
	var count int64
	if err := r.client.Table(orm.ValidatorSeqModel{}.TableName()).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count of validator sequences", errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting count of validator sequences", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetByHeight(h types.Height) ([]*validatordomain.ValidatorSeq, errors.ApplicationError) {
	query := heightQuery(h)
	var seqs []orm.ValidatorSeqModel

	if err := r.client.Where(&query).Find(&seqs).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find validator sequences with height %d", h), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting validator sequences", errors.QueryError, err)
	}

	var resp []*validatordomain.ValidatorSeq
	for _, s := range seqs {
		vs, err := validatorseqmapper.FromPersistence(s)
		if err != nil {
			return nil, err
		}

		resp = append(resp, vs)
	}
	return resp, nil
}

func (r *dbRepo) Save(sv *validatordomain.ValidatorSeq) errors.ApplicationError {
	pr, err := validatorseqmapper.ToPersistence(sv)
	if err != nil {
		return err
	}

	if err := r.client.Save(pr).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not save validator sequence", errors.SaveError, err)
	}
	return nil
}

func (r *dbRepo) Create(sv *validatordomain.ValidatorSeq) errors.ApplicationError {
	b, err := validatorseqmapper.ToPersistence(sv)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create validator sequence", errors.CreateError, err)
	}
	return nil
}

/*************** Private ***************/

func heightQuery(h types.Height) orm.ValidatorSeqModel {
	return orm.ValidatorSeqModel{
		SequenceModel: orm.SequenceModel{
			Height: h,
		},
	}
}