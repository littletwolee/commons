package commons

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ConsConfigHelper *ConfigHelper
	conflist         []map[string]map[string]string
)

type ConfigHelper struct {
	Filepath string
}

func init() {
	ConsConfigHelper = new(ConfigHelper)
	ConsConfigHelper.Filepath = "conf/app.conf"
	ConsConfigHelper.ReadList()
}

// @Title GetValue
// @Description get value from config by section & name
// @Parameters
//            section            string          group name
//            name               string          node name
// @Returns err:error
func (c *ConfigHelper) GetValue(section, name string) string {
	// if conflist == nil{
	// 	c.ReadList()
	// }
	for _, v := range conflist {
		for key, value := range v {
			if key == section {
				return value[name]
			}
		}
	}
	return "no value"
}

// @Title ReadList
// @Description create a kv list from config
// @Returns config list:[]map[string]map[string]string
func (c *ConfigHelper) ReadList() []map[string]map[string]string {
	file, err := os.Open(c.Filepath)
	if err != nil {
		CheckErr(err)
	}
	defer file.Close()
	var data map[string]map[string]string
	var section string
	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				CheckErr(err)
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1 : len(line)])
			data[section][strings.TrimSpace(line[0:i])] = value
			if c.uniquappend(section) == true {
				conflist = append(conflist, data)
			}
		}

	}
	return conflist
}

// @Title CheckErr
// @Description check error
// @Parameters
//            err            error          error
// @Returns err string:error
func CheckErr(err error) string {
	if err != nil {
		return fmt.Sprintf("Error is :'%s'", err.Error())
	}
	return "Notfound this error"
}

// @Title uniquappend
// @Description check section is unique
// @Parameters
//            conf            string          section name
// @Returns result:bool
func (c *ConfigHelper) uniquappend(conf string) bool {
	for _, v := range conflist {
		for k, _ := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}
