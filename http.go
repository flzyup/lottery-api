/**
 * Copyright © 2017-2018 Yonnie @ i4o.xyz . All rights reserved.
 *
 * FileName: main/http.go
 *
 * Author: FLZYUP Lu
 * Email: yonnie.lu.inc@gmail.com
 * Date: 2018-12-26 18:13
 * Description:
 * History:
 *   <Author>      <Time>    <version>    <desc>
 *   YonnieLu      2018-12-26 18:13    1.0          Create
 */
package main

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type AwardedList struct {
	Award
	Users []*User `json:"users"`
}

func routeApi(root *gin.Engine) *gin.RouterGroup {
	api := root.Group("/api")

	getUser(api)
	getAward(api)
	lottery(api)
	getAwardUsers(api)

	return api
}

func routeHtml(root *gin.Engine) {
	root.Use(static.Serve("/", static.LocalFile("./dist", true)))
	//root.StaticFS("/dashboard", gin.Dir("templates", true))
	//root.StaticFS("/assets", gin.Dir("templates/assets", true))
}

// get, post, put, delete

func getUser(api *gin.RouterGroup) {
	api.GET("/user", func(ctx *gin.Context) {
		if result, err := queryUser(); err != nil {
			jsonError(ctx, err)
		} else {
			jsonData(ctx, result)
		}
	})
}

func getAward(api *gin.RouterGroup) {
	api.GET("/award", func(ctx *gin.Context) {
		if result, err := queryAward(); err != nil {
			jsonError(ctx, err)
		} else {
			jsonData(ctx, result)
		}
	})
}

func getAwardUsers(api *gin.RouterGroup) {
	api.GET("/award/users", func(ctx *gin.Context) {
		if result, err := queryAwardUsers(); err != nil {
			jsonError(ctx, err)
		} else {
			ret := make([]*AwardedList, 0)

			for i := 0; i < len(result); i++ {
				au := result[i]

				var al *AwardedList
				for j := 0; j < len(ret); j++ {
					if ret[j].Id == au.AwardId {
						al = ret[j]
						break
					}
				}

				if al == nil {
					al = &AwardedList{}
					al.Id = au.AwardId
					al.Name = au.AwardName

					al.Users = make([]*User, 0)
					ret = append(ret, al)
				}

				al.Users = append(al.Users, &User{Id: au.UserId, Name: au.UserName})
			}

			jsonData(ctx, ret)
		}
	})
}

// Lottery
func lottery(api *gin.RouterGroup) {
	api.GET("/lottery", func(ctx *gin.Context) {
		awardId, _ := strconv.Atoi(ctx.Query("awardId"))
		stop, _ := strconv.Atoi(ctx.Query("stop"))
		count, _ := strconv.Atoi(ctx.Query("count"))

		if awardId > 0 {
			result := make([]int, 0)

			award, err := queryAwardById(awardId)

			if err != nil {
				jsonError(ctx, errors.New("奖品信息不存在！"))
				return
			}

			awardedCount, err := queryAwardUserCount(awardId)

			if err != nil {
				jsonError(ctx, err)
				return
			}

			limit := award.Count - awardedCount

			limit = int(math.Min(float64(count), float64(limit)))

			if limit == 0 {
				jsonError(ctx, errors.New(fmt.Sprintf("%d个%s已全部抽完！", award.Count, award.Name)))
				return
			}

			var users []*User
			switch stop {
			case 1:
				{ // 如果是停止抽奖
					users, err = lotteryUserEnd(award.CarModel, award.BatchId, limit)
					if err != nil {
						jsonError(ctx, err)
						return
					}
				}
			default:
				{ // 如果是随机展示
					users, err = lotteryUser(limit)
					if err != nil {
						jsonError(ctx, err)
						return
					}
				}
			}

			for _, v := range users {
				if stop == 1 && award.BatchId == 2 && v.Id == yuanxingUserId {
					// 老大不要奖品，哈哈哈
					continue
				}
				result = append(result, v.Id)
				if stop == 1 {
					log().Infof("lottery user awarded: %v", v)
					_, err = insertUserAward(v.Id, awardId, award.BatchId)
					if err != nil {
						jsonError(ctx, err)
						return
					}
				} else {
					log().Infof("lottery user temp: %v", v)
				}
			}

			jsonData(ctx, result)
		} else {
			jsonError(ctx, errors.New("illegal request"))
		}
	})
}

func jsonData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Msg:  "ok",
		Data: data,
	})
}

func jsonError(ctx *gin.Context, err error) {
	log().Error(err)
	ctx.JSON(http.StatusBadRequest, Response{
		Msg:  "fail",
		Data: err.Error(),
	})
}
