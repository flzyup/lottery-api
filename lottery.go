/**
 * Copyright © 2017-2018 Yonnie @ i4o.xyz . All rights reserved.
 *
 * FileName: rx8_lottery_api/lottery.go
 *
 * Author: FLZYUP Lu
 * Email: yonnie.lu.inc@gmail.com
 * Date: 2018-12-26 17:35
 * Description:
 * History:
 *   <Author>      <Time>    <version>    <desc>
 *   YonnieLu      2018-12-26 17:35    1.0          Create
 */
package main

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
	"i4o.xyz/rx8lottery/logger"
	"io/ioutil"
	"strconv"
	"time"
)

var (
	db *sql.DB
)

func main() {
	conf := &Config{}

	configFile, err := ioutil.ReadFile("config.yml")

	if err != nil {
		log().Error("read config file error.", err)
	}

	if err = yaml.Unmarshal(configFile, conf); err != nil {
		log().Error("read yaml failed", err)
	} else {

	}

	ginHttpServer := gin.Default()
	// Test cors and allow all request from any origin
	ginHttpServer.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routeApi(ginHttpServer)
	routeHtml(ginHttpServer)

	db, err = sql.Open("mysql", conf.Mysql.Dsn)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	ginHttpServer.Run(conf.Http.Ip + ":" + strconv.Itoa(conf.Http.Port))
}

func log() *logging.Logger {
	return logger.GetLogger("lottery")
}
