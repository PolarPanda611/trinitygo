package application

import (
	"reflect"
	"sync"
)

// RepositoryPool service pool
type RepositoryPool struct {
	poolMap  map[reflect.Type]*sync.Pool
	repoType []reflect.Type
}

// NewRepositoryPool new pool with init map
func NewRepositoryPool() *RepositoryPool {
	repoPool := new(RepositoryPool)
	repoPool.poolMap = make(map[reflect.Type]*sync.Pool)

	return repoPool

}

// NewRepository add new service
func (s *RepositoryPool) NewRepository(repoType reflect.Type, repoPool *sync.Pool) {
	s.poolMap[repoType] = repoPool
	s.repoType = append(s.repoType, repoType)
}

// GetRepositoryType get repository type list
func (s *RepositoryPool) GetRepositoryType() []reflect.Type {
	return s.repoType
}

// GetRepository get repository with di
func (s *RepositoryPool) GetRepository(repoType reflect.Type, tctx Context, app Application) interface{} {
	repoPool, ok := s.poolMap[repoType]
	if !ok {
		panic("unknown repository name")
	}
	repo := repoPool.Get()
	DiTCtx(repo, tctx)
	return repo
}

// Release release repo
func (s *RepositoryPool) Release(repo interface{}) {
	t := reflect.TypeOf(repo)
	syncpool, ok := s.poolMap[t]
	if !ok {
		return
	}
	DiFree(repo)
	syncpool.Put(repo)
}
