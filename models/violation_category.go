package models

import "time"

type ViolationCategory struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetViolationCategories() ([]ViolationCategory, error) {
	var categories []ViolationCategory
	err := DB.Find(&categories).Error
	if err != nil {
		return categories, err
	}
	return categories, nil
}
