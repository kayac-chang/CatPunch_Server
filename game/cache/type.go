package cache

import (
	"fmt"
	"time"

	"github.com/YWJSonic/ServerUtility/cacheinfo"
	"github.com/YWJSonic/ServerUtility/code"
	"github.com/YWJSonic/ServerUtility/messagehandle"
	"github.com/gomodule/redigo/redis"
)

// Setting ...
type Setting struct {
	// ConnectTimeout, ReadTimeout, WriteTimeout time.Duration
	CacheDeleteTime time.Duration
	URL             string
}

// GameCache ICache
type GameCache struct {
	cachePool *redis.Pool
	Setting   Setting
}

// GetCachePool ...
func (c *GameCache) GetCachePool() *redis.Pool {
	if c.cachePool == nil {
		c.cachePool = &redis.Pool{
			MaxIdle:     50,
			IdleTimeout: 240 * time.Second,
			MaxActive:   50,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", c.Setting.URL)
				// redis.DialConnectTimeout(c.Setting.ConnectTimeout),
				// redis.DialReadTimeout(c.Setting.ReadTimeout),
				// redis.DialWriteTimeout(c.Setting.WriteTimeout))
				if err != nil {
					messagehandle.ErrorLogPrintln("newCachePool-1", c, err)
					return nil, fmt.Errorf("redis connection error: %s", err)
				}
				//验证redis密码
				// if _, authErr := c.Do("AUTH", RedisPassword); authErr != nil {
				// 	return nil, fmt.Errorf("redis auth password error: %s", authErr)
				// }
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				if err != nil {
					return fmt.Errorf("ping redis error: %s", err)
				}
				return nil
			},
		}
	}
	return c.cachePool
}

// SetAttach ...
func (c *GameCache) SetAttach(playerid string, gameIndex int64, value interface{}) {
	key := fmt.Sprintf("attach/%s/%d", playerid, gameIndex)
	if err := cacheinfo.RunSet(c.GetCachePool(), key, value, c.Setting.CacheDeleteTime); err != nil {
		fmt.Printf("SetAttach: err %s: playerid %s: gameindex %d: value %v", err, playerid, gameIndex, value)
	}
}

// GetAttach game data request
func (c *GameCache) GetAttach(playerid string, gameIndex int64) interface{} {
	err := messagehandle.New()
	key := fmt.Sprintf("attach/%s/%d", playerid, gameIndex)
	info, errMsg := cacheinfo.Get(c.GetCachePool(), key)

	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		messagehandle.ErrorLogPrintln("GetAttach-1", errMsg, key)
		return nil
	}

	return info
}
