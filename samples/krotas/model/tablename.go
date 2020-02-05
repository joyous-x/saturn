package model

import (
	"fmt"
)

var tnUserInfo = func(appname string) string {
	return fmt.Sprintf("t_%v_user_info", appname)
}

var tnAttendance = func(appname string) string {
	return fmt.Sprintf("t_%v_attendance", appname)
}

var tnUserRelation = func(appname string) string {
	return fmt.Sprintf("t_%v_user_relation", appname)
}

var tnBalloonDatas = func(appname string) string {
	return fmt.Sprintf("t_%v_balloon_datas", appname)
}
