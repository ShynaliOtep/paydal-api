package models

import (
	"errors"
	"fmt"
	"github.com/ShynaliOtep/paydal-api/utils/token"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strconv"
	"strings"
	"time"
)

type User struct {
	gorm.Model
	Phone    string     `gorm:"size:255;not null;unique" json:"phone"`
	Password string     `gorm:"size:255;not null;" json:"password"`
	Fio      string     `gorm:"size:255;not null;" json:"fio"`
	Iin      string     `gorm:"size:12;null;" json:"iin"`
	IsPolice bool       `json:"is_police"`
	IsAdmin  bool       `json:"is_admin"`
	Birthday *time.Time `json:"birthday"`
	Avatar   *File      `json:"avatar"`
	AvatarId *uint
	Wallet   Wallet
}

func (u *User) SaveUser() (*User, error) {

	var err error
	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateUser() (*User, error) {

	var err error
	err = DB.Save(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Fio = html.EscapeString(strings.TrimSpace(u.Fio))

	return nil

}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(phone string, password string) (accessToken string, refreshToken string, err error) {

	u := User{}

	err = DB.Model(User{}).Where("phone = ?", phone).Take(&u).Error

	if err != nil {
		return "", "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", "", err
	}

	accessToken, refreshToken, err = token.GenerateTokenPair(u.ID)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := DB.
		Preload("Wallet").
		Preload("Avatar").
		First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func GetUserByUsername(username string) (User, error) {

	var u User

	if err := DB.Where("username=?", username).First(&u).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func GetUserByPhone(phone string) (User, error) {

	var u User

	if err := DB.Where("phone=?", phone).First(&u).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}

// ExtractBirthDateFromIIN извлекает дату рождения из ИИН
func ExtractBirthDateFromIIN(iin string) (time.Time, error) {
	if len(iin) != 12 {
		return time.Time{}, fmt.Errorf("invalid IIN length")
	}

	yearPrefix := "18"
	if iin[6] == '3' || iin[6] == '4' {
		yearPrefix = "19"
	} else if iin[6] == '5' || iin[6] == '6' {
		yearPrefix = "20"
	}

	year, err := strconv.Atoi(yearPrefix + iin[0:2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year format in IIN: %v", err)
	}

	month, err := strconv.Atoi(iin[2:4])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month format in IIN: %v", err)
	}

	day, err := strconv.Atoi(iin[4:6])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day format in IIN: %v", err)
	}

	birthDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return birthDate, nil
}
