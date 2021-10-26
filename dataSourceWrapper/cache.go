package data_source_wrapper

import (
	"time"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

type Cache struct {
	cache       interface{}
	updatedAt   time.Time
	lifeTime    time.Duration
	updateCache func() (interface{}, error)
	group       singleflight.Group
}

type CacheMaker struct {
	lifeTime time.Duration
}

func NewCacheMaker(lifeTime time.Duration) *CacheMaker {
	return &CacheMaker{
		lifeTime: lifeTime,
	}
}

func (m CacheMaker) NewDataSourceWrapper(f func() (interface{}, error)) DataSourceWrapper {
	ret := &Cache{
		lifeTime: m.lifeTime,
	}
	ret.updateCache = func() (interface{}, error) {
		var err error
		ret.cache, err = f()
		log.Print("UPDATED CACHE!!")
		ret.updatedAt = time.Now()
		return ret.cache, err
	}
	return ret
}

func (s *Cache) Get() interface{} {
	if time.Since(s.updatedAt) > s.lifeTime {
		_, err, _ := s.group.Do("cache", s.updateCache)
		if err != nil {
			log.Error(errors.Wrap(err, "failed to update cache"))
		}
	}
	return s.cache
}
