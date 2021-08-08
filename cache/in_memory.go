package cache

import (
	"errors"
	"snackable/domain/file"

	"github.com/rs/zerolog/log"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type FileCache interface {
	Get(id string) (file.Model, error)
	Set(id string, value file.Model) error
}

type inMemoryFileCache struct {
	cache map[string]file.Model
}

func NewInMemoryFileCache(cache map[string]file.Model) FileCache {
	internalCache := cache
	if internalCache == nil {
		internalCache = make(map[string]file.Model)
	}

	return inMemoryFileCache{
		cache: internalCache,
	}
}

func (imc inMemoryFileCache) Get(id string) (file.Model, error) {
	if v, ok := imc.cache[id]; ok {
		return v, nil
	}

	return file.Model{}, ErrCacheMiss
}

func (imc inMemoryFileCache) Set(id string, value file.Model) error {
	log.Debug().Str("id", id).Msg("saving to cache")
	imc.cache[id] = value
	return nil
}
