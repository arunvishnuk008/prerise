package persistence

import (
	"crimson-sunrise.site/pkg/common"
	"crimson-sunrise.site/pkg/model"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"
)

func GetAllHotTakes(db *sql.DB) ([]model.HotTake, error) {
	const query = "" +
		"SELECT ht.id, ht.title, ht.details, ht.created_at from hot_takes ht"
	resultSet, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var hotTakes []model.HotTake
	defer func(resultSet *sql.Rows) {
		err := resultSet.Close()
		if err != nil {
			log.Printf("error closing resultset %s", err.Error())
		}
	}(resultSet)
	for resultSet.Next() {
		var hotTake model.HotTake
		err = resultSet.Scan(&hotTake.ID, &hotTake.Title, &hotTake.Details, &hotTake.CreatedAtInternal)
		if err != nil {
			return nil, err
		}
		loc, _ := time.LoadLocation("Asia/Kolkata")
		createdAt, err := common.EpochSecondToTime(strconv.Itoa(hotTake.CreatedAtInternal))
		if err != nil {
			return nil, err
		}
		createdAt = createdAt.In(loc)
		hotTake.CreatedAt = createdAt.Format(time.RFC3339)
		hotTakes = append(hotTakes, hotTake)
	}
	return hotTakes, nil
}

func GetHotTakeByID(db *sql.DB, id int64) (*model.HotTake, error) {
	const query = "" +
		"SELECT ht.id, ht.title, ht.details, ht.created_at from hot_takes ht where ht.id=?"
	resultSet, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer func(resultSet *sql.Rows) {
		err := resultSet.Close()
		if err != nil {
			log.Printf("error closing resultset %s", err.Error())
		}
	}(resultSet)
	if resultSet.Next() {
		var hotTake model.HotTake
		err = resultSet.Scan(&hotTake.ID, &hotTake.Title, &hotTake.Details, &hotTake.CreatedAtInternal)
		if err != nil {
			return nil, err
		}
		loc, _ := time.LoadLocation("Asia/Kolkata")
		createdAt, err := common.EpochSecondToTime(strconv.Itoa(hotTake.CreatedAtInternal))
		if err != nil {
			return nil, err
		}
		createdAt = createdAt.In(loc)
		hotTake.CreatedAt = createdAt.Format(time.RFC3339)
		return &hotTake, nil
	}
	return nil, errors.New("hot take not found")
}

func AddHotTake(db *sql.DB, hotTake model.HotTake) (*model.HotTake, error) {
	const query =
		"INSERT INTO hot_takes(title, details, created_at) VALUES (?,?,?)"
	insert, err := db.Exec(query, hotTake.Title, hotTake.Details, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	lastInsertedId, err := insert.LastInsertId()
	if err != nil {
		return nil, err
	}
	return GetHotTakeByID(db, lastInsertedId)
}

func DeleteHotTake(db *sql.DB, id int64) (bool, error) {
	const query =
		"DELETE FROM hot_takes WHERE id=?"
	deleted, err := db.Exec(query, id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := deleted.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}