package cache

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/code"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/messagehandle"
)

// Setting ...
type Setting struct {
	ConnectTimeout, ReadTimeout, WriteTimeout, CacheDeleteTime time.Duration
	URL                                                        string
}

// GameCache ICache
type GameCache struct {
	cachePool *redis.Pool
	cacheMap  map[string]interface{}
	Setting   Setting
}

// SetAttach ...
func (c *GameCache) SetAttach(playerid int64, value interface{}) {
	key := fmt.Sprintf("attach%d", playerid)

	c.cacheMap[key] = value
	// cacheinfo.RunSet(c.GetCachePool(), key, value, c.Setting.CacheDeleteTime)
}

// GetAttach game data request
func (c *GameCache) GetAttach(playerid int64) interface{} {
	err := messagehandle.New()
	key := fmt.Sprintf("attach%d", playerid)
	// info, errMsg := cacheinfo.Get(c.GetCachePool(), key)
	info, exist := c.cacheMap[key]
	if !exist {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintf("cache key `%s` not exist", key)
		messagehandle.ErrorLogPrintln("GetAttach-1", key)
		return nil
	}

	return info
}
