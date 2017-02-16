package pg_test

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/conversable/woodhouse/server/pg"
)

func getStr(str string, def string) string {
	if s := os.Getenv(str); s != "" {
		return s
	}

	return def
}

func getInt(str string, def int) int {
	if s := getStr(str, ""); s != "" {
		v, _ := strconv.Atoi(s)
		return v
	}

	return def
}

// ConfigToEnv returns a set of env vars with the pg env
// vars set
func ConfigToEnv(prefix string, conf pg.Config) []string {
	return []string{
		prefix + "HOST=" + conf.Host,
		prefix + "PORT=" + strconv.Itoa(int(conf.Port)),
		prefix + "USER=" + conf.User,
		prefix + "PASSWORD=" + conf.Password,
		prefix + "DATABASE=" + conf.Database,
		prefix + "SSL_MODE=" + conf.SslMode,
	}
}

// CreateTestDB returns a pg.DB with values pulled from the environment
func CreateTestDB(envPrefix string) *pg.Db {
	if envPrefix == "" {
		envPrefix = "PG_TEST_"
	}

	var c = pg.Config{
		Host:          getStr(envPrefix+"HOST", "localhost"),
		Port:          uint16(getInt(envPrefix+"PORT", 5432)),
		User:          getStr(envPrefix+"USER", ""),
		Password:      getStr(envPrefix+"PASSWORD", ""),
		Database:      getStr(envPrefix+"DATABASE", ""),
		SslMode:       getStr(envPrefix+"SSL_MODE", ""),
		SlowThreshold: 1000,
	}

	return pg.NewDb(logrus.New(), c)
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// CreateRandomDB creates a random database and returns a drop
// function to remove it
func CreateRandomDB(envPrefix string) (*pg.Db, func() error) {
	if envPrefix == "" {
		envPrefix = "PG_TEST_"
	}

	var c = pg.Config{
		Host:          getStr(envPrefix+"HOST", "localhost"),
		Port:          uint16(getInt(envPrefix+"PORT", 5432)),
		User:          getStr(envPrefix+"USER", ""),
		Password:      getStr(envPrefix+"PASSWORD", ""),
		Database:      getStr(envPrefix+"DATABASE", "postgres"),
		SslMode:       getStr(envPrefix+"SSL_MODE", ""),
		SlowThreshold: 1000,
	}

	db := pg.NewDb(logrus.New(), c)

	if err := db.Connect(); err != nil {
		panic(err)
	}

	dbName := fmt.Sprintf("test_db_%d", r.Int())

	_, err := db.Exec(fmt.Sprintf(`create database "%s"`, dbName), nil)

	if err != nil {
		panic(err)
	}

	c2 := c
	c2.Database = dbName

	db2 := pg.NewDb(logrus.New(), c2)

	return db2, func() error {
		if err := db2.Close(); err != nil {
			return err
		}

		_, err := db.Exec(fmt.Sprintf(`drop database "%s"`, dbName), nil)
		return err
	}
}
