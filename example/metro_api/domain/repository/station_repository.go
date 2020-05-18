package repository

import (
	"errors"
	"math"
	"metro_api/domain/model"
	"metro_api/infra/db"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/queryutil"
	"github.com/jinzhu/gorm"
)

var (
	_              StationRepository      = new(stationRepositoryImpl)
	_stationConfig *queryutil.QueryConfig = &queryutil.QueryConfig{

		FilterBackend: []func(db *gorm.DB) *gorm.DB{},
		PageSize:      20,
		FilterList:    []string{"deleted_time__isnull", "code__ilike", "name__ilike", "description__ilike", "next_station_id", "line_id"},
		OrderByList:   []string{},
		SearchByList:  []string{},
		PreloadList: map[string]func(db *gorm.DB) *gorm.DB{
			"Line":        nil,
			"NextStation": nil,
			"CreateUser":  nil,
			"UpdateUser":  nil,
			"DeleteUser":  nil,
		},
		FilterCustomizeFunc: map[string]interface{}{},
	}
)

func init() {
	trinitygo.RegisterInstance(func() interface{} {
		repo := new(stationRepositoryImpl)
		repo.queryHandler = queryutil.New(_stationConfig)
		return repo
	}, "StationRepository")
}

// StationRepository station repository
type StationRepository interface {
	GetNextSeq() (string, error)
	GetStationByID(id int64) (*model.Station, error)
	GetStationList(query string) ([]model.Station, bool, error)
	CreateStation(*model.Station) (*model.Station, error)
	UpdateStationByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteStationByID(id int64, dVersion string) error
	GetStationCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error)
}

type stationRepositoryImpl struct {
	Tctx         application.Context `autowired:"true" `
	queryHandler queryutil.QueryHandler
}

func (r *stationRepositoryImpl) GetNextSeq() (string, error) {
	createsequencesql := "CREATE SEQUENCE IF NOT EXISTS seq_station_number MINVALUE 0 MAXVALUE 9999999 START 1  ;"
	if err := r.Tctx.DB().Exec(createsequencesql).Error; err != nil {
		return "", err
	}
	var seq db.SequenceResult
	if err := r.Tctx.DB().Raw("SELECT nextval('seq_station_number');").Scan(&seq).Error; err != nil {
		return "", err
	}
	return seq.Nextval, nil
}

func (r *stationRepositoryImpl) GetStationByID(id int64) (*model.Station, error) {
	var station model.Station
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ?", id).First(&station).Error; err != nil {
		return nil, err
	}
	return &station, nil
}

func (r *stationRepositoryImpl) GetStationList(query string) ([]model.Station, bool, error) {
	var stationList []model.Station
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&stationList).Error; err != nil {
		return nil, false, err
	}
	return stationList, r.queryHandler.IsPaginationOff(), nil
}

func (r *stationRepositoryImpl) GetStationCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Model(&model.Station{}).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	return count, r.queryHandler.PageNum(), int(math.Ceil(float64(count) / float64(r.queryHandler.PageSize()))), r.queryHandler.PageSize(), nil
}

func (r *stationRepositoryImpl) CreateStation(newStation *model.Station) (*model.Station, error) {
	if err := r.Tctx.DB().Create(newStation).Error; err != nil {
		return nil, err
	}
	return newStation, nil

}

func (r *stationRepositoryImpl) UpdateStationByID(id int64, dVersion string, change map[string]interface{}) error {

	updateQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Table(r.Tctx.DB().NewScope(&model.Station{}).TableName()).Update(change)
	if err := updateQuery.Error; err != nil {
		return err
	}
	if updateQuery.RowsAffected != 1 {
		return errors.New("update affected zero lines , please refresh the data")
	}

	return nil
}
func (r *stationRepositoryImpl) DeleteStationByID(id int64, dVersion string) error {
	deleteQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Delete(&model.Station{})
	if err := deleteQuery.Error; err != nil {
		return err
	}
	if deleteQuery.RowsAffected != 1 {
		return errors.New("delete affected zero lines , please refresh the data")
	}

	return nil
}
