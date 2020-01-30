package settings

import (
	"os"
	"path/filepath"
)

//Values given in conf file or by flags
type Values struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Root string `yaml:"rootdir"`
}

//GetInfo erver
func GetInfo() error {
	var vals Values

	if os.Getppid() == 1 { //already a daemon
		return nil
	}
	//get filePath from flags
	err := unmarshalConf("./server_conf.yaml", &vals)
	if err != nil {
		return err
	}
	vals.Store()
	return nil
}

//Store configuration setting in environement vars
func (v Values) Store() {
	os.Setenv("WEB_SERVER_HOST", v.Host)
	os.Setenv("WEB_SERVER_PORT", v.Port)
	root, err := filepath.Abs(v.Root) // get the absolute path of given directory
	if err == nil {
		os.Setenv("WEB_SERVER_ROOT", root)
	}
}

//Restore setting values from environement
func Restore() (v Values) {
	return Values{
		Host: os.Getenv("WEB_SERVER_HOST"),
		Port: os.Getenv("WEB_SERVER_PORT"),
		Root: os.Getenv("WEB_SERVER_ROOT"),
	}
}
