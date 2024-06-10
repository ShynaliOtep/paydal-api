package controllers

import (
	"fmt"
	"github.com/ShynaliOtep/paydal-api/models"
	"github.com/ShynaliOtep/paydal-api/services/filesystem"
	"github.com/ShynaliOtep/paydal-api/utils/token"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ViolationApplicationCreateInput struct {
	Lat        float64   `json:"lat" binding:"required"`
	Lon        float64   `json:"lon" binding:"required"`
	CategoryId uint      `json:"category_id" binding:"required"`
	Phone      string    `json:"phone" binding:"required"`
	Address    string    `json:"address" binding:"required"`
	Comment    string    `json:"comment"`
	CityId     uint      `json:"city_id"`
	Date       time.Time `json:"date" binding:"required"`
}

type ViolationsInput struct {
	Lat  float64 `form:"lat" binding:"required"`
	Lon  float64 `form:"lon" binding:"required"`
	Dist float64 `form:"dist" binding:"required"`
}
type ViolationUploadMediaInput struct {
	ViolationId uint `form:"violation_id" binding:"required"`
}

type ViolationRejectInput struct {
	ViolationId uint   `json:"violation_id" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
}

type ViolationResponse struct {
	ID          uint                      `json:"id"`
	VolunteerId uint                      `json:"volunteer_id"`
	Category    ViolationCategoryResponse `json:"category"`
	Location    LocationResponse          `json:"location"`
	Phone       string                    `json:"phone"`
	Address     string                    `json:"address"`
	Comment     string                    `json:"comment"`
	City        CityResponse              `json:"city"`
	Status      string                    `json:"status"`
	Reason      string                    `json:"reason"`
	Date        time.Time                 `json:"date"`
	Medias      []ViolationMediaResponse  `json:"medias"`
}

type ViolationCategoryResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type LocationResponse struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type ViolationMediaResponse struct {
	ID          uint         `json:"id"`
	ViolationId uint         `json:"violation_id"`
	File        FileResponse `json:"file"`
}

func GetViolationCategories(c *gin.Context) {
	categories, err := models.GetViolationCategories()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": categories})
}

func ViolationApplicationCreate(c *gin.Context) {

	var input ViolationApplicationCreateInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("F")
	violation := models.Violation{}

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	volunteer, err := models.GetVolunteerByUserId(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	violation.VolunteerId = volunteer.ID
	// категорияға валидация қосу керек
	violation.CategoryId = input.CategoryId
	violation.Address = input.Address
	violation.Comment = input.Comment
	violation.Phone = input.Phone
	violation.Date = input.Date
	violation.Status = "created"
	violation.CityId = input.CityId

	location := models.Location{
		Lat: input.Lat,
		Lon: input.Lon,
	}
	_, err = location.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	violation.LocationId = location.ID

	_, err = violation.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": violationConvertToResponse(violation)})
}

func ViolationUploadMedia(c *gin.Context) {
	var input ViolationUploadMediaInput
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	media := models.ViolationMedia{}

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = models.GetVolunteerByUserId(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	media.ViolationId = input.ViolationId

	file, fileHeader, err := c.Request.FormFile("media")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	imagePath, err := filesystem.UploadFile(filesystem.GeneratePath("violation", filesystem.GetExtension(fileHeader.Filename)), file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	mediaFile, err := models.SaveFile(imagePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	media.FileId = mediaFile.ID

	_, err = media.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": nil})
}

func GetViolation(c *gin.Context) {

	number, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert string to uint: %v", err)
	}
	id := uint(number)

	violation, err := models.GetViolationById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": violationConvertToResponse(violation)})
}

func GetViolations(c *gin.Context) {
	var input ViolationsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	violations, err := models.ViolationSearchInMap(input.Lat, input.Lon, input.Dist)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": violationsConvertToResponse(violations)})
}

func GetViolationsForPolice(c *gin.Context) {
	violations, err := models.ViolationsForPolice()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": violationsConvertToResponse(violations)})
}

func GetViolationForPolice(c *gin.Context) {
	number, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert string to uint: %v", err)
	}
	id := uint(number)

	violation, err := models.GetViolationById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if violation.Status == "created" {
		violation.Status = "waiting"
		_, err = violation.Save()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": violationConvertToResponse(violation)})
}

func ViolationReject(c *gin.Context) {
	var input ViolationRejectInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	violation, err := models.GetViolationById(input.ViolationId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	violation.Status = "rejected"
	violation.Reason = input.Reason
	_, err = violation.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": nil})
}

func GerVolunteerViolations(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volunteer, err := models.GetVolunteerByUserId(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	violations, err := models.VolunteerViolations(volunteer.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": violationsConvertToResponse(violations)})
}

func violationsConvertToResponse(violations []models.Violation) []ViolationResponse {
	var response []ViolationResponse
	for _, violation := range violations {
		response = append(response, violationConvertToResponse(violation))
	}
	return response
}

func violationConvertToResponse(violation models.Violation) ViolationResponse {
	response := ViolationResponse{
		ID:          violation.ID,
		VolunteerId: violation.VolunteerId,
		Location: LocationResponse{
			Lat: violation.Location.Lat,
			Lon: violation.Location.Lon,
		},
		Category: ViolationCategoryResponse{
			ID:    violation.Category.ID,
			Title: violation.Category.Title,
		},
		City: CityResponse{
			ID:    violation.City.ID,
			Title: violation.City.Title,
		},
		Address: violation.Address,
		Comment: violation.Comment,
		Phone:   violation.Phone,
		Date:    violation.Date,
		Status:  violation.Status,
		Reason:  violation.Reason,
	}
	var mediasResponse []ViolationMediaResponse
	for _, media := range violation.Medias {
		mediasResponse = append(mediasResponse, convertViolationMediaResponse(media))
	}
	response.Medias = mediasResponse
	return response
}

func convertViolationMediaResponse(media models.ViolationMedia) ViolationMediaResponse {
	return ViolationMediaResponse{
		ID:          media.ID,
		ViolationId: media.ViolationId,
		File: FileResponse{
			Url:        filesystem.GerFileUrl(media.File.Path),
			UploadedAt: media.File.UploadedAt,
		},
	}
}
