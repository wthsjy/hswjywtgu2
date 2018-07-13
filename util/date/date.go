package util

import "time"

func Date2ts(dateStr string, withDefault bool) int64 {
	if dateStr == "" {
		if withDefault {
			return time.Now().Unix()
		} else {
			return 0
		}
	}
	timeLayout := "2006-01-02 15:04:05"                            //转化所需模板
	loc, _ := time.LoadLocation("Local")                           //重要：获取时区
	theTime, err := time.ParseInLocation(timeLayout, dateStr, loc) //使用模板在对应时区转化为time.time类型
	if err != nil {
		if withDefault {
			return time.Now().Unix()
		} else {
			return 0
		}
	}
	return theTime.Unix()

}
