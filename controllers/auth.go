package controllers

import (
	"fmt"
	"github.com/ShynaliOtep/paydal-api/models"
	"github.com/ShynaliOtep/paydal-api/services/filesystem"
	"github.com/ShynaliOtep/paydal-api/services/payments"
	"github.com/ShynaliOtep/paydal-api/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

type LoginInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterByPhoneInput struct {
	Phone    string `json:"phone" binding:"required"`
	Iin      string `json:"iin" binding:"required"`
	Fio      string `json:"fio" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       uint               `json:"id"`
	Fio      string             `json:"fio"`
	Phone    string             `json:"phone"`
	Iin      string             `json:"iin"`
	Type     string             `json:"type"`
	Birthday *time.Time         `json:"birthday"`
	Avatar   *FileResponse      `json:"avatar"`
	Wallet   UserWalletResponse `json:"wallet"`
}

type UserWalletResponse struct {
	Balance float64 `json:"balance"`
	Bonus   float64 `json:"bonus"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredDate  time.Time `json:"expired_date"`
	Type         string    `json:"type"`
}

// ByPhoneRegister godoc
// @Summary      Register
// @Description  Register new users
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  RegisterByPhoneInput   true  "Register Data"
// @Success      200  {object}  UserResponse
// @Router       /register [post]
func ByPhoneRegister(c *gin.Context) {
	var input RegisterByPhoneInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := models.GetUserByPhone(input.Phone); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким номером существует "})
		return
	}

	u := models.User{}

	u.Fio = input.Fio
	u.Password = input.Password
	u.Phone = input.Phone
	u.Iin = input.Iin
	birthday, err := models.ExtractBirthDateFromIIN(input.Iin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.Birthday = &birthday

	_, err = u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = payments.TransactionBonus(u.ID, 100.00, payments.Transaction_type_refill)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

// Login godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request  body  LoginInput    true  "Login Data"
// @Success      200  {object}  LoginResponse
// @Router       /login [post]
func Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Phone = input.Phone
	u.Password = input.Password

	accessToken, refreshToken, err := models.LoginCheck(u.Phone, u.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	tokenLifespan, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_HOUR_LIFESPAN"))
	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredDate:  time.Now().Add(time.Hour * time.Duration(tokenLifespan)),
		Type:         "Bearer",
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func RefreshTokenHandler(c *gin.Context) {
	accessToken, refreshToken, err := token.RefreshToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	tokenLifespan, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_HOUR_LIFESPAN"))
	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredDate:  time.Now().Add(time.Hour * time.Duration(tokenLifespan)),
		Type:         "Bearer",
	}
	c.JSON(http.StatusOK, gin.H{"access_token": response})
}

// CurrentUser godoc
// @Summary      Get current app user
// @Description  get string by ID
// @Tags   User Authentication
// @Accept       json
// @Produce      json
// @Success      200  {object}  UserResponse
// @Router       /p/user [get]
func CurrentUser(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)
	fmt.Println(user_id)
	fmt.Println(err)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	fmt.Println(u)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": userConvertToResponse(u)})
}

func userConvertToResponse(user models.User) UserResponse {
	userType := "user"

	_, err := models.GetVolunteerByUserId(user.ID)
	fmt.Println(err)
	if err == nil {
		userType = "volunteer"
	}

	if user.IsPolice {
		userType = "police"
	}

	if user.IsAdmin {
		userType = "admin"
	}

	response := UserResponse{
		ID:       user.ID,
		Fio:      user.Fio,
		Phone:    user.Phone,
		Iin:      user.Iin,
		Type:     userType,
		Birthday: user.Birthday,
		Wallet: UserWalletResponse{
			Balance: user.Wallet.Balance,
			Bonus:   user.Wallet.BonusBalance,
		},
	}

	if user.AvatarId != nil {
		response.Avatar = &FileResponse{
			Url:        filesystem.GerFileUrl(user.Avatar.Path),
			UploadedAt: user.Avatar.UploadedAt,
		}
	}
	return response
}
