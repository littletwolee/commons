package commons

import (
	"bytes"
	"fmt"
	"sort"
	"time"
)

var (
	constimes *times
)

type times struct{}

func GetTimes() *times {
	if constimes == nil {
		constimes = &times{}
	}
	return constimes
}

// @Title TimeToStamp
// @Description get stamp from time
// @Parameters
//       inputTime         time.Time       time
// @Returns stamp:int64
func (t *times) TimeToStamp(inputTime time.Time) int64 {
	return inputTime.Unix()
}

// @Title StampToTime
// @Description get time from stamp
// @Parameters
//       stamp         int64       stamp
// @Returns time:time.Time
func (t *times) StampToTime(stamp int64) time.Time {
	return time.Unix(stamp, 0)
}

// @Title StringToTime
// @Description get time from a time of string type
// @Parameters
//       strTime         string        time string
// @Returns time:*time.Time
func (t *times) StringToTime(strTime string) (*time.Time, error) {
	timeLayout := "2006-01-02 15:04:05"
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, err
	}
	theTime, err := time.ParseInLocation(timeLayout, strTime, loc)
	if err != nil {
		return nil, err
	}
	return &theTime, nil
}

// @Title TimeToString
// @Description time to string
// @Parameters
//       inputTime         time.Time        time
// @Returns time:*time.Time
func (t *times) TimeToString(inputTime time.Time) string {
	return inputTime.Format("2006-01-02 15:04:05")
}

type timeFmtRule struct {
	TimeType        int
	TimeDescription string
}
type timeFmtRules []timeFmtRule

func (u timeFmtRules) Len() int {
	return len(u)
}

func (u timeFmtRules) Less(i, j int) bool {
	if u[i].TimeType == u[j].TimeType {
		return u[i].TimeDescription > u[j].TimeDescription
	}
	return u[i].TimeType < u[j].TimeType
}

func (u timeFmtRules) Swap(i, j int) {
	u[i].TimeType, u[j].TimeType = u[j].TimeType, u[i].TimeType
	u[i].TimeDescription, u[j].TimeDescription = u[j].TimeDescription, u[i].TimeDescription
}

// @Title TimeFmt
// @Description time format from rule
// @Parameters
//       startTime         time.Time        start time
//       endTime           time.Time        end time
//       rules             map[int]string   format rule
// @Returns time:*time.Time
func (t *times) TimeFmt(startTime, endTime time.Time, rules map[int]string) string {
	var (
		str       *bytes.Buffer
		sortRules []timeFmtRule
		diff      int
	)
	str = &bytes.Buffer{}
	sortRules = []timeFmtRule{}
	if !startTime.Before(endTime) {
		startTime, endTime = endTime, startTime
	}
	for k, v := range rules {
		sortRules = append(sortRules, timeFmtRule{
			TimeType:        k,
			TimeDescription: v,
		})
	}
	sort.Sort(timeFmtRules(sortRules))
	diff = int(endTime.Sub(startTime).Seconds())
	for _, v := range sortRules {
		switch v.TimeType {
		case YEAR:
			rule := 365 * 24 * 60 * 60
			y := diff / rule
			str = addTimeStr(y, v.TimeDescription, str)
			diff -= y * rule
		case MONTH:
			rule := 30 * 24 * 60 * 60
			m := diff / rule
			str = addTimeStr(m, v.TimeDescription, str)
			diff -= m * rule
		case DAY:
			rule := 24 * 60 * 60
			d := diff / rule
			str = addTimeStr(d, v.TimeDescription, str)
			diff -= d * rule
		case HOUR:
			rule := 60 * 60
			h := diff / rule
			str = addTimeStr(h, v.TimeDescription, str)
			diff -= h * rule
		case MINUTE:
			rule := 60
			m := diff / rule
			str = addTimeStr(m, v.TimeDescription, str)
			diff -= m * rule
		case SECOND:
			str = addTimeStr(diff, v.TimeDescription, str)
		default:
			continue
		}
	}
	return str.String()
}
func addTimeStr(diff int, timeType string, str *bytes.Buffer) *bytes.Buffer {
	if diff > 0 {
		str.WriteString(fmt.Sprintf("%d%s", diff, timeType))
	}
	return str
}
