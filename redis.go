package melock

import (
	"github.com/gomodule/redigo/redis"
)

const unlockScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`

type RedisMeLock struct {
	conn           redis.Conn
	resource       string
	identifier     string
	timeoutSeconds uint
}

func NewRedisMeLock(conn redis.Conn, resource string, identifier string, timeoutSeconds uint) *RedisMeLock {
	return &RedisMeLock{
		conn:           conn,
		resource:       resource,
		identifier:     identifier,
		timeoutSeconds: timeoutSeconds,
	}
}

func (lock *RedisMeLock) acquire() (bool, error) {
	reply, err := redis.String(lock.conn.Do("SET", lock.resource, lock.identifier, "NX", "EX", lock.timeoutSeconds))
	//acquire lock success
	if err == nil && reply == "OK" {
		return true, nil
	}
	//acquire lock fail, lock already exist
	if err == redis.ErrNil {
		return false, nil
	}
	//error
	return false, err
}

func (lock *RedisMeLock) Release() (flag bool, err error) {
	script := redis.NewScript(1, unlockScript)
	reply, err := redis.Int(script.Do(lock.conn, lock.resource, lock.identifier))
	if err == nil && reply > 0 {
		flag = true
	}
	return
}

func Acquire(conn redis.Conn, resource string, identifier string, timeoutSeconds uint) (lock *RedisMeLock, ok bool, err error) {
	lock = NewRedisMeLock(conn, resource, identifier, timeoutSeconds)
	ok, err = lock.acquire()
	return
}
