package dto

import "time"

type WeatherRequestDTO struct {
	MinTemp       string
	MaxTemp       string
	Rainfall      string
	Evaporation   string
	Sunshine      string
	WindGustDir   string
	WindGustSpeed string
	WindDir9am    string
	WindDir3pm    string
	WindSpeed9am  string
	WindSpeed3pm  string
	Humidity9am   string
	Humidity3pm   string
	Pressure9am   string
	Pressure3pm   string
	Cloud9am      string
	Cloud3pm      string
	Temp9am       string
	Temp3pm       string
	RainToday     string
	RISK_MM       string
	RainTomorrow  string
}

type ResponseBaseDTO struct {
	StatusCode    string
	ReasonCode    string
	ReasonMessage string
	TransactionID string
}

type FilterWeather struct {
	UserID        string
	CreatedAfter  time.Time
	CreatedBefore time.Time
}

