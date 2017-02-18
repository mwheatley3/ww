package main

import (
	"github.com/mwheatley3/ww/server/config"
	goconfig "github.com/mwheatley3/ww/server/config/config"
	"github.com/mwheatley3/ww/server/log"
	"github.com/mwheatley3/ww/server/personal/web"
	"github.com/mwheatley3/ww/server/pg"
)

const (
	// DefaultConfPath default configuration file path
	DefaultConfPath = "conf/personal.conf"
)

var confPath string

// Config is portal service configuration
type Config struct {
	Port string
	Web  struct {
		web.Config
		Cookie struct {
			BlockKey goconfig.HexBytes
			HashKey  goconfig.HexBytes
		}
	}
	Postgres pg.Config
	Log      log.Config
}

func loadConfig() Config {
	var c Config
	config.MustLoad(&c, config.FromFileWithOverride(confPath))

	// c.Web.Config.Cookie.BlockKey = c.Web.Cookie.BlockKey
	// c.Web.Config.Cookie.HashKey = c.Web.Cookie.HashKey

	return c
}

func pgLoadConfig() (pg.Config, log.Config) {
	c := loadConfig()
	return c.Postgres, c.Log
}
