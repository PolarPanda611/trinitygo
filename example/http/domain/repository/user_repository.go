package repository

import (
	"errors"
	"http/domain/model"
	"math"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/queryutil"
	"github.com/jinzhu/gorm"
)

var (
	_           UserRepository         = new(userRepositoryImpl)
	_userConfig *queryutil.QueryConfig = &queryutil.QueryConfig{
		FilterBackend:       []func(db *gorm.DB) *gorm.DB{},
		PageSize:            20,
		FilterList:          []string{},
		OrderByList:         []string{},
		SearchByList:        []string{},
		PreloadList:         map[string]func(db *gorm.DB) *gorm.DB{},
		FilterCustomizeFunc: map[string]interface{}{},
	}
)

func init() {
	trinitygo.RegisterInstance(func() interface{} {
		repo := new(userRepositoryImpl)
		repo.queryHandler = queryutil.New(_userConfig)
		return repo
	}, "UserRepository")
}

// UserRepository user repository
type UserRepository interface {
	GetUserByID(id int64) (*model.User, error)
	GetUserList(query string) ([]model.User, bool, error)
	CreateUser(*model.User) (*model.User, error)
	UpdateUserByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteUserByID(id int64, dVersion string) error
	GetUserCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error)
}

type userRepositoryImpl struct {
	Tctx         application.Context `autowired:"true" `
	queryHandler queryutil.QueryHandler
}

func (r *userRepositoryImpl) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetUserList(query string) ([]model.User, bool, error) {
	var userList []model.User
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&userList).Error; err != nil {
		return nil, false, err
	}
	return userList, r.queryHandler.IsPaginationOff(), nil
}

func (r *userRepositoryImpl) GetUserCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Model(&model.User{}).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	return count, r.queryHandler.PageNum(), int(math.Ceil(float64(count) / float64(r.queryHandler.PageSize()))), r.queryHandler.PageSize(), nil
}

func (r *userRepositoryImpl) CreateUser(newUser *model.User) (*model.User, error) {
	if err := r.Tctx.DB().Create(newUser).Error; err != nil {
		return nil, err
	}
	return newUser, nil
}

func (r *userRepositoryImpl) UpdateUserByID(id int64, dVersion string, change map[string]interface{}) error {

	updateQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Table(r.Tctx.DB().NewScope(&model.User{}).TableName()).Update(change)
	if err := updateQuery.Error; err != nil {
		return err
	}
	if updateQuery.RowsAffected != 1 {
		return errors.New("update failed , affected zero lines , please refresh the data and retry")
	}
	return nil
}
func (r *userRepositoryImpl) DeleteUserByID(id int64, dVersion string) error {
	deleteQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Delete(&model.User{})
	if err := deleteQuery.Error; err != nil {
		return err
	}
	if deleteQuery.RowsAffected != 1 {
		return errors.New("delete affected zero lines , please refresh the data")
	}
	return nil
}
