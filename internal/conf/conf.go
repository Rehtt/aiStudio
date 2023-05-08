package conf

import (
	"github.com/Rehtt/Kit/i18n"
	"github.com/Rehtt/Kit/log/logs"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var Config = new(Conf)

func Init(path string) error {
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) != ".yaml" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(data, Config)
	})
}

func GetServer() *Server {
	return Config.Server
}

func GenConfig(path string) error {
	if !strings.HasSuffix(path, "yaml") {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		path = filepath.Join(path, "config.yaml")
	}
	var tmp Conf
	tmp.Server = new(Server)
	data, _ := yaml.Marshal(tmp)
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	logs.Info(i18n.GetText("模板生成完毕，路径：%s"), path)
	return nil
}
