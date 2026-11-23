package handler

import (
	"Campus-errand-platform/config"
	Http "Campus-errand-platform/Http"
	"Campus-errand-platform/middleware"
	"Campus-errand-platform/model"
	"errors"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Regist(ctx *gin.Context, db *gorm.DB, cfg *config.Config) {
	var req Http.RegistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Username = strings.TrimSpace(req.Username)

	data := model.Users{
		Username:     req.Username,
		PasswordHash: req.Password,
		SchoolID:     &req.SchoolID,
		Phone:        req.Phone,
		Status:       model.UserStatusActive,
	}

	switch req.Role {
	case "publisher":
		data.Role = model.RoleUser
	case "runner":
		data.Role = model.RoleRunner
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户注册身份错误"})
		return
	}

	if utf8.RuneCountInString(data.Username) < 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名至少3位"})
		return
	}
	if utf8.RuneCountInString(data.Username) > 50 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能超过50位"})
		return
	}
	if data.PasswordHash == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
		return
	}
	if len(data.PasswordHash) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码不能小于8位"})
		return
	}
	if len(data.PasswordHash) > 72 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码不能超过72位"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[regist]hash err:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	data.PasswordHash = string(hash)

	result := db.Model(&model.Users{}).Create(&data)
	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			ctx.JSON(http.StatusConflict, gin.H{"error": "用户名重复"})
			return
		}
		log.Println("[regist]db err:", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "注册用户失败"})
		return
	}

	tokenStr, err := middleware.SigntToken(data.ID, cfg)
	if err != nil {
		log.Println("[regist]jwt err:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}

	res := Http.RegistResponse{
		AccessToken: tokenStr,
		TokenType:   "Bearer",
		ExpiresIn:   604800,
		User:        data,
	}

	ctx.JSON(http.StatusCreated, res)
}

func Login(ctx *gin.Context, db *gorm.DB, cfg *config.Config) {
	var req Http.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "登录请求参数错误"})
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名与密码不能为空"})
		return
	}

	var data model.Users
	if err := db.Model(&model.Users{}).Where("username=?", req.Username).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		log.Println("[login]db err:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(data.PasswordHash), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	
	if data.Status != model.UserStatusActive {
		ctx.JSON(http.StatusConflict, gin.H{"error": "账户已被封禁"})
		return
	}

	

	tokenStr, err := middleware.SigntToken(data.ID, cfg)
	if err != nil {
		log.Println("[login]jwt err:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}

	res := Http.LoginResponse{
		AccessToken: tokenStr,
		TokenType:   "Bearer",
		ExpiresIn:   604800,
		User:        data,
	}
	ctx.JSON(http.StatusOK, res)
}
