-- +goose Up
CREATE TABLE weathers(
    id int primary key auto_increment,
    min_temp TEXT,
    max_temp TEXT,
    rainfall TEXT,
    evaporation TEXT,
    sunshine TEXT,
    wind_gust_dir TEXT,
    wind_gust_speed TEXT,
    wind_dir_9_am TEXT,
    wind_dir_3_pm TEXT,
    wind_speed_9_am TEXT,
    wind_speed_3_pm TEXT,
    humidity_9_am TEXT,
    humidity_3_pm TEXT,
    pressure_9_am TEXT,
    pressure_3_pm TEXT,
    cloud_9_am TEXT,
    cloud_3_pm TEXT,
    temp_9_am TEXT,
    temp_3_pm TEXT,
    rain_today TEXT,
    risk_mm TEXT,
    rain_tomorrow TEXT
); 

-- +goose Down
DROP TABLE weathers;
