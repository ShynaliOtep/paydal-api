package models

import (
	"errors"
	"fmt"
	"time"
)

type Violation struct {
	ID          uint
	VolunteerId uint
	Volunteer   Volunteer
	LocationId  uint
	Location    Location
	CategoryId  uint
	Category    ViolationCategory
	Phone       string
	Address     string
	Comment     string
	Date        time.Time
	CityId      uint
	City        City
	Status      string
	Reason      string
	Medias      []ViolationMedia
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ViolationMedia struct {
	ID          uint
	ViolationId uint
	FileId      uint
	Violation   Violation
	File        File
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (v *Violation) Create() (*Violation, error) {

	var err error

	fmt.Println(v)
	err = DB.Create(&v).Error
	if err != nil {
		return &Violation{}, err
	}
	return v, nil
}

func (v *Violation) Save() (*Violation, error) {

	var err error
	fmt.Println(v)
	err = DB.Save(&v).Error
	if err != nil {
		return &Violation{}, err
	}
	return v, nil
}

func (m *ViolationMedia) Create() (*ViolationMedia, error) {

	var err error

	fmt.Println(m)
	err = DB.Create(&m).Error
	if err != nil {
		return &ViolationMedia{}, err
	}
	return m, nil
}

func GetViolationById(id uint) (Violation, error) {

	var violation Violation

	if err := DB.
		Preload("Medias").
		Preload("Category").
		Preload("Location").
		Preload("City").
		Preload("Medias.File").
		First(&violation, id).Error; err != nil {
		return violation, errors.New("VНарушение не найден!!")
	}

	return violation, nil

}

func ViolationSearchInMap(lat float64, lon float64, dist float64) ([]Violation, error) {

	var violations []Violation
	// надо добавить status="confirmed"
	if err := DB.
		//Preload("Medias").
		Joins("JOIN violation_categories ON violations.category_id = violation_categories.id").
		Joins("JOIN cities ON violations.city_id = cities.id").
		//Preload("Medias.File").
		Joins("JOIN locations ON location_id = locations.id").
		Select("violations.*, violation_categories.*, cities.*, locations.id as l, locations.lon, locations.lat,  (6371 * ACOS(COS(RADIANS(?)) * COS(RADIANS(locations.lat)) * COS(RADIANS(locations.lon) - RADIANS(?)) + SIN(RADIANS(?)) * SIN(RADIANS(locations.lat)))) AS distance",
			lat, lon, lat).
		Table("violations").
		Having("distance < ?", dist).
		//Where("ST_Distance_Sphere(locations.point, POINT(?,?)) < ?", lat, lon, dist).
		Scan(&violations).Error; err != nil {
		return violations, errors.New("Нарушение не найден!")
	}

	return violations, nil

}

func ViolationsForPolice() ([]Violation, error) {

	var violations []Violation
	if err := DB.
		Preload("Medias").
		Preload("Category").
		Preload("Location").
		Preload("City").
		Preload("Medias.File").
		Where("status = ?", "created").
		Find(&violations).Error; err != nil {
		return violations, errors.New("Нарушение не найден!!")
	}

	return violations, nil

}

func VolunteerViolations(volunteerId uint) ([]Violation, error) {

	var violations []Violation
	if err := DB.
		Preload("Medias").
		Preload("Category").
		Preload("Location").
		Preload("City").
		Preload("Medias.File").
		Where("volunteer_id = ?", volunteerId).
		Find(&violations).Error; err != nil {
		return violations, errors.New("Нарушение не найден!!")
	}

	return violations, nil

}
