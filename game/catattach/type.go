package catattach

import (
	"encoding/json"
	"fmt"

	"github.com/YWJSonic/GameServer/catpunch/game/cache"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/attach"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/foundation"
)

// UserAttach ...
type UserAttach struct {
	cache   *cache.GameCache
	userID  int64
	dataMap map[int64]map[int64]*attach.Info
}

// LoadData ...
func (us *UserAttach) LoadData() {
	// test Data
	if us.dataMap == nil {
		us.dataMap = make(map[int64]map[int64]*attach.Info)
	}
	// TODO get cache data insert map
	attStr := us.cache.GetAttach(us.userID)
	if attStr == nil {
		fmt.Println("UserAttach LoadData error", attStr)
		return
	}
	if errMsg := json.Unmarshal([]byte(attStr.(string)), &us.dataMap); errMsg != nil {
		fmt.Println("UserAttach LoadData", errMsg)
	}
}

// Save ...
func (us *UserAttach) Save() {
	us.cache.SetAttach(us.userID, foundation.JSONToString(us.dataMap))
}

// Get ...
func (us *UserAttach) Get(attachkind int64, attachtype int64) *attach.Info {
	if _, ok := (*us.GetType(attachkind))[attachtype]; !ok {
		us.SetValue(attachkind, attachtype, "", 0)
	}
	return us.dataMap[attachkind][attachtype]
}

// GetType ...
func (us *UserAttach) GetType(attachkind int64) *map[int64]*attach.Info {
	if _, ok := us.dataMap[attachkind]; !ok {
		us.dataMap[attachkind] = make(map[int64]*attach.Info)
	}
	result := us.dataMap[attachkind]
	return &result
}

// SetDBValue ...
func (us *UserAttach) SetDBValue(attachKind, attachType int64, SValue string, IValue int64) {
}

// SetValue ...
func (us *UserAttach) SetValue(attachKind, attachType int64, SValue string, IValue int64) {

	if att, ok := (*us.GetType(attachKind))[attachType]; !ok {
		att = attach.NewInfo(attachKind, attachType, false)
		att.SetSValue(SValue)
		att.SetIValue(IValue)
		us.dataMap[attachKind][attachType] = att
	} else {
		att.SetSValue(SValue)
		att.SetIValue(IValue)
	}
}

// SetAttach ...
func (us *UserAttach) SetAttach(info *attach.Info) {
	us.dataMap[info.GetKind()][info.GetTypes()] = info
}
