-- +goose Up
CREATE TABLE weathers(
    id int primary key auto_increment,
    min_temp DOUBLE,
    max_temp DOUBLE,
    rainfall DOUBLE,
    evaporation DOUBLE,
    sunshine DOUBLE,
    wind_gust_dir TEXT,
    wind_gust_speed DOUBLE,
    wind_dir_9_am TEXT,
    wind_dir_3_pm TEXT,
    wind_speed_9_am DOUBLE,
    wind_speed_3_pm DOUBLE,
    humidity_9_am DOUBLE,
    humidity_3_pm DOUBLE,
    pressure_9_am DOUBLE,
    pressure_3_pm DOUBLE,
    cloud_9_am DOUBLE,
    cloud_3_pm DOUBLE,
    temp_9_am DOUBLE,
    temp_3_pm DOUBLE,
    rain_today BOOLEAN,
    risk_mm DOUBLE,
    rain_tomorrow BOOLEAN
); 

-- +goose Down
DROP TABLE weathers;
