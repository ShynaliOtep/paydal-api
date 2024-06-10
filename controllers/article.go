package controllers

import (
	"github.com/ShynaliOtep/paydal-api/models"
	"github.com/ShynaliOtep/paydal-api/utils/mrp"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ArticleResponse struct {
	ID          uint   `json:"id"`
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Mrp         *int   `json:"mrp"`
	Price       *int   `json:"price"`
}

func ArticleSearch(c *gin.Context) {
	articles, err := models.GetArticles(c.Query("search"), c.Query("price_sort"))
	var response []ArticleResponse
	for _, article := range articles {
		response = append(response, ArticleResponse{
			ID:          article.ID,
			Number:      article.Number,
			Title:       article.Title,
			Description: article.Description,
			Mrp:         article.Mrp,
			Price:       mrp.GetMrpPrice(article.Mrp),
		})
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": response})
}
