package repository

import (
	"fmt"

	"github.com/PolarPanda611/trinitygo"

	"github.com/PolarPanda611/trinitygo/application"
)

var _ UserRepo = new(UserRepoImpl)

func init() {
	trinitygo.RegisterInstance(UserRepoImpl{})
}

// UserRepo user repo
type UserRepo interface {
	Print() string
}

// UserRepoImpl user repo impl
type UserRepoImpl struct {
	TContext application.Context `autowired:"true"`
}

// Print user repo impl print method
func (r *UserRepoImpl) Print() string {
	return "userrepo" + fmt.Sprintf("%v", r.TContext.Runtime())
}
