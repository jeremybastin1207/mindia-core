package redis

import (
	redigo "github.com/gomodule/redigo/redis"
)

func NewPool(addr string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", addr)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

}
