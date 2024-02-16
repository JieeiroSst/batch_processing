package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/JIeeiroSst/upload/model"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	CreateOrderInstantCardTx(tx *sql.Tx, weathers []model.Weather) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (repo *repository) CreateOrderInstantCardTx(tx *sql.Tx, weathers []model.Weather) error {
	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, w := range weathers {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

		valueArgs = append(valueArgs, w.MinTemp)
		valueArgs = append(valueArgs, w.MaxTemp)
		valueArgs = append(valueArgs, w.Rainfall)
		valueArgs = append(valueArgs, w.Evaporation)
		valueArgs = append(valueArgs, w.Sunshine)
		valueArgs = append(valueArgs, w.WindGustDir)
		valueArgs = append(valueArgs, w.WindGustSpeed)
		valueArgs = append(valueArgs, w.WindDir9am)
		valueArgs = append(valueArgs, w.WindDir3pm)
		valueArgs = append(valueArgs, w.WindSpeed9am)
		valueArgs = append(valueArgs, w.WindSpeed3pm)
		valueArgs = append(valueArgs, w.Humidity9am)
		valueArgs = append(valueArgs, w.Humidity3pm)
		valueArgs = append(valueArgs, w.Pressure9am)
		valueArgs = append(valueArgs, w.Pressure3pm)
		valueArgs = append(valueArgs, w.Cloud9am)
		valueArgs = append(valueArgs, w.Cloud3pm)
		valueArgs = append(valueArgs, w.Temp9am)
		valueArgs = append(valueArgs, w.Temp3pm)
		valueArgs = append(valueArgs, w.RainToday)
		valueArgs = append(valueArgs, w.RISK_MM)
		valueArgs = append(valueArgs, w.RainTomorrow)
	}
	smt := `INSERT INTO weathers(min_temp,max_temp,rainfall,evaporation,sunshine,wind_gust_dir
		,wind_gust_speed,wind_dir_9_am,wind_dir_3_pm,wind_speed_9_am,wind_speed_3_pm
		,humidity_9_am,humidity_3_pm,pressure_9_am,pressure_3_pm,cloud_9_am
		,cloud_3_pm,temp_9_am,temp_3_pm,rain_today,risk_mm,rain_tomorrow) VALUES %s`

	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	_, err := tx.Exec(smt, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}
