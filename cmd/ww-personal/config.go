package main

import (
	"github.com/conversable/woodhouse/server/config"
	goconfig "github.com/conversable/woodhouse/server/config/config"
	"github.com/conversable/woodhouse/server/log"
	"github.com/conversable/woodhouse/server/metrics"
	"github.com/conversable/woodhouse/server/pg"
	"github.com/conversable/woodhouse/server/portal/web"
	"github.com/conversable/woodhouse/server/vault"
)

const (
	// DefaultConfPath default configuration file path
	DefaultConfPath = "conf/portal.conf"
)

var confPath string

// Config is portal service configuration
type Config struct {
	Web struct {
		web.Config
		Cookie struct {
			BlockKey goconfig.HexBytes
			HashKey  goconfig.HexBytes
		}
	}
	Postgres pg.Config
	Log      log.Config
	Librato  metrics.LibratoConfig
	Stats    metrics.StatsConfig
	Vault    vault.Config
}

func loadConfig() Config {
	var c Config
	config.MustLoad(&c, config.FromFileWithOverride(confPath), config.FromVault(&c.Vault))

	c.Web.Config.Cookie.BlockKey = c.Web.Cookie.BlockKey
	c.Web.Config.Cookie.HashKey = c.Web.Cookie.HashKey

	return c
}

func pgLoadConfig() (pg.Config, log.Config) {
	c := loadConfig()
	return c.Postgres, c.Log
}
