package controllers

import (
	"github.com/ShynaliOtep/paydal-api/models"
	"github.com/ShynaliOtep/paydal-api/services/payments"
	"github.com/ShynaliOtep/paydal-api/utils/mrp"
	"github.com/ShynaliOtep/paydal-api/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ProtocolCreateInput struct {
	ViolationId uint                        `json:"violation_id" binding:"required"`
	ArticleId   uint                        `json:"article_id" binding:"required"`
	Violator    ProtocolCreateViolatorInput `json:"violator" binding:"required"`
	Car         ProtocolCreateCarInput      `json:"car" binding:"required"`
	Mrp         int                         `json:"mrp" binding:"required"`
}

type ProtocolCreateViolatorInput struct {
	Iin      string    `json:"iin" binding:"required"`
	Address  string    `json:"address" binding:"required"`
	Birthday time.Time `json:"birthday" time_format:"2006-01-02" binding:"required"`
}

type ProtocolCreateCarInput struct {
	Number   string `json:"number"`
	Passport string `json:"passport"`
	Mark     string `json:"mark"`
	Color    string `json:"color"`
}

type ProtocolSearchInput struct {
	Iin string `form:"iin" binding:"required"`
}

type ProtocolResponse struct {
	ID          uint                     `json:"id"`
	VolunteerId uint                     `json:"volunteer_id"`
	PoliceId    uint                     `json:"police_id"`
	Violation   ViolationResponse        `json:"violation"`
	Article     ArticleResponse          `json:"article"`
	Violator    ProtocolViolatorResponse `json:"violator"`
	Car         ProtocolCarResponse      `json:"car"`
	Mrp         int                      `json:"mrp"`
	Price       int                      `json:"price"`
	Status      string                   `json:"status"`
}

type ProtocolViolatorResponse struct {
	Iin      string    `json:"iin"`
	Address  string    `json:"address"`
	Birthday time.Time `json:"birthday"`
}

type ProtocolCarResponse struct {
	Number   string `json:"number"`
	Passport string `json:"passport"`
	Mark     string `json:"mark"`
	Color    string `json:"color"`
}

func ProtocolCreate(c *gin.Context) {
	var input ProtocolCreateInput
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	protocol := models.Protocol{}

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	police, err := models.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	protocol.PoliceId = police.ID
	protocol.ArticleId = input.ArticleId
	protocol.ViolationId = input.ViolationId
	protocol.Mrp = input.Mrp
	protocol.Price = *mrp.GetMrpPrice(&input.Mrp)

	violator := models.Violator{
		UserId:   nil,
		Iin:      input.Violator.Iin,
		Address:  input.Violator.Address,
		Birthday: input.Violator.Birthday,
	}
	_, err = violator.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	protocol.ViolatorId = violator.ID

	violation, err := models.GetViolationById(input.ViolationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if violation.CategoryId == 1 {
		car := models.Car{
			Number:   input.Car.Number,
			Passport: input.Car.Passport,
			Mark:     input.Car.Mark,
			Color:    input.Car.Color,
		}
		_, err = car.Create()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		protocol.CarId = &car.ID
	} else {
		protocol.CarId = nil
	}

	protocol.Status = "created"

	_, err = protocol.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	violation, err = models.GetViolationById(protocol.ViolationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	violation.Status = "confirmed"
	_, err = violation.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = payments.TransactionBalance(police.ID, 1000.00, payments.Transaction_type_refill)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": protocolConvertToResponse(protocol)})
}

func ProtocolSearch(c *gin.Context) {
	var input ProtocolSearchInput
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	protocols, err := models.GetProtocolsByIin(input.Iin)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": protocolsConvertToResponse(protocols)})
}

func GetPoliceProtocols(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	police, err := models.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	protocols, err := models.GetPoliceProtocols(police.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": protocolsConvertToResponse(protocols)})
}

func protocolsConvertToResponse(protocols []models.Protocol) []ProtocolResponse {
	var response []ProtocolResponse
	for _, protocol := range protocols {
		response = append(response, protocolConvertToResponse(protocol))
	}
	return response
}

func protocolConvertToResponse(protocol models.Protocol) ProtocolResponse {
	response := ProtocolResponse{
		ID:          protocol.ID,
		VolunteerId: protocol.Violation.VolunteerId,
		PoliceId:    protocol.PoliceId,
		Violation:   violationConvertToResponse(protocol.Violation),
		Article: ArticleResponse{
			ID:          protocol.Article.ID,
			Number:      protocol.Article.Number,
			Title:       protocol.Article.Title,
			Description: protocol.Article.Description,
			Mrp:         protocol.Article.Mrp,
			Price:       mrp.GetMrpPrice(protocol.Article.Mrp),
		},
		Violator: ProtocolViolatorResponse{
			Iin:      protocol.Violator.Iin,
			Address:  protocol.Violator.Address,
			Birthday: protocol.Violator.Birthday,
		},
		Mrp:    protocol.Mrp,
		Price:  protocol.Price,
		Status: protocol.Status,
	}

	if protocol.Violation.CategoryId == 1 {
		response.Car = ProtocolCarResponse{
			Number:   protocol.Car.Number,
			Passport: protocol.Car.Passport,
			Mark:     protocol.Car.Mark,
			Color:    protocol.Car.Color,
		}
	}
	return response
}
