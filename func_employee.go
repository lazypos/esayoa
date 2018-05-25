package esayoa

import (
	"log"
)

func EmployeeList(token int) string {
	type STEmpInfo struct {
		Eid string
	}

	type STEmpList struct {
		Del         bool
		ArrEmployee []*STEmpInfo
	}
	rows, err := GDBpgopt.Query(SQL_LOAD_USERINFO)
	if err != nil {
		log.Println("加载用户信息失败", err)
		return ""
	}

	for rows.Next() {

	}

	return ""
}

func FrameEmployee(sInfo *SessionInfo) interface{} {
	type STFrmEmp struct {
		Token  int
		Counts string
	}

	counts, err := GDBpgopt.QueryVal(SQL_COUNT_USERINFO)
	if err != nil {
		log.Println("加载用户信息总数失败", err)
		return nil
	}

	st := &STFrmEmp{Token: sInfo.UserInfo.Auth, Counts: counts}
	return st
}
