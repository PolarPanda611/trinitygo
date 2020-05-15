package service

import (
	"strconv"

	"github.com/PolarPanda611/trinitygo/example/http/domain/model"

	"github.com/PolarPanda611/trinitygo/example/http/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
)

var _ LanguageService = new(languageServiceImpl)

func init() {
	trinitygo.RegisterInstance(languageServiceImpl{}, "LanguageService")
}

// LanguageService  service interface
type LanguageService interface {
	GetLanguageByID(id int64) (*model.Language, error)
	GetLanguageList(query string) (interface{}, error)
	CreateLanguage(*model.Language) (*model.Language, error)
	UpdateLanguageByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteLanguageByID(id int64, dVersion string) error
}

type languageServiceImpl struct {
	LanguageRepo repository.LanguageRepository `autowired:"true"  resource:"LanguageRepository"`
	Tctx         application.Context           `autowired:"true"`
}

func (s *languageServiceImpl) GetLanguageByID(id int64) (*model.Language, error) {
	return s.LanguageRepo.GetLanguageByID(id)
}
func (s *languageServiceImpl) GetLanguageList(query string) (interface{}, error) {
	res, err := s.LanguageRepo.GetLanguageList(query)
	if err != nil {
		return nil, err
	}
	IsOff, _ := strconv.ParseBool(s.Tctx.GinCtx().Query("PaginationOff"))
	if IsOff {
		return res, nil
	}
	count, currentPage, totalPage, pageSize, err := s.LanguageRepo.GetLanguageCount(query)
	if err != nil {
		return nil, err
	}
	resWithPagination := map[string]interface{}{
		"data":         res,
		"current_page": currentPage,
		"total_count":  count,
		"total_page":   totalPage,
		"page_size":    pageSize,
	}
	return resWithPagination, nil
}

func (s *languageServiceImpl) CreateLanguage(newLanguage *model.Language) (*model.Language, error) {
	return s.LanguageRepo.CreateLanguage(newLanguage)
}

func (s *languageServiceImpl) UpdateLanguageByID(id int64, dVersion string, change map[string]interface{}) error {
	return s.LanguageRepo.UpdateLanguageByID(id, dVersion, change)
}

func (s *languageServiceImpl) DeleteLanguageByID(id int64, dVersion string) error {
	return s.LanguageRepo.DeleteLanguageByID(id, dVersion)
}
