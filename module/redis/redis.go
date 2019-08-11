// 参考:
// https://github.com/garyburd/redigo
// https://github.com/boj/redistore

// Copyright 2012 Brian "bojo" Jones. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package redis

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

var ErrNil = redis.ErrNil

// RediStore stores sessions in a redis backend.
type RediStore struct {
	Pool *redis.Pool
}

func dial(network, address, password string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

// NewRediStore returns a new RediStore.
// size: maximum number of idle connections.
func NewRediStore(size int, network, address, password string) (*RediStore, error) {
	return NewRediStoreWithPool(&redis.Pool{
		MaxIdle:     size,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dial(network, address, password)
		},
	})
}

func dialWithDB(network, address, password, DB string) (redis.Conn, error) {
	c, err := dial(network, address, password)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("SELECT", DB); err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}

// NewRediStoreWithDB - like NewRedisStore but accepts `DB` parameter to select
// redis DB instead of using the default one ("0")
func NewRediStoreWithDB(size int, network, address, password, DB string) (*RediStore, error) {
	return NewRediStoreWithPool(&redis.Pool{
		MaxIdle:     size,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dialWithDB(network, address, password, DB)
		},
	})
}

// NewRediStoreWithPool instantiates a RediStore with a *redis.Pool passed in.
func NewRediStoreWithPool(pool *redis.Pool) (*RediStore, error) {
	rs := &RediStore{
		// http://godoc.org/github.com/garyburd/redigo/redis#Pool
		Pool: pool,
	}
	_, err := rs.ping()
	return rs, err
}

// Close closes the underlying *redis.Pool
func (s *RediStore) Close() error {
	return s.Pool.Close()
}

// Delete removes the session from redis.
//
func (s *RediStore) Delete(key string) error {
	conn := s.Pool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", key); err != nil {
		return err
	}
	return nil
}

// ping does an internal ping against a server to check if it is alive.
func (s *RediStore) ping() (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()
	data, err := conn.Do("PING")
	if err != nil || data == nil {
		return false, err
	}
	return (data == "PONG"), nil
}

// save stores the session in redis.
// 以json方式存储
// age -- seconds, 0永久性存储
func (s *RediStore) Set(key string, val interface{}, age int64) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	conn := s.Pool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return err
	}
	_, err = conn.Do("SETEX", key, age, b)
	return err
}

// 读取数据, 请参考json的读法
func (s *RediStore) Scan(key string, val interface{}) error {
	conn := s.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	reply, err := conn.Do("GET", key)
	if err != nil {
		return err
	}
	out, err := redis.Bytes(reply, err)
	if err != nil {
		return err
	}
	return json.Unmarshal(out, val)
}
