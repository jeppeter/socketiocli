package socketiocli

import (
	"encoding/json"
	"fmt"
	logging "github.com/jeppeter/go-logging"
	"reflect"
	"strings"
)

const (
	SOCKET_IO_SID           = "sid"
	SOCKET_IO_PROTOCOL      = "upgrades"
	SOCKET_IO_HEART_TIMEOUT = "pingInterval"
	SOCKET_IO_TIMEOUT       = "pingTimeout"
)

func getJsonValue(path string, v map[string]interface{}) (val string, err error) {
	var pathext []string
	var curmap map[string]interface{}

	val = ""
	err = nil
	pathext = strings.Split(path, "/")
	if len(pathext) == 0 {
		err = fmt.Errorf("can not split (%s)", path)
		return
	}

	curmap = v
	for i, curpath := range pathext {
		curval, ok := curmap[curpath]
		if !ok {
			err = fmt.Errorf("can not find (%s) in %s", curpath, path)
			return
		}

		if i == (len(pathext) - 1) {
			switch reflect.TypeOf(curval).Kind() {
			case reflect.Int:
				val = fmt.Sprintf("%d", curval)
			case reflect.Uint32:
				val = fmt.Sprintf("%d", curval)
			case reflect.Uint64:
				val = fmt.Sprintf("%d", curval)
			case reflect.Float64:
				val = fmt.Sprintf("%f", curval)
			case reflect.Float32:
				val = fmt.Sprintf("%f", curval)
			case reflect.String:
				val = fmt.Sprintf("%s", curval)
			case reflect.Slice:
				val = ""
				for i := 0; i < reflect.ValueOf(curval).Len(); i++ {
					a := reflect.ValueOf(curval).Index(i)
					if i != 0 {
						val = fmt.Sprintf("%s,", val)
					}
					val = fmt.Sprintf("%s%s", val, a)
				}
			default:
				logging.Debugf("type %s", reflect.TypeOf(curval).Kind())
				val = fmt.Sprintf("%q", curval)
				logging.Debugf("%s => %s", path, val)
			}
			err = nil
			return
		}

		cval, ok := curval.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("can not parse in (%s) for path(%s)", curpath, path)
			return
		}

		curmap = cval
	}

	err = fmt.Errorf("can not find (%s) all over", path)
	return
}

func getJsonValueDefault(instr string, key string, defval string) string {
	var v map[string]interface{}
	err := json.Unmarshal([]byte(instr), &v)
	if err != nil {
		return defval
	}

	val, err := getJsonValue(key, v)
	if err != nil {
		return defval
	}
	return val
}
