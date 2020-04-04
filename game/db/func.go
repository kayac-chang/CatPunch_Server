package db

import (
	"database/sql"
	"fmt"

	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/dbinfo"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/foundation"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/messagehandle"
)

// SetLog new goruting set log
func SetLog(db *sql.DB, account string, playerID, time int64, activityEvent uint8, iValue1, iValue2, iValue3 int64, sValue1, sValue2, sValue3, msg string) messagehandle.ErrorMsg {
	tableName := foundation.ServerNow().Format("20060102")
	query := fmt.Sprintf("INSERT INTO `%s` VALUE(NULL,\"%s\",%d,%d, %d, %d,%d,%d,\"%s\",\"%s\",\"%s\",\"%s\");", tableName, account, playerID, time, activityEvent, iValue1, iValue2, iValue3, sValue1, sValue2, sValue3, msg)
	_, err := dbinfo.CallWrite(db, query)
	return err
}
