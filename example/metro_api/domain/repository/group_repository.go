package repository

import (
	"errors"
	"math"
	"metro_api/domain/model"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/queryutil"
	"github.com/jinzhu/gorm"
)

var (
	_            GroupRepository        = new(groupRepositoryImpl)
	_groupConfig *queryutil.QueryConfig = &queryutil.QueryConfig{
		FilterBackend: []func(db *gorm.DB) *gorm.DB{},
		PageSize:      20,
		FilterList:    []string{"deleted_time__isnull", "code__ilike", "name__ilike"},
		OrderByList:   []string{},
		SearchByList:  []string{},
		PreloadList: map[string]func(db *gorm.DB) *gorm.DB{
			"CreateUser": nil,
			"UpdateUser": nil,
			"DeleteUser": nil,
		},
		FilterCustomizeFunc: map[string]interface{}{},
	}
)

func init() {
	trinitygo.RegisterInstance(func() interface{} {
		repo := new(groupRepositoryImpl)
		repo.queryHandler = queryutil.New(_groupConfig)
		return repo
	}, "GroupRepository")
}

// GroupRepository group repository
type GroupRepository interface {
	GetGroupByID(id int64) (*model.Group, error)
	GetGroupList(query string) ([]model.Group, bool, error)
	CreateGroup(*model.Group) (*model.Group, error)
	UpdateGroupByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteGroupByID(id int64, dVersion string) error
	GetGroupCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error)
}

type groupRepositoryImpl struct {
	Tctx         application.Context `autowired:"true" `
	queryHandler queryutil.QueryHandler
}

func (r *groupRepositoryImpl) GetGroupByID(id int64) (*model.Group, error) {
	var group model.Group
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ?", id).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepositoryImpl) GetGroupList(query string) ([]model.Group, bool, error) {
	var groupList []model.Group
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&groupList).Error; err != nil {
		return nil, false, err
	}
	return groupList, r.queryHandler.IsPaginationOff(), nil
}

func (r *groupRepositoryImpl) GetGroupCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Model(&model.Group{}).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	return count, r.queryHandler.PageNum(), int(math.Ceil(float64(count) / float64(r.queryHandler.PageSize()))), r.queryHandler.PageSize(), nil
}

func (r *groupRepositoryImpl) CreateGroup(newGroup *model.Group) (*model.Group, error) {
	if err := r.Tctx.DB().Create(newGroup).Error; err != nil {
		return nil, err
	}
	return newGroup, nil
}

func (r *groupRepositoryImpl) UpdateGroupByID(id int64, dVersion string, change map[string]interface{}) error {

	updateQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Table(r.Tctx.DB().NewScope(&model.Group{}).TableName()).Update(change)
	if err := updateQuery.Error; err != nil {
		return err
	}
	if updateQuery.RowsAffected != 1 {
		return errors.New("update affected zero lines , please refresh the data")
	}
	return nil
}
func (r *groupRepositoryImpl) DeleteGroupByID(id int64, dVersion string) error {
	deleteQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Delete(&model.Group{})
	if err := deleteQuery.Error; err != nil {
		return err
	}
	if deleteQuery.RowsAffected != 1 {
		return errors.New("delete affected zero lines , please refresh the data")
	}
	return nil
}
