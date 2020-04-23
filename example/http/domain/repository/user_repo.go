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
		FilterBackend: nil,
		PageSize:      20,
		FilterList:    []string{"user_name"},
		OrderByList:   []string{"id"},
		SearchByList:  []string{"user_name", "email"},
		PreloadList: map[string]func(db *gorm.DB) *gorm.DB{
			"Languages": nil,
		},
		FilterCustomizeFunc: map[string]interface{}{
			"test": func(db *gorm.DB, queryValue string) *gorm.DB {
				fmt.Println("Where xxxxx = ?", queryValue)
				return db.Where("xxxxx = ?", queryValue)
			},
		},
	}
)

func init() {
	trinitygo.BindContainer(reflect.TypeOf(&userRepositoryImpl{}), &sync.Pool{
		New: func() interface{} {
			repo := new(userRepositoryImpl)
			repo.queryHandler = queryutil.New(_userConfig)
			return repo
		},
	},
		"UserRepository",
	)
}

// UserRepository user repository
type UserRepository interface {
	GetUserByID(id int) (*object.User, error)
	GetUserList(query string) ([]object.User, error)
}

type userRepositoryImpl struct {
	Tctx         application.Context `autowired:"true"`
	queryHandler queryutil.QueryHandler
}

func (r *userRepositoryImpl) GetUserByID(id int) (*object.User, error) {
	fmt.Println("repo run ")
	var user object.User
	if err := r.Tctx.DB().Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetUserList(query string) ([]object.User, error) {
	fmt.Println("repo run ")
	var user []object.User
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
