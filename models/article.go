package models

import (
	"fmt"
	"time"
)

type Article struct {
	ID          uint   `json:"id"`
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Mrp         *int   `json:"mrp"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func GetArticles(search string, priceSort string) ([]Article, error) {
	var articles []Article
	var err error
	if priceSort == "" {
		priceSort = "desc"
	}

	if search != "" {
		err = DB.Order(fmt.Sprintf("mrp %s", priceSort)).Where("title LIKE '%" + search + "%' OR description LIKE '%" + search + "%'").Find(&articles).Error
	} else {
		err = DB.Order(fmt.Sprintf("mrp %s", priceSort)).Find(&articles).Error
	}

	if err != nil {
		return nil, err
	}

	return articles, nil
}
