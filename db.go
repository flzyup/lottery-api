/**
 * Copyright © 2017-2018 Yonnie @ i4o.xyz . All rights reserved.
 *
 * FileName: main/db.go
 *
 * Author: FLZYUP Lu
 * Email: yonnie.lu.inc@gmail.com
 * Date: 2018-12-26 17:52
 * Description:
 * History:
 *   <Author>      <Time>    <version>    <desc>
 *   YonnieLu      2018-12-26 17:52    1.0          Create
 */
package main

import (
	"database/sql"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	CarModel string `json:"car_model"'`
}

type Award struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Count        int    `json:"count"`
	AwardedCount int    `json:"awarded_count"`
	CarModel     string `json:"car_model"'`
	BatchId      int    `json:"batch_id"`
}

type UserAward struct {
	Id           int    `json:"id"`
	UserId       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	AwardId      int    `json:"award_id"`
	AwardBatchId int    `json:"award_batch_id"`
	AwardName    string `json:"award_name"`
}

func queryUser() ([]*User, error) {
	stmtQueryUser, err := db.Prepare("SELECT id, name, car_model FROM lt_user")
	defer stmtQueryUser.Close()

	rows, err := stmtQueryUser.Query()
	defer rows.Close()

	result := make([]*User, 0)

	if err != nil {
		log().Errorf("error: %v", err)
		return nil, err
	} else {
		for rows.Next() {
			user := User{}
			err = rows.Scan(&user.Id, &user.Name, &user.CarModel)
			result = append(result, &user)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func lotteryUser(limit int) ([]*User, error) {
	stmt, err := db.Prepare("SELECT id, name, car_model FROM lt_user" +
		" ORDER BY RAND() LIMIT ?")

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows *sql.Rows

	rows, err = stmt.Query(limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		rows.Scan(&user.Id, &user.Name, &user.CarModel)
		result = append(result, user)
	}
	return result, nil
}

func lotteryUserEnd(carModel string, awardBatchId int, limit int) ([]*User, error) {
	queryString := ""
	if carModel != "" {
		queryString = "SELECT id, name, car_model FROM lt_user" +
			" WHERE	id NOT IN ( SELECT user_id FROM lt_user_award WHERE award_batch_id = ? )" +
			" and car_model = ?"
		queryString += " ORDER BY RAND() LIMIT ?"

	} else {
		queryString = "SELECT id, name, car_model FROM lt_user" +
			" WHERE	id NOT IN ( SELECT user_id FROM lt_user_award WHERE award_batch_id = ? )" +
			" ORDER BY RAND() LIMIT ?"
	}
	stmt, err := db.Prepare(queryString)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows *sql.Rows

	if carModel != "" {
		rows, err = stmt.Query(awardBatchId, carModel, limit)
	} else {
		rows, err = stmt.Query(awardBatchId, limit)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		rows.Scan(&user.Id, &user.Name, &user.CarModel)
		result = append(result, user)
	}
	return result, nil
}

func queryAward() ([]*Award, error) {
	stmtQueryUser, err := db.Prepare("SELECT id, name, count, car_model, batch_id, (select count(*) from lt_user_award where award_id = a.id) as awarded_count FROM lt_award as a")
	defer stmtQueryUser.Close()

	rows, err := stmtQueryUser.Query()
	defer rows.Close()

	result := make([]*Award, 0)

	if err != nil {
		log().Errorf("error: %v", err)
		return nil, err
	} else {
		for rows.Next() {
			award := Award{}
			err = rows.Scan(&award.Id, &award.Name, &award.Count, &award.CarModel, &award.BatchId, &award.AwardedCount)
			result = append(result, &award)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func queryAwardById(awardId int) (*Award, error) {
	stmtQueryUser, err := db.Prepare("SELECT id, name, count, car_model, batch_id FROM lt_award where id = ?")
	defer stmtQueryUser.Close()

	rows, err := stmtQueryUser.Query(awardId)
	defer rows.Close()

	if err != nil {
		log().Errorf("error: %v", err)
		return nil, err
	} else {
		if rows.Next() {
			award := Award{}
			err = rows.Scan(&award.Id, &award.Name, &award.Count, &award.CarModel, &award.BatchId)
			return &award, nil
		} else {
			return nil, nil
		}

	}
}

func queryAwardUsers() ([]*UserAward, error) {
	stmt, err := db.Prepare("select ua.user_id, ua.award_id, u.`name`, a.`name` from lt_user_award" +
		" as ua left JOIN lt_user as u on ua.user_id = u.id" +
		" left JOIN lt_award as a on ua.award_id = a.id" +
		" ORDER BY ua.create_time desc, a.id ASC")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	result := make([]*UserAward, 0)
	for rows.Next() {
		ret := UserAward{}
		rows.Scan(&ret.UserId, &ret.AwardId, &ret.UserName, &ret.AwardName)
		result = append(result, &ret)
	}
	return result, nil
}

func queryAwardUserCount(awardId int) (int, error) {
	stmt, err := db.Prepare("select count(*) from lt_user_award where award_id = ?")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(awardId)

	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if rows.Next() {
		count := 0
		rows.Scan(&count)
		return count, nil
	} else {
		return 0, nil
	}

	return 1, nil
}

func insertUserAward(userId int, awardId int, awardBatchId int) (int, error) {
	stmt, err := db.Prepare("INSERT INTO lt_user_award (user_id, award_id, award_batch_id) values (?, ?, ?)")
	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(userId, awardId, awardBatchId)

	if err != nil {
		return -1, err
	}

	return 1, nil
}
