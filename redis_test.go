package melock

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedisMeLock_Acquire(t *testing.T) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ass := assert.New(t)
	rml := NewRedisMeLock(conn, "stock", "uuid", 2)
	// should lock success
	flag, err := rml.acquire()
	ass.Nil(err)
	ass.True(flag)
	// should lock fail
	flag, err = rml.acquire()
	ass.False(flag)
	ass.Nil(err)
	time.Sleep(time.Duration(2 * time.Second))
	// should lock success after the lock is timeout
	flag, err = rml.acquire()
	ass.Nil(err)
	ass.True(flag)
	flag, err = rml.Release()
	ass.Nil(err)
	ass.True(flag)
}

func TestRedisMeLock_Release(t *testing.T) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ass := assert.New(t)
	rml := NewRedisMeLock(conn, "stock", "uuid", 2)
	// should lock success
	flag, err := rml.acquire()
	ass.Nil(err)
	ass.True(flag)
	flag, err = rml.Release()
	ass.True(flag)
	ass.Nil(err)
}
