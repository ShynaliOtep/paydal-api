package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateToken(user_id uint) (string, error) {

	token_lifespan, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(c *gin.Context) (uint, error) {

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}

func GenerateTokenPair(user_id uint) (accessToken string, refreshToken string, err error) {
	// Генерируем access token
	accessToken, err = generateAccessToken(user_id)
	if err != nil {
		return "", "", err
	}

	// Генерируем refresh token
	refreshToken, err = generateRefreshToken(user_id)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateAccessToken(user_id uint) (string, error) {
	tokenLifespan, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
}

func generateRefreshToken(user_id uint) (string, error) {
	refreshTokenLifespan, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(refreshTokenLifespan)).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
}

func RefreshToken(c *gin.Context) (string, string, error) {
	refreshTokenString := ExtractRefreshToken(c)
	user_id, err := ExtractTokenIDFromRefreshToken(refreshTokenString)
	fmt.Println(err)
	if err != nil {
		return "", "", err
	}

	// Генерируем новую пару токенов
	accessToken, refreshToken, err := GenerateTokenPair(user_id)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ExtractRefreshToken(c *gin.Context) string {
	refreshToken := c.Query("refresh_token")
	if refreshToken != "" {
		return refreshToken
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenIDFromRefreshToken(refreshTokenString string) (uint, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}
