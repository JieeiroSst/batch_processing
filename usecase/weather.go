package usecase

import (
	"context"
	"errors"
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
	upsertBigQueryWeather(weathers []dto.WeatherRequestDTO) error
	processWeather(weathers <-chan dto.WeatherRequestDTO, batchSize int)
	produceWeather(weathers []dto.WeatherRequestDTO, to chan dto.WeatherRequestDTO)
	InsertWeather(weathers []dto.WeatherRequestDTO)

	ValidateWeatherData(weathers []dto.WeatherRequestDTO) (*dto.ValidationResult, error)
	validateWeatherBatch(weathers []dto.WeatherRequestDTO, startIndex int) ([]dto.WeatherRequestDTO, []dto.WeatherRequestDTO, []dto.ValidationError)
	validateSingleWeather(weather dto.WeatherRequestDTO, index int) []dto.ValidationError
}

func NewUsecase(repository repository.Repository) Usecase {
	return &usecase{
		repository: repository,
	}
}

func (u *usecase) upsertBigQueryWeather(weathers []dto.WeatherRequestDTO) error {
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

	ctx := context.Background()
	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v\n", err)
			tx.Rollback()
		}
	}()

	if err = u.repository.CreateOrderInstantCardTx(tx, weatherModels); err != nil {
		log.Printf("Failed to insert weather data: %v\n", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return err
	}

	fmt.Printf("Successfully committed batch of %d records\n", len(weathers))
	return nil
}

func (u *usecase) processWeather(weathers <-chan dto.WeatherRequestDTO, batchSize int) {
	var batch []dto.WeatherRequestDTO
	for weather := range weathers {
		batch = append(batch, weather)
		if len(batch) == batchSize {
			if err := u.upsertBigQueryWeather(batch); err != nil {
				log.Printf("Error processing batch: %v\n", err)
			}
			batch = []dto.WeatherRequestDTO{}
		}
	}
	if len(batch) > 0 {
		if err := u.upsertBigQueryWeather(batch); err != nil {
			log.Printf("Error processing final batch: %v\n", err)
		}
	}
}

func (u *usecase) produceWeather(weathers []dto.WeatherRequestDTO, to chan dto.WeatherRequestDTO) {
	defer close(to)
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

var (
	validationBatchSize = 1000
)

func (u *usecase) ValidateWeatherData(weathers []dto.WeatherRequestDTO) (*dto.ValidationResult, error) {
	if len(weathers) == 0 {
		return nil, errors.New("empty weather data")
	}

	fmt.Printf("Starting validation for %d records\n", len(weathers))

	result := &dto.ValidationResult{
		ValidData:   make([]dto.WeatherRequestDTO, 0),
		InvalidData: make([]dto.WeatherRequestDTO, 0),
		Errors:      make([]dto.ValidationError, 0),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	totalBatches := (len(weathers) + validationBatchSize - 1) / validationBatchSize

	type batchResult struct {
		valid   []dto.WeatherRequestDTO
		invalid []dto.WeatherRequestDTO
		errors  []dto.ValidationError
	}
	resultChan := make(chan batchResult, totalBatches)

	for i := 0; i < len(weathers); i += validationBatchSize {
		end := i + validationBatchSize
		if end > len(weathers) {
			end = len(weathers)
		}

		batch := weathers[i:end]
		startIndex := i

		wg.Add(1)
		go func(batch []dto.WeatherRequestDTO, startIdx int) {
			defer wg.Done()

			valid, invalid, errs := u.validateWeatherBatch(batch, startIdx)

			resultChan <- batchResult{
				valid:   valid,
				invalid: invalid,
				errors:  errs,
			}
		}(batch, startIndex)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for batchRes := range resultChan {
		mu.Lock()
		result.ValidData = append(result.ValidData, batchRes.valid...)
		result.InvalidData = append(result.InvalidData, batchRes.invalid...)
		result.Errors = append(result.Errors, batchRes.errors...)
		mu.Unlock()
	}

	fmt.Printf("Validation completed: %d valid, %d invalid, %d errors\n",
		len(result.ValidData), len(result.InvalidData), len(result.Errors))

	return result, nil
}

func (u *usecase) validateWeatherBatch(weathers []dto.WeatherRequestDTO, startIndex int) ([]dto.WeatherRequestDTO, []dto.WeatherRequestDTO, []dto.ValidationError) {
	fmt.Printf("Validating batch starting at index %d with %d records\n", startIndex, len(weathers))

	validData := make([]dto.WeatherRequestDTO, 0)
	invalidData := make([]dto.WeatherRequestDTO, 0)
	errors := make([]dto.ValidationError, 0)

	for i, weather := range weathers {
		actualIndex := startIndex + i
		validationErrors := u.validateSingleWeather(weather, actualIndex)

		if len(validationErrors) > 0 {
			invalidData = append(invalidData, weather)
			errors = append(errors, validationErrors...)
		} else {
			validData = append(validData, weather)
		}
	}

	return validData, invalidData, errors
}

func (u *usecase) validateSingleWeather(weather dto.WeatherRequestDTO, index int) []dto.ValidationError {
	var errors []dto.ValidationError

	if weather.MinTemp < "-50" || weather.MinTemp > "60" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "MinTemp",
			Message: fmt.Sprintf("MinTemp must be between -50 and 60, got %.2f", weather.MinTemp),
		})
	}

	if weather.MaxTemp < "-50" || weather.MaxTemp > "60" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "MaxTemp",
			Message: fmt.Sprintf("MaxTemp must be between -50 and 60, got %.2f", weather.MaxTemp),
		})
	}

	if weather.MinTemp > weather.MaxTemp {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "MinTemp/MaxTemp",
			Message: fmt.Sprintf("MinTemp (%.2f) cannot be greater than MaxTemp (%.2f)", weather.MinTemp, weather.MaxTemp),
		})
	}

	if weather.Rainfall < "0" || weather.Rainfall > "1000" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "Rainfall",
			Message: fmt.Sprintf("Rainfall must be between 0 and 1000, got %.2f", weather.Rainfall),
		})
	}

	if weather.Evaporation < "0" || weather.Evaporation > "100" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "Evaporation",
			Message: fmt.Sprintf("Evaporation must be between 0 and 100, got %.2f", weather.Evaporation),
		})
	}

	if weather.Sunshine < "0" || weather.Sunshine > "24" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "Sunshine",
			Message: fmt.Sprintf("Sunshine must be between 0 and 24 hours, got %.2f", weather.Sunshine),
		})
	}

	if weather.WindGustSpeed < "0" || weather.WindGustSpeed > "200" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "WindGustSpeed",
			Message: fmt.Sprintf("WindGustSpeed must be between 0 and 200, got %.2f", weather.WindGustSpeed),
		})
	}

	if weather.WindSpeed9am < "0" || weather.WindSpeed9am > "200" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "WindSpeed9am",
			Message: fmt.Sprintf("WindSpeed9am must be between 0 and 200, got %.2f", weather.WindSpeed9am),
		})
	}

	if weather.WindSpeed3pm < "0" || weather.WindSpeed3pm > "200" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "WindSpeed3pm",
			Message: fmt.Sprintf("WindSpeed3pm must be between 0 and 200, got %.2f", weather.WindSpeed3pm),
		})
	}

	if weather.Humidity9am < "0" || weather.Humidity9am > "100" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "Humidity9am",
			Message: fmt.Sprintf("Humidity9am must be between 0 and 100, got %.2f", weather.Humidity9am),
		})
	}

	if weather.Humidity3pm < "0" || weather.Humidity3pm > "100" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "Humidity3pm",
			Message: fmt.Sprintf("Humidity3pm must be between 0 and 100, got %.2f", weather.Humidity3pm),
		})
	}

	if weather.Pressure9am < "900" || weather.Pressure9am > "1100" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "Pressure9am",
			Message: fmt.Sprintf("Pressure9am must be between 900 and 1100, got %.2f", weather.Pressure9am),
		})
	}

	if weather.Pressure3pm < "900" || weather.Pressure3pm > "1100" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "Pressure3pm",
			Message: fmt.Sprintf("Pressure3pm must be between 900 and 1100, got %.2f", weather.Pressure3pm),
		})
	}

	if weather.Cloud9am < "0" || weather.Cloud9am > "9" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "Cloud9am",
			Message: fmt.Sprintf("Cloud9am must be between 0 and 9, got %.2f", weather.Cloud9am),
		})
	}

	if weather.Cloud3pm < "0" || weather.Cloud3pm > "9" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "Cloud3pm",
			Message: fmt.Sprintf("Cloud3pm must be between 0 and 9, got %.2f", weather.Cloud3pm),
		})
	}

	validDirections := map[string]bool{
		"N": true, "NNE": true, "NE": true, "ENE": true,
		"E": true, "ESE": true, "SE": true, "SSE": true,
		"S": true, "SSW": true, "SW": true, "WSW": true,
		"W": true, "WNW": true, "NW": true, "NNW": true,
	}

	if weather.WindGustDir != "" && !validDirections[weather.WindGustDir] {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "WindGustDir",
			Message: fmt.Sprintf("Invalid wind direction: %s", weather.WindGustDir),
		})
	}

	if weather.WindDir9am != "" && !validDirections[weather.WindDir9am] {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "WindDir9am",
			Message: fmt.Sprintf("Invalid wind direction: %s", weather.WindDir9am),
		})
	}

	if weather.WindDir3pm != "" && !validDirections[weather.WindDir3pm] {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "WindDir3pm",
			Message: fmt.Sprintf("Invalid wind direction: %s", weather.WindDir3pm),
		})
	}

	if weather.RISK_MM < "0" {
		errors = append(errors, dto.ValidationError{
			Index:   index,
			Field:   "RISK_MM",
			Message: fmt.Sprintf("RISK_MM cannot be negative, got %.2f", weather.RISK_MM),
		})
	}

	return errors
}
