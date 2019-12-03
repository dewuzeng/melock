# melock
Distributed Locks using Go and Redis.

```go
package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/somecodeio/melock"
)

func main() {

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	rml, ok, err := melock.Acquire(conn, "order", "uuid", 2)
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("failed to lock")
		return
	}
	fmt.Println("locked successfully")
	ok, err = rml.Release()
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("failed to unlock")
		return
	}
	fmt.Println("unlocked successfully")
}

```
