package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"

	"net/http"

	"github.com/fxkr/openview/backend/util/handler"
)

// RedisCache is Cache implementation that stores keys on a Redis server.
type RedisCache struct {
	db     redis.Conn
	config RedisCacheConfig
}

type RedisCacheConfig struct {
	Host     string  `json:"host"`
	Network  string  `json:"network"`
	Password *string `json:"password"`
	Prefix   string  `json:"prefix"`
}

// Statically assert that *RedisCache implements Cache.
var _ Cache = (*RedisCache)(nil)

func NewRedisCache(config RedisCacheConfig) (*RedisCache, error) {
	c, err := redis.Dial(config.Network, config.Host)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if config.Password != nil {
		_, err = redis.String(c.Do("AUTH", *config.Password))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &RedisCache{
		db:     c,
		config: config,
	}, nil
}

func (c *RedisCache) Close() {
	c.Close()
}

func (c *RedisCache) Put(key Key, value []byte) error {
	_, err := c.db.Do("SET", c.config.Prefix+key.String(), value)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *RedisCache) GetBytes(key Key, filler func() ([]byte, error)) ([]byte, error) {
	value, err := redis.Bytes(c.db.Do("GET", c.config.Prefix+key.String()))
	if err == nil {
		return value, nil // Cache hit
	}

	value, err = filler()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = c.Put(key, value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return value, nil
}

func (c *RedisCache) GetHandler(key Key, filler func() ([]byte, error), contentType string) (http.Handler, error) {
	bytes, err := c.GetBytes(key, filler)
	if err != nil {
		return nil, err
	}

	return &handler.ByteHandler{
		Bytes:       bytes,
		ContentType: contentType,
	}, nil
}
