package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"

	"net/http"

	"bytes"

	"github.com/fxkr/openview/backend/util/handler"
	"github.com/fxkr/openview/backend/util/safe"
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

func (c *RedisCache) Put(key Key, version Version, value []byte) error {
	dataKey := c.config.Prefix + key.String()
	versionKey := c.config.Prefix + safe.NewKey(key.String(), "ver").String()

	_, err := c.db.Do("MSET", dataKey, value, versionKey, []byte(version.String()))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *RedisCache) GetBytes(key Key, version Version, filler func() (Version, []byte, error)) ([]byte, error) {
	dataKey := c.config.Prefix + key.String()
	versionKey := c.config.Prefix + safe.NewKey(key.String(), "ver").String()

	values, err := redis.ByteSlices(c.db.Do("MGET", versionKey, dataKey))
	if err == nil && len(values) == 2 { // Cache hit?
		if bytes.Equal(values[0], []byte(version.String())) { // Cache up to date?
			cachedBytes := values[1]
			return cachedBytes, nil
		}
	}

	version, value, err := filler()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = c.Put(key, version, value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return value, nil
}

func (c *RedisCache) GetHandler(key Key, version Version, filler func() (Version, []byte, error), contentType string) (http.Handler, error) {
	bytes, err := c.GetBytes(key, version, filler)
	if err != nil {
		return nil, err
	}

	return &handler.ByteHandler{
		Bytes:       bytes,
		ContentType: contentType,
	}, nil
}
