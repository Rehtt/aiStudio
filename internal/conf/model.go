package conf

import "fmt"

type Conf struct {
	Server *Server `yaml:"server"`
}
type Server struct {
	Listen string `yaml:"listen"`
	Redis  Redis  `yaml:"redis"`
	Mysql  Mysql  `yaml:"mysql"`
	Midj   Midj   `yaml:"midj"`
}
type Redis struct {
	Addr     string `yaml:"addr"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
	Username string `yaml:"username"`
}
type Mysql struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Addr     string `yaml:"addr"`
	Port     int    `yaml:"port"`
	DB       string `yaml:"db"`
	Option   string `yaml:"option"`
}

type Midj struct {
	BotToken  string `yaml:"bot_token"`
	UserToken string `yaml:"user_token"`
	ServerID  string `yaml:"server_id"`
	ChannelID string `yaml:"channel_id"`
}

func (m Mysql) ToDSNString() string {
	if m.Option == "" {
		m.Option = "charset=utf8mb4&parseTime=True&loc=Local"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.Username, m.Password, m.Addr, m.Port, m.DB, m.Option)
}
