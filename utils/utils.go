package utils

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/JIeeiroSst/upload/dto"
)

func stringTofloat64(str string) float64 {
	feetFloat, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		log.Println(err)
	}

	return feetFloat
}

func stringToBool(str string) bool {
	boolValue, err := strconv.ParseBool(str)
	if err != nil {
		log.Println(err)
	}
	return boolValue
}

func readData(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func CreateWeatherList(file string) (weatherList []dto.WeatherRequestDTO) {
	records, err := readData(file)

	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {

		weather := dto.WeatherRequestDTO{
			MinTemp:       stringTofloat64(record[0]),
			MaxTemp:       stringTofloat64(record[1]),
			Rainfall:      stringTofloat64(record[2]),
			Evaporation:   stringTofloat64(record[3]),
			Sunshine:      stringTofloat64(record[4]),
			WindGustDir:   record[5],
			WindGustSpeed: stringTofloat64(record[6]),
			WindDir9am:    record[7],
			WindDir3pm:    record[8],
			WindSpeed9am:  stringTofloat64(record[9]),
			WindSpeed3pm:  stringTofloat64(record[10]),
			Humidity9am:   stringTofloat64(record[11]),
			Humidity3pm:   stringTofloat64(record[12]),
			Pressure9am:   stringTofloat64(record[13]),
			Pressure3pm:   stringTofloat64(record[14]),
			Cloud9am:      stringTofloat64(record[15]),
			Cloud3pm:      stringTofloat64(record[16]),
			Temp9am:       stringTofloat64(record[17]),
			Temp3pm:       stringTofloat64(record[18]),
			RainToday:     stringToBool(record[19]),
			RISK_MM:       stringTofloat64(record[20]),
			RainTomorrow:  stringToBool(record[21]),
		}
		weatherList = append(weatherList, weather)

	} 
	return 
}
