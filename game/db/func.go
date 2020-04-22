package db

import (
	"database/sql"
	"fmt"

	"github.com/YWJSonic/ServerUtility/dbinfo"
	"github.com/YWJSonic/ServerUtility/dbservice"
	"github.com/YWJSonic/ServerUtility/foundation"
	"github.com/YWJSonic/ServerUtility/messagehandle"
)

// GetAttachKind get db attach kind
func GetAttachKind(db *sql.DB, playeridstr string, kind int64) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	// result, err := dbservice.CallReadOutMap(db, "AttachKindGet_Read", playerid, kind)
	result, err := dbservice.QueryOutMap(db, "select * from attach where PlayerIDStr = ? AND Kind = ?;", playeridstr, kind)
	return result, err
}

// GetAttachType ...
func GetAttachType(db *sql.DB, playeridstr string, kind int64, attType int64) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	// result, err := dbservice.CallReadOutMap(db, "AttachTypeGet_Read", playerid, kind, attType)
	result, err := dbservice.QueryOutMap(db, "seletct * form attach where PlayerIDStr = ? AND Kind = ? AND Type = ?;", playeridstr, kind, attType)

	return result, err
}

// NewAttach ...
func NewAttach(db *sql.DB, args ...interface{}) (sql.Result, messagehandle.ErrorMsg) {
	result, err := dbservice.CallWrite(
		db,
		dbservice.MakeProcedureQueryStr("AttachNew_Write", len(args)),
		args...,
	)
	return result, err
}

// UpdateAttach ...
func UpdateAttach(db *sql.DB, args ...interface{}) messagehandle.ErrorMsg {
	_, err := dbservice.CallWrite(db, dbservice.MakeProcedureQueryStr("AttachSet_Write", len(args)), args...)
	return err
}

// SetLog new goruting set log
func SetLog(db *sql.DB, account string, playerID, time int64, activityEvent uint8, iValue1, iValue2, iValue3 int64, sValue1, sValue2, sValue3, msg string) messagehandle.ErrorMsg {
	tableName := foundation.ServerNow().Format("20060102")
	query := fmt.Sprintf("INSERT INTO `%s` VALUE(NULL,\"%s\",%d,%d, %d, %d,%d,%d,\"%s\",\"%s\",\"%s\",\"%s\");", tableName, account, playerID, time, activityEvent, iValue1, iValue2, iValue3, sValue1, sValue2, sValue3, msg)
	_, err := dbinfo.CallWrite(db, query)
	return err
}
