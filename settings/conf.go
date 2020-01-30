package settings

import (
	"os"

	"gopkg.in/yaml.v2"
)

/*	read file of configuartion file and return readed data or error in case no such file found */
func readConf(filePath string) (data []byte, err error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return
	}
	stat, err := file.Stat()
	if err != nil {
		return
	}
	data = make([]byte, stat.Size())
	_, err = file.Read(data)
	return
}

/* example of conf file should be like:
server:
	host: localhost
	port: 80
	rootdir: directoryPath/
*/

/* parsing conf file and return it in settings struct*/
func unmarshalConf(filePath string, values *Values) error {
	var m map[interface{}]interface{}
	var server []byte

	fileData, err := readConf(filePath)
	if err != nil {
		return err //the file path given does not exist or don't have the permissions of read
	}
	m = make(map[interface{}]interface{})
	err = yaml.Unmarshal(fileData, &m)
	if err != nil {
		return err //file given is not a yaml file
	}
	server, err = yaml.Marshal(m["server"])
	if err != nil {
		return err //server field not found in conf file
	}
	return yaml.Unmarshal(server, values)
}
