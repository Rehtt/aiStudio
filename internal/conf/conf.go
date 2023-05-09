package conf

import (
	"fmt"
	"github.com/Rehtt/Kit/i18n"
	"github.com/Rehtt/Kit/log/logs"
	kit_yaml "github.com/Rehtt/Kit/yaml"
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

// GenConfig 生成配置模板
func GenConfig(path string) error {
	if !strings.HasSuffix(path, "yaml") {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		path = filepath.Join(path, "config.yaml")
	}

	fileData, _ := os.ReadFile(path)
	var tmp Conf
	if len(fileData) != 0 {
		var v bool
		for !v {
			var f string
			fmt.Printf("配置文件 %s 已存在，是否覆盖？[Y/n]:", path)
			fmt.Scanln(&f)
			switch strings.ToLower(f) {
			case "y", "":
				yaml.Unmarshal(fileData, &tmp)
				v = true
			case "n":
				os.Exit(0)
			default:
				fmt.Println("输入错误")
			}
		}
	}
	data, _ := kit_yaml.GenYamlTemplate(tmp)
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	logs.Info(i18n.GetText("模板生成完毕，路径：%s"), path)
	return nil
}
