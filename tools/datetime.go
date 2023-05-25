package tools

import (
	"fmt"
	"strconv"
	"time"
)

// 今天距离1970.1.1多少天
func TodayTo1970() string {
	nowtime := time.Now().Unix() + 28800
	return strconv.FormatInt(nowtime/(24*60*60), 10)
}

// 今天距离1970.1.1多少天
func TodayTo1970Int64() int64 {
	nowtime := time.Now().Unix() + 28800
	return nowtime / (24 * 60 * 60)
}

// 某天距离1970.1.1多少天
func OnedayTo1970(t int64) int64 {
	return (t + 28800) / (24 * 60 * 60)
}

func TimeDifferenceByDays(NewTime int64, OldTime int64) int {
	//timeLayout := "2006-01-02 15:04:05"

	NewDateTime := time.Unix(NewTime, 0).Format("2006-01-02")
	OldDateTime := time.Unix(OldTime, 0).Format("2006-01-02")

	New, _ := time.ParseInLocation("2006-01-02", NewDateTime, time.Local)
	Old, _ := time.ParseInLocation("2006-01-02", OldDateTime, time.Local)

	second := New.Unix() - Old.Unix()

	return int(second / (60 * 60 * 24))
}

func TimeDifferenceBySecond(NowTime time.Time, OldTime time.Time) int64 {
	return NowTime.Unix() - OldTime.Unix()
}

func TimeDifferenceByDays2(NewTime int64, OldTime int64) int {
	//timeLayout := "2006-01-02 15:04:05"

	NewDateTime := time.Unix(NewTime, 0)
	OldDateTime := time.Unix(OldTime, 0)

	New := time.Date(NewDateTime.Year(), NewDateTime.Month(), NewDateTime.Day(), 0, 0, 0, 0, NewDateTime.Location())
	Old := time.Date(OldDateTime.Year(), OldDateTime.Month(), OldDateTime.Day(), 0, 0, 0, 0, OldDateTime.Location())

	second := New.Unix() - Old.Unix()

	return int(second / (60 * 60 * 24))
}

func TimestampToDateTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func TimeToMysqlDataTimeString(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func StringTimeToTime(timeStr, layout string) time.Time {
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}

func TimestampToTime(timestampStr string) time.Time {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(timestamp, 0)
}

func TimesToYearMonth(t time.Time) string {
	return t.Format("200601")
}

func DateTimeToTimestamp(dataTime string) int64 {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", dataTime, time.Local)
	return t.Unix()
}

func DateTimeToTimestamp1(dataTime string) int64 {
	t, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", dataTime, time.Local)
	return t.Unix()
}

func TimeToDateTimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func NowDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func NowDate() string {
	return time.Now().Format("2006-01-02")
}

// NowDateAdd 计算跟当前时间相隔n秒的时间,负数位之前, 正数为之后, 返回年月日
func NowDateAdd(seconds int) string {
	h := fmt.Sprintf("%vs", seconds)
	d, _ := time.ParseDuration(h)
	return time.Now().Add(d).Format("2006-01-02")
}

func NowDateTimeAdd(seconds int) string {
	h := fmt.Sprintf("%vs", seconds)
	d, _ := time.ParseDuration(h)
	return time.Now().Add(d).Format("2006-01-02 15:04:05")
}

// 获取当天0点时间戳
func TodayZeroTimestamp() int64 {
	nowTime := time.Now()
	return time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, nowTime.Location()).Unix()
}

// 获取明天0点时间戳
func TomorrowZeroTimestamp() time.Time {
	nowTime := time.Now()
	return time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day()+1, 0, 0, 0, 0, nowTime.Location())
}

func GetTodayLeftSecond() int64 {
	now := time.Now()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return TimeDifferenceBySecond(endTime, now)
}
