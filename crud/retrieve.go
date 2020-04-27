package crud

import (
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/queryutil"
)

// DefaultCRUD default crud entity
type DefaultCRUD struct {
	Model           interface{}
	ModelSlice      interface{}
	HasAuthCtl      bool
	EnableChangeLog bool
}

// Retrieve default retrieve func
func Retrieve(tctx application.Context, queryHandler queryutil.QueryHandler) (interface{}, error) {
	return nil, nil

}
