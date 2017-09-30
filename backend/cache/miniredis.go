package cache

import (
	"github.com/alicebob/miniredis"
	"github.com/pkg/errors"

	"github.com/fxkr/openview/backend/util"
)

// MiniRedisCache is Cache implementation that uses a miniredis-based embedded fake Redis server.
type MiniRedisCache struct {
	*RedisCache

	miniredis *miniredis.Miniredis
}

type MiniRedisCacheConfig struct {
}

// Statically assert that *MiniRedisCache implements Cache.
var _ Cache = (*MiniRedisCache)(nil)

func NewMiniRedisCache(config MiniRedisCacheConfig) (*MiniRedisCache, error) {
	s := miniredis.NewMiniRedis()

	randomPassword, err := util.NewPassword()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	s.RequireAuth(randomPassword)

	err = s.Start()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rc, err := NewRedisCache(RedisCacheConfig{
		Network:  "tcp",
		Host:     s.Addr(),
		Password: &randomPassword,
	})
	if err != nil {
		s.Close()
		return nil, errors.WithStack(err)
	}

	return &MiniRedisCache{
		RedisCache: rc,
		miniredis:  s,
	}, nil
}

func (c *MiniRedisCache) Close() {
	c.RedisCache.Close()
	c.miniredis.Close()
}
