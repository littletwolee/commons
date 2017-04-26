package commons

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var (
	ConsConfig *config
	confList   []map[string]map[string]string
	filePath   string
)

type config struct{}

func init() {
	ConsConfig = new(config)
	filePath = "conf/app.conf"
	readList()
}

// @Title getValue
// @Description get value from config by section & name
// @Parameters
//            section            string          group name
//            name               string          node name
// @Returns err:error
func (c *config) getValue(section, name string) string {
	// if conflist == nil{
	// 	c.ReadList()
	// }
	for _, v := range confList {
		for key, value := range v {
			if key == section {
				return value[name]
			}
		}
	}
	return "no value"
}

// @Title readList
// @Description create a kv list from config
// @Returns config list:[]map[string]map[string]string
func readList() {
	file, err := os.Open(filePath)
	if err != nil {
		ConsLogger.LogErr(err)
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
				ConsLogger.LogErr(err)
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
			if uniquappend(section) == true {
				confList = append(confList, data)
			}
		}

	}
}

// @Title uniquappend
// @Description check section is unique
// @Parameters
//            conf            string          section name
// @Returns result:bool
func uniquappend(conf string) bool {
	for _, v := range confList {
		for k, _ := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}
