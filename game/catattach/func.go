package catattach

import (
	"gitlab.fbk168.com/gamedevjp/cat/server/game/cache"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/attach"
)

// NewUserAttach ...
func NewUserAttach(cache *cache.GameCache, userID int64) *UserAttach {
	attach := &UserAttach{
		cache:   cache,
		userID:  userID,
		dataMap: make(map[int64]map[int64]*attach.Info),
	}
	// attach.InitData(userID)
	return attach
}
