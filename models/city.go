package models

import "time"

type City struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetCities() ([]City, error) {
	var city []City
	err := DB.Find(&city).Error
	if err != nil {
		return city, err
	}
	return city, nil
}
