/**
 * Copyright © 2017-2018 Shanghai Lebai Robotic Co., Ltd. All rights reserved.
 *
 * FileName: main/http.go
 *
 * Author: Yonnie Lu
 * Email: zhangyong.lu@lebai.ltd
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
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func routeApi(root *gin.Engine) *gin.RouterGroup {
	api := root.Group("/api")

	getUser(api)
	getAward(api)
	lottery(api)
	getAwardUsers(api)

	return api
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
			jsonData(ctx, result)
		}
	})
}

const (
	footPlanAwardId = 1
	yuanxingUserId  = 1
	renCiUserId     = 9
)

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
					switch awardId {
					case footPlanAwardId:
						result = append(result, renCiUserId)
						_, err = insertUserAward(renCiUserId, footPlanAwardId, award.BatchId)

						if err != nil {
							jsonError(ctx, err)
							return
						}
					default:
						users, err = lotteryUserEnd(award.CarModel, award.BatchId, limit)
						if err != nil {
							jsonError(ctx, err)
							return
						}
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
