package repository

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/example/http/domain/object"
	queryutil "github.com/PolarPanda611/trinitygo/queryutils"
	"github.com/jinzhu/gorm"
)

var (
	_           UserRepository         = new(userRepositoryImpl)
	_userConfig *queryutil.QueryConfig = &queryutil.QueryConfig{
		TablePrefix:  "",
		DbBackend:    nil,
		PageSize:     20,
		FilterList:   []string{"user_name", "user_name__ilike"},
		OrderByList:  []string{"id"},
		SearchByList: []string{"user_name", "email"},
		FilterCustomizeFunc: map[string]interface{}{
			"test": func(db *gorm.DB, queryValue string) *gorm.DB {
				fmt.Println("Where xxxxx = ?", queryValue)
				return db.Where("xxxxx = ?", queryValue)
			},
		},
		IsDebug: false,
	}
)

func init() {
	trinitygo.BindContainer(reflect.TypeOf(&userRepositoryImpl{}), &sync.Pool{
		New: func() interface{} {
			repo := new(userRepositoryImpl)
			repo.queryHandler = queryutil.New(_userConfig)
			return repo
		},
	})
}

// UserRepository user repository
type UserRepository interface {
	GetUserByID(id int) (*object.User, error)
	GetUserList(query string) ([]object.User, error)
}

type userRepositoryImpl struct {
	TCtx         application.Context
	queryHandler queryutil.QueryHandler
}

func (r *userRepositoryImpl) GetUserByID(id int) (*object.User, error) {
	var user object.User
	if err := r.TCtx.DB().Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetUserList(query string) ([]object.User, error) {
	var user []object.User
	if err := r.TCtx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
