package utils

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil
}

func GetLoginFromHeader(ctx *gin.Context) (string, error) {
	reqBasic := ctx.Request.Header.Get("Authorization")
	splitBasic := strings.Split(reqBasic, "Basic ")
	reqBasic = splitBasic[1]
	decodedHeader, err := base64.StdEncoding.DecodeString(reqBasic)
	if err != nil {
		return "", err
	}
	splitDecodedHeader := strings.Split(string(decodedHeader), ":")
	return splitDecodedHeader[0], nil
}

func GetPasswordFromHeader(ctx *gin.Context) (string, error) {
	reqBasic := ctx.Request.Header.Get("Authorization")
	splitBasic := strings.Split(reqBasic, "Basic ")
	reqBasic = splitBasic[1]
	decodedHeader, err := base64.StdEncoding.DecodeString(reqBasic)
	if err != nil {
		return "", err
	}
	splitDecodedHeader := strings.Split(string(decodedHeader), ":")
	return splitDecodedHeader[1], nil
}