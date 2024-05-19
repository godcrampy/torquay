package counter

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-zookeeper/zk"
)

type Counter struct {
	zkConn *zk.Conn
	zkPath string
}

func NewCounterWithRetry(servers []string, path string) (*Counter, error) {
	const RetryCount = 5
	const RetryDelayS = 5

	var (
		c   *Counter
		err error
	)

	for attempt := 1; attempt <= RetryCount; attempt++ {
		c, err = NewCounter(servers, path)

		if err != nil {
			log.Printf("ERROR: Failed to connect to ZooKeeper (attempt %d/%d): %v", attempt, RetryCount, err)
			time.Sleep(time.Second * RetryDelayS)
			continue
		}

		return c, nil
	}

	return nil, err
}

func NewCounter(servers []string, path string) (*Counter, error) {
	conn, _, err := zk.Connect(servers, time.Second*5)
	if err != nil {
		return nil, err
	}

	c := &Counter{
		zkConn: conn,
		zkPath: path,
	}

	if exists, _, err := conn.Exists(path); err != nil {
		return nil, err
	} else if !exists {
		_, err := conn.Create(path, []byte("0"), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Counter) GetAndIncrement() (int, error) {
	for {
		data, stat, err := c.zkConn.Get(c.zkPath)
		if err != nil {
			return 0, err
		}

		currentValue, err := strconv.Atoi(string(data))
		if err != nil {
			return 0, err
		}

		newValue := currentValue + 1
		newData := []byte(fmt.Sprintf("%d", newValue))

		_, err = c.zkConn.Set(c.zkPath, newData, stat.Version)
		if err == zk.ErrBadVersion {
			continue
		} else if err != nil {
			return 0, err
		}

		return newValue, nil
	}
}

func (c *Counter) Close() {
	c.zkConn.Close()
}
