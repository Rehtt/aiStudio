package conf

import "fmt"

type Conf struct {
	Server *Server `yaml:"server" comment:"服务器配置"`
}
type Server struct {
	// todo 临时鉴权
	AdminKey string `yaml:"admin_key" comment:"key"`
	Listen   string `yaml:"listen" comment:"监听地址"`
	Redis    Redis  `yaml:"redis"`
	Mysql    Mysql  `yaml:"mysql"`
	Midj     []Midj `yaml:"midj"`
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
	BotToken  string   `yaml:"bot_token" comment:"机器人token"`
	UserToken string   `yaml:"user_token" comment:"用户token"`
	UserID    string   `yaml:"user_id" comment:"用户id"`
	GuildID   string   `yaml:"guild_id" comment:"服务器id"`
	ChannelID []string `yaml:"channel_id" comment:"频道id"`
}

func (m Mysql) ToDSNString() string {
	if m.Option == "" {
		m.Option = "charset=utf8mb4&parseTime=True&loc=Local"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.Username, m.Password, m.Addr, m.Port, m.DB, m.Option)
}
