package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/JIeeiroSst/upload/dto"
	"github.com/JIeeiroSst/upload/model"
	"github.com/JIeeiroSst/upload/repository"
)

type usecase struct {
	repository repository.Repository
}

type Usecase interface {
	upsertBigQueryWeather(weathers []dto.WeatherRequestDTO)
	processWeather(weathers <-chan dto.WeatherRequestDTO, batchSize int)
	produceWeather(weathers []dto.WeatherRequestDTO, to chan dto.WeatherRequestDTO)
	InsertWeather(weathers []dto.WeatherRequestDTO)
}

func NewUsecase(repository repository.Repository) Usecase {
	return &usecase{
		repository: repository,
	}
}

func (u *usecase) upsertBigQueryWeather(weathers []dto.WeatherRequestDTO) {
	fmt.Printf("Processing batch of %d\n", len(weathers))
	var weatherModels []model.Weather
	for _, weather := range weathers {
		weatherModel := model.Weather{
			MinTemp:       weather.MinTemp,
			MaxTemp:       weather.MaxTemp,
			Rainfall:      weather.Rainfall,
			Evaporation:   weather.Evaporation,
			Sunshine:      weather.Sunshine,
			WindGustDir:   weather.WindGustDir,
			WindGustSpeed: weather.WindGustSpeed,
			WindDir9am:    weather.WindDir9am,
			WindDir3pm:    weather.WindDir3pm,
			WindSpeed9am:  weather.WindSpeed9am,
			WindSpeed3pm:  weather.WindSpeed3pm,
			Humidity9am:   weather.Humidity9am,
			Humidity3pm:   weather.Humidity3pm,
			Pressure9am:   weather.Pressure9am,
			Pressure3pm:   weather.Pressure3pm,
			Cloud9am:      weather.Cloud9am,
			Cloud3pm:      weather.Cloud3pm,
			Temp9am:       weather.Temp9am,
			Temp3pm:       weather.Temp3pm,
			RainToday:     weather.RainToday,
			RISK_MM:       weather.RISK_MM,
			RainTomorrow:  weather.RainTomorrow,
		}
		weatherModels = append(weatherModels, weatherModel)
	}
	if err := u.repository.CreateOrderInstantCard(context.Background(), weatherModels); err != nil {
		log.Println(err)
	}
	fmt.Println()
}

func (u *usecase) processWeather(weathers <-chan dto.WeatherRequestDTO, batchSize int) {
	var batch []dto.WeatherRequestDTO
	for weather := range weathers {
		batch = append(batch, weather)
		if len(batch) == batchSize {
			u.upsertBigQueryWeather(batch)
			batch = []dto.WeatherRequestDTO{}
		}
	}
	if len(batch) > 0 {
		u.upsertBigQueryWeather(batch)
	}
}

func (u *usecase) produceWeather(weathers []dto.WeatherRequestDTO, to chan dto.WeatherRequestDTO) {
	for _, weather := range weathers {
		to <- dto.WeatherRequestDTO{
			MinTemp:       weather.MinTemp,
			MaxTemp:       weather.MaxTemp,
			Rainfall:      weather.Rainfall,
			Evaporation:   weather.Evaporation,
			Sunshine:      weather.Sunshine,
			WindGustDir:   weather.WindGustDir,
			WindGustSpeed: weather.WindGustSpeed,
			WindDir9am:    weather.WindDir9am,
			WindDir3pm:    weather.WindDir3pm,
			WindSpeed9am:  weather.WindSpeed9am,
			WindSpeed3pm:  weather.WindSpeed3pm,
			Humidity9am:   weather.Humidity9am,
			Humidity3pm:   weather.Humidity3pm,
			Pressure9am:   weather.Pressure9am,
			Pressure3pm:   weather.Pressure3pm,
			Cloud9am:      weather.Cloud9am,
			Cloud3pm:      weather.Cloud3pm,
			Temp9am:       weather.Temp9am,
			Temp3pm:       weather.Temp3pm,
			RainToday:     weather.RainToday,
			RISK_MM:       weather.RISK_MM,
			RainTomorrow:  weather.RainTomorrow,
		}
	}
}

const batchSize = 1000

func (u *usecase) InsertWeather(weathers []dto.WeatherRequestDTO) {
	var wg sync.WaitGroup
	audits := make(chan dto.WeatherRequestDTO)
	wg.Add(1)
	go func() {
		defer wg.Done()
		u.processWeather(audits, batchSize)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		u.produceWeather(weathers, audits)
		close(audits)
	}()
	wg.Wait()
	fmt.Println("Complete")
}
