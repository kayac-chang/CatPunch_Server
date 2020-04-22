package catattach

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	attach "github.com/YWJSonic/ServerUtility/attach"
	"github.com/YWJSonic/ServerUtility/code"
	"github.com/YWJSonic/ServerUtility/foundation"
	"gitlab.fbk168.com/gamedevjp/cat/server/game/cache"
	"gitlab.fbk168.com/gamedevjp/cat/server/game/db"
)

// Setting ...
type Setting struct {
	UserIDStr string
	Kind      int64
	DB        *sql.DB
	Redis     *cache.GameCache
}

// UserAttach ...
type UserAttach struct {
	userID    int64
	userIDStr string
	kind      int64
	db        *sql.DB
	dataMap   map[int64]map[int64]*attach.Info

	redis *cache.GameCache
}

// LoadData ...
func (us *UserAttach) LoadData() {
	// test Data
	us.dataMap = make(map[int64]map[int64]*attach.Info)
	//---
	// redis load data
	if attcache := us.redis.GetAttach(us.userIDStr, us.kind); attcache != nil {
		if errMsg := json.Unmarshal(attcache.([]byte), &us.dataMap); errMsg != nil {
			fmt.Println("errMsg ", errMsg)
		} else {
			return
		}
	}

	// if fail sql load data
	result, err := db.GetAttachKind(us.db, us.userIDStr, us.kind)
	if err.ErrorCode != code.OK {
		fmt.Println(err)
	}
	for _, row := range result {
		att := &attach.Info{
			Kind:     us.kind,
			Types:    row["Type"].(int64),
			IValue:   row["IValue"].(int64),
			IsDBData: true,
		}
		us.SetAttach(att)
	}
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

	if att, ok := (*us.GetType(attachKind))[attachType]; !ok {
		att = attach.NewInfo(attachKind, attachType, true)
		att.SetSValue(SValue)
		att.SetIValue(IValue)
		us.dataMap[attachKind][attachType] = att
	} else {
		att.SetSValue(SValue)
		att.SetIValue(IValue)
	}
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
	info.IsDirty = true
	if _, ok := us.dataMap[info.GetKind()]; !ok {
		us.dataMap[info.GetKind()] = make(map[int64]*attach.Info)
	}
	us.dataMap[info.GetKind()][info.GetTypes()] = info
}

// Save ...
func (us *UserAttach) Save() {
	// set to redis
	attStr := foundation.JSONToString(us.dataMap)
	us.redis.SetAttach(us.userIDStr, us.kind, attStr)

	// set to db
	quarys := []string{}
	for kind, typeAtt := range us.dataMap {
		for t, att := range typeAtt {
			if att.IsDirty {
				quarys = append(quarys, fmt.Sprintf(us.setQuary(), us.userID, us.userIDStr, kind, t, att.IValue, att.IValue))
			}
		}
	}

	ctx := context.TODO()
	sqlConn, err := us.db.Conn(ctx)
	defer sqlConn.Close()
	if err != nil {
		return
	}

	for _, quart := range quarys {
		_, err := sqlConn.ExecContext(ctx, quart)
		if err != nil {
			fmt.Println("att set into sql error: ", err)
		}
	}
}

// Clear ...
func (us *UserAttach) Clear() {

}

func (us *UserAttach) setQuary() string {
	return "INSERT INTO attach VALUES (%d, \"%s\", %d, %d, %d) ON DUPLICATE KEY UPDATE IValue = %d;\n"
}
