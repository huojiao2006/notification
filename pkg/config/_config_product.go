// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/koding/multiconfig"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/constants"
)

type Config struct {
	Log  LogConfig
	Grpc GrpcConfig

	Mysql struct {
		Host     string `default:"notification-db"`
		Port     int    `default:"3306"`
		User     string `default:"root"`
		Password string `default:"password"`
		Database string `default:"notification"`
		Disable  bool   `default:"false"`
		LogMode  bool   `default:"false"`
	}

	Queue struct {
		Addr string `default:"redis://notification-redis:6379"`
		Type string `default:"redis"`
	}

	Email struct {
		Protocol      string `default:"SMTP"`
		EmailHost     string `default:"mail.app-center.cn"`
		Port          int    `default:"25"`
		DisplaySender string `default:"admin_openpitrix"`
		Email         string `default:"openpitrix@app-center.cn"`
		Password      string `default:"openpitrix"`
		SSLEnable     bool   `default:"false"`
		FromEmailAddr string `default:"openpitrix@app-center.cn"`
	}

	App struct {
		Host string `default:"notification-manager"`
		Port int    `default:"9201"`

		ApiHost string `default:"notification-manager"`
		ApiPort int    `default:"9200"`

		MaxWorkingNotifications int `default:"15"`
		MaxWorkingTasks         int `default:"15"`
		MaxTaskRetryTimes       int `default:"1"`
	}

	Websocket struct {
		//Service string `default:"op,ks"`
		Service string `default:"none"`
	}
}

var instance *Config

var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
		instance.LoadConf()
	})

	return instance
}

type LogConfig struct {
	//Level string `default:"error"` // debug, info, warn, error, fatal
	Level string `default:"info"`
}

type GrpcConfig struct {
	ShowErrorCause bool `default:"false"` // show grpc error cause to frontend
}

func (c *Config) PrintUsage() {
	fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprint(os.Stdout, "\nSupported environment variables:\n")
	e := newLoader(constants.ServiceName)
	e.PrintEnvs(new(Config))
	fmt.Println("")
}

func (c *Config) GetFlagSet() *flag.FlagSet {
	flag.CommandLine.Usage = c.PrintUsage
	return flag.CommandLine
}

func (c *Config) ParseFlag() {
	c.GetFlagSet().Parse(os.Args[1:])
}

func (c *Config) LoadConf() *Config {
	c.ParseFlag()
	config := instance

	m := &multiconfig.DefaultLoader{}
	m.Loader = multiconfig.MultiLoader(newLoader(constants.ServiceName))
	m.Validator = multiconfig.MultiValidator(
		&multiconfig.RequiredValidator{},
	)
	err := m.Load(config)
	if err != nil {
		panic(err)
	}

	loglevel := config.Log.Level
	logger.SetLevelByString(loglevel)
	logger.Debugf(nil, "load conf,config=[%+v]", config)
	return config
}
