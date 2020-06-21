package util

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type ToolConfig struct {
	ListenAddress  string
	ListenPort     int
	LocalResources bool
	//LogPath
	//LogLevel
	DataDirectory string
	TempDirectory string
	ForceRescan   bool
	//CacheSize

}

func ReadConfig(filename string) (*ToolConfig, error) {
	//FIXME all printf shoudl be logs
	ret := &ToolConfig{
		ListenAddress:  "localhost",
		ListenPort:     8666,
		LocalResources: false,
		DataDirectory:  ".",
		TempDirectory:  "/tmp/",
		ForceRescan:    false,
	}
	of := reflect.ValueOf(ret).Elem()
	file, e := os.Open(filename)
	if e != nil {
		return ret, e
	}

	//FIXME check if the TempFir exists

	scanner := bufio.NewScanner(file)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue // Ignore comments
		} else if strings.Contains(line, "=") {
			split := strings.SplitN(line, "=", 2)
			val := of.FieldByName(split[0])
			stringVal := split[1]
			if val.IsValid() {
				if val.CanAddr() {
					switch val.Interface().(type) {
					case int:
						iv, err := strconv.Atoi(stringVal)
						if err != nil { //FIXME error should be handled more global
							panic(err)
						}
						val.SetInt(int64(iv))
					case string:
						val.SetString(stringVal)
					case bool:
						if "true" == strings.ToLower(stringVal) {
							val.SetBool(true)
						} else {
							val.SetBool(false)
						}
					}
				} else {
					fmt.Printf("Cannot Set!! %v \n", split[0])
				}
			} else {
				fmt.Printf("Not a config: %s\n", split[0])
			}
		} else {
			fmt.Printf("Wrong Config in line %d: (%s)\n", lineNum, line)
			continue
		}
	}

	return ret, nil
}
