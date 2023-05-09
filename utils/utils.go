package utils

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/JIeeiroSst/upload/dto"
)

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
			MinTemp:       record[0],
			MaxTemp:       record[1],
			Rainfall:      record[2],
			Evaporation:   record[3],
			Sunshine:      record[4],
			WindGustDir:   record[5],
			WindGustSpeed: record[6],
			WindDir9am:    record[7],
			WindDir3pm:    record[8],
			WindSpeed9am:  record[9],
			WindSpeed3pm:  record[10],
			Humidity9am:   record[11],
			Humidity3pm:   record[12],
			Pressure9am:   record[13],
			Pressure3pm:   record[14],
			Cloud9am:      record[15],
			Cloud3pm:      record[16],
			Temp9am:       record[17],
			Temp3pm:       record[18],
			RainToday:     record[19],
			RISK_MM:       record[20],
			RainTomorrow:  record[21],
		}
		weatherList = append(weatherList, weather)

	}
	return
}
