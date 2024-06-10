package controllers

import (
	"github.com/ShynaliOtep/paydal-api/models"
	"github.com/ShynaliOtep/paydal-api/services/filesystem"
	"github.com/ShynaliOtep/paydal-api/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserUploadAvatar(c *gin.Context) {

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

	file, fileHeader, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	imagePath, err := filesystem.UploadFile(filesystem.GeneratePath("user/avatar", filesystem.GetExtension(fileHeader.Filename)), file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	avatarFile, err := models.SaveFile(imagePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user.AvatarId = &avatarFile.ID

	_, err = user.UpdateUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": nil})
}
