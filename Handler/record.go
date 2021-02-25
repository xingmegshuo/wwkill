/***************************
@File        : record.go
@Time        : 2021/01/12 15:40:10
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 战绩信息
****************************/

package Handler

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"wwKill/Mydb"
)

// 获取最近战绩
func GetRecord(mes []byte) string {
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "获取战绩失败,数据无法解析")
	}
	User := Mydb.User{
		OpenID: user.OpenID,
	}
	thisUser, has := ctrlUser.GetUser(User)
	if has {
		back := Mydb.Record{
			User: int(thisUser.Id),
		}
		backs := ctrlRecord.GetUser(back)
		return RecordToString("ok", backs, "获取战绩成功")
	} else {
		return ToMes("error", "获取战绩失败，找不到用户")
	}
}

// 获取战绩信息
func GetRecordAll(mes []byte) string {
	err := json.Unmarshal(mes, &user)
	log.Println("获取战绩信息func")

	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "获取战绩失败,数据无法解析")
	}
	User := Mydb.User{
		OpenID: user.OpenID,
	}
	thisUser, has := ctrlUser.GetUser(User)
	if has {
		back := Mydb.Record{
			User: int(thisUser.Id),
		}
		backs := ctrlRecord.GetUser(back)
		return AllRecordToString("ok", backs, "获取战绩成功")
	} else {
		return ToMes("error", "获取战绩失败，找不到用户")
	}
}

// record to str
func RecordToString(status string, records []Mydb.Record, mes string) string {
	str := "{'status':'" + status + "','mes':'" + mes + "','data':["
	for l, item := range records {
		if len(item.GameMode) > 0 && l < 10 {
			if l == len(records)-1 {
				str = str + "{'gameType':'" + item.GameMode + "','runAway':'" + strconv.Itoa(item.RunAway) + "','gameTime':'" + item.GameTime + "','identity':'" + item.Identity + "','result':'" + item.Result + "'}"
			} else {
				str = str + "{'gameType':'" + item.GameMode + "','runAway':'" + strconv.Itoa(item.RunAway) + "','gameTime':'" + item.GameTime + "','identity':'" + item.Identity + "','result':'" + item.Result + "'},"
			}
		}
	}
	str = str + "]}"
	str = strings.Replace(str, "'", "\"", -1)
	return str
}

// allRecord to str
func AllRecordToString(status string, records []Mydb.Record, mes string) string {
	str := "{'status':'" + status + "','mes':'" + mes + "','data':{"
	count := 0
	winCount := 0
	runAway := 0
	maxWin := 0
	win := 0
	for _, item := range records {
		if len(item.GameMode) > 0 {
			if item.Result == "胜利" {
				winCount++
				win++
			} else {
				win = 0
			}
			if item.RunAway == 1 {
				runAway++
			}

			count++
		}
		if win > maxWin {
			maxWin = win
		}
	}
	if winCount > 0 {
		str = str + "'count':'" + strconv.Itoa(count) + "','runAway':'" + strconv.Itoa(runAway) + "','maxWin':'" + strconv.Itoa(maxWin) + "','winRate':'" + strconv.Itoa(winCount/count*100) + "'}}"
	} else {
		str = str + "'count':'" + strconv.Itoa(count) + "','runAway':'" + strconv.Itoa(runAway) + "','maxWin':'" + strconv.Itoa(maxWin) + "','winRate':'" + strconv.Itoa(0) + "'}}"
	}
	str = strings.Replace(str, "'", "\"", -1)
	return str
}
