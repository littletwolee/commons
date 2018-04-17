package commons

import (
	"fmt"
	"reflect"
	"regexp"
)

var (
	consCheck        *check
	defaultCheckList string = "select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute"
)

type check struct {
	CheckList string
}

func GetCheck() *check {
	if consCheck == nil {
		consCheck = &check{
			CheckList: GetConfig().GetString("check.sqlinjection"),
		}
		if consCheck.CheckList == "" {
			consCheck.CheckList = defaultCheckList
		}
	}
	return consCheck
}

// // @Title SQLInjectionAttack
// // @Description check interface sql injection attack
// // @Parameters
// //          value          interface{}           parameter
// // @Returns value:interface{} err:error
func (m *check) SQLInjectionAttack(value interface{}) (interface{}, error) {
	if value == nil || value == "" {
		return nil, fmt.Errorf(ERROR_PARAMETER_EMPTY)
	}
	switch value.(type) {
	case string:
		return m.attackCheck(value.(string))
	default:
		object := reflect.ValueOf(value)
		for i := 0; i < object.NumField(); i++ {
			field := object.Field(i)
			_, err := m.attackCheck(field.Interface())
			if err != nil {
				return nil, err
			}
		}
	}
	return value, nil
}

func (m *check) attackCheck(value interface{}) (interface{}, error) {
	var (
		re  *regexp.Regexp
		err error
		v   string
	)
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(` + m.CheckList + `)\b)`
	re, err = regexp.Compile(str)
	if err != nil {
		value = ""
		goto RESULT
	}
	v = fmt.Sprint(value)
	if v == "" {
		return "", fmt.Errorf(ERROR_PARAMETER_EMPTY)
	}
	if re.MatchString(v) {
		err = fmt.Errorf(ERROR_PARAMETER_SQL_ATTACK, value)
		value = ""
		goto RESULT
	}
	goto RESULT
RESULT:
	return value, err
}
