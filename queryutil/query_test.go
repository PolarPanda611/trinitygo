package queryutil

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestQueryHandler(t *testing.T) {
	query := "query=123&Query=1234&code=sasfaf&query__ilike=234&OrderBy=id&SearchBy=124"
	config := &QueryConfig{
		FilterBackend:       nil,
		PageSize:            20,
		FilterList:          []string{"query", "query__ilike"},
		OrderByList:         []string{"id"},
		SearchByList:        []string{"code", "name"},
		FilterCustomizeFunc: make(map[string]interface{}),
		PreloadList: map[string]func(db *gorm.DB) *gorm.DB{
			"Test": nil,
		},
		IsDebug: true,
	}
	x := New(config).Handle(query)
	fmt.Println(len(x))
}

func TestQueryHandlerWithPagi(t *testing.T) {
	query := "query=123&Query=1234&code=sasfaf&query__ilike=234&OrderBy=id&SearchBy=124&PageSize=30&PageNum=4&test=wqrqwr"
	config := &QueryConfig{
		FilterBackend: nil,
		PageSize:      20,
		FilterList:    []string{"query", "query__ilike"},
		OrderByList:   []string{"id"},
		SearchByList:  []string{"code", "name"},
		PreloadList: map[string]func(db *gorm.DB) *gorm.DB{
			"Test": func(db *gorm.DB) *gorm.DB {
				return db.Where("delete is null ")
			},
		},
		FilterCustomizeFunc: map[string]interface{}{
			"test": func(db *gorm.DB, queryValue string) *gorm.DB {
				fmt.Println("Where xxxxx = ?", queryValue)
				return db.Where("xxxxx = ?", queryValue)
			},
		},
		IsDebug: true,
	}
	x := New(config).HandleWithPagination(query)

	fmt.Println(len(x))
}

type TestStruct struct {
	ModelID int
	Model   *TestModelStruct `remote_resource:"UserResource" remote_condition:"xxx=xxx"`
}

type TestModelStruct struct {
	ID   int
	Name string
}

func TestHandleRemotePreloaderType(t *testing.T) {

	q := new(queryRepositoryImpl)

	testPtr := TestStruct{
		ModelID: 1,
	}
	err := q.HandleRemotePreloader(testPtr)
	assert.Equal(t, errors.New("must be ptr or slice "), err, "result should be ptr ")
	a := 1
	testPtrStruct := &a
	err = q.HandleRemotePreloader(testPtrStruct)
	assert.Equal(t, errors.New("must be ptr of struct or slice "), err, "result should be ptr of struct ")

	cc := []TestStruct{TestStruct{
		ModelID: 1,
	}}
	err = q.HandleRemotePreloader(&cc)
	assert.Equal(t, nil, err, "result should be ptr of struct ")

	// b := []string{"1234", "2345"}
	// err = q.HandleRemotePreloader(&b)
	// assert.Equal(t, errors.New("must be ptr of struct or slice "), err, "result should be ptr of struct ")
}

func TestHandleRemotePreloaderStruct(t *testing.T) {
	remoteDB := []TestModelStruct{
		TestModelStruct{
			ID:   1,
			Name: "test1",
		},
		TestModelStruct{
			ID:   2,
			Name: "test2",
		},
		TestModelStruct{
			ID:   3,
			Name: "test3",
		},
		TestModelStruct{
			ID:   4,
			Name: "test4",
		},
	}
	q := new(queryRepositoryImpl)
	q.remotePreloadlist = []RemotePreloader{
		RemotePreloader{
			Column: "Model",
			PreloadFunc: func(args interface{}, Conditions ...string) (interface{}, error) {
				for _, v := range remoteDB {
					if v.ID == args.(int) {
						return &v, nil
					}
					continue
				}
				return nil, errors.New("not found")
			},
			Condition: "string",
		},
	}

	result := &TestStruct{
		ModelID: 2,
	}
	err := q.HandleRemotePreloader(result)
	assert.Equal(t, nil, err, "result should be ptr of struct ")
	fmt.Println(result)
	assert.Equal(t, 2, result.Model.ID, "result should be ptr of struct ")
	assert.Equal(t, "test2", result.Model.Name, "result should be ptr of struct ")
}
