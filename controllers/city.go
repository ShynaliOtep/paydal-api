package controllers

import (
	"github.com/ShynaliOtep/paydal-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CityResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func GetCities(c *gin.Context) {
	cities, err := models.GetCities()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cities})
}
