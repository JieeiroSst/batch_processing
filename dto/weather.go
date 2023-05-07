package dto

import "time"

type WeatherRequestDTO struct {
	MinTemp       float64
	MaxTemp       float64
	Rainfall      float64
	Evaporation   float64
	Sunshine      float64
	WindGustDir   string
	WindGustSpeed float64
	WindDir9am    string
	WindDir3pm    string
	WindSpeed9am  float64
	WindSpeed3pm  float64
	Humidity9am   float64
	Humidity3pm   float64
	Pressure9am   float64
	Pressure3pm   float64
	Cloud9am      float64
	Cloud3pm      float64
	Temp9am       float64
	Temp3pm       float64
	RainToday     bool
	RISK_MM       float64
	RainTomorrow  bool
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

