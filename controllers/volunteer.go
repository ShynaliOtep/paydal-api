package controllers

import (
	"github.com/ShynaliOtep/paydal-api/models"
	"github.com/ShynaliOtep/paydal-api/services/filesystem"
	"github.com/ShynaliOtep/paydal-api/services/payments"
	"github.com/ShynaliOtep/paydal-api/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type VolunteerContractCreateInput struct {
	Fio    string `form:"fio" binding:"required"`
	CityId uint   `form:"city_id" binding:"required"`
}

type VolunteerContractApproveInput struct {
	ContractId uint `form:"contract_id" binding:"required"`
}

type VolunteerContractResponse struct {
	ID        uint                     `json:"id"`
	UserId    uint                     `json:"user_id"`
	ImageId   uint                     `json:"image_id"`
	Fio       string                   `json:"fio"`
	Document  ContractDocumentResponse `json:"document"`
	Status    string                   `json:"status"`
	Reason    string                   `json:"reason"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type FileResponse struct {
	Url        string    `json:"url"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type ContractDocumentResponse struct {
	FrontImage      FileResponse `json:"front_image"`
	BackImage       FileResponse `json:"bac_image"`
	FrontWithAvatar FileResponse `json:"front_with_avatar"`
	BackWithAvatar  FileResponse `json:"back_with_avatar"`
}

func VolunteerContractCreate(c *gin.Context) {
	var input VolunteerContractCreateInput
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract := models.VolunteerContract{}

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := models.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract.UserId = user.ID

	frontImage, frontImageHeader, err := c.Request.FormFile("front_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	backImage, backImageHeader, err := c.Request.FormFile("back_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	frontWithAvatar, frontWithAvatarHeader, err := c.Request.FormFile("front_with_avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	backWithAvatar, backImageAvatarHeader, err := c.Request.FormFile("back_with_avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	imagePath, err := filesystem.UploadFile(filesystem.GeneratePath("front_image", filesystem.GetExtension(frontImageHeader.Filename)), frontImage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	frontImageFile, err := models.SaveFile(imagePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	imagePath, err = filesystem.UploadFile(filesystem.GeneratePath("image", filesystem.GetExtension(backImageHeader.Filename)), backImage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	backImageFile, err := models.SaveFile(imagePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	imagePath, err = filesystem.UploadFile(filesystem.GeneratePath("front_with_avatar", filesystem.GetExtension(frontWithAvatarHeader.Filename)), frontWithAvatar)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	frontWithAvatarFile, err := models.SaveFile(imagePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	imagePath, err = filesystem.UploadFile(filesystem.GeneratePath("back_with_avatar", filesystem.GetExtension(backImageAvatarHeader.Filename)), backWithAvatar)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	backWithAvatarFile, err := models.SaveFile(imagePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	doc := models.VolunteerContractDocument{}
	doc.FrontImageId = frontImageFile.ID
	doc.BackImageId = backImageFile.ID
	doc.FrontWithAvatarId = frontWithAvatarFile.ID
	doc.BackWithAvatarId = backWithAvatarFile.ID

	_, err = doc.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contract.DocumentId = doc.ID
	contract.Fio = input.Fio
	contract.CityId = input.CityId

	_, err = contract.Create()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = payments.TransactionBonus(contract.UserId, 1000.00, payments.Transaction_type_refill)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Volunteer contract success", "data": nil})
}

func VolunteerContractApprove(c *gin.Context) {
	var input VolunteerContractApproveInput

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract, err := models.GetVolunteerContractById(input.ContractId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract.Status = "confirmed"
	_, err = contract.Save()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volunteer := models.Volunteer{
		UserId:     contract.UserId,
		ContractId: contract.ID,
	}

	_, err = volunteer.Create()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = payments.TransactionBonus(volunteer.UserId, 1000.00, payments.Transaction_type_refill)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": contractConvertToResponse(contract)})
}

func VolunteerContractStatus(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract, err := models.GetVolunteerContractByUserId(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": contractConvertToResponse(contract)})
}

func contractConvertToResponse(contract models.VolunteerContract) VolunteerContractResponse {
	return VolunteerContractResponse{
		ID:     contract.ID,
		UserId: contract.UserId,
		Fio:    contract.Fio,
		Document: ContractDocumentResponse{
			FrontImage: FileResponse{
				Url:        filesystem.GerFileUrl(contract.Document.FrontImage.Path),
				UploadedAt: contract.Document.FrontImage.UploadedAt,
			},
			BackImage: FileResponse{
				Url:        filesystem.GerFileUrl(contract.Document.BackImage.Path),
				UploadedAt: contract.Document.BackImage.UploadedAt,
			},
			FrontWithAvatar: FileResponse{
				Url:        filesystem.GerFileUrl(contract.Document.FrontWithAvatar.Path),
				UploadedAt: contract.Document.FrontImage.UploadedAt,
			},
			BackWithAvatar: FileResponse{
				Url:        filesystem.GerFileUrl(contract.Document.BackWithAvatar.Path),
				UploadedAt: contract.Document.BackImage.UploadedAt,
			},
		},
		Status:    contract.Status,
		Reason:    contract.Reason,
		CreatedAt: contract.CreatedAt,
		UpdatedAt: contract.UpdatedAt,
	}
}
