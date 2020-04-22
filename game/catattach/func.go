package catattach

import (
	"github.com/YWJSonic/ServerUtility/attach"
)

// NewAttach ...
func NewAttach(attSetting Setting) attach.IAttach {
	attach := &UserAttach{
		userIDStr: attSetting.UserIDStr,
		kind:      attSetting.Kind,
		db:        attSetting.DB,
		dataMap:   make(map[int64]map[int64]*attach.Info),
		redis:     attSetting.Redis,
	}
	// attach.InitData(userID)
	return attach
}
