package middleware

import (
	"Campus-errand-platform/config"
	"Campus-errand-platform/model"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const CtxUserIDKey = "user_id"

type MyClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func SigntToken(userID uint64, cfg *config.Config) (string, error) {
	claims := MyClaims{
		UserID: userID,
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(cfg.JWTTTL))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
	return tokenStr, err
}

type authStatus uint64

const (
	EmptyHead authStatus = iota
	BadHead
	InvalidToken
	StatusOk
)

func headerCheck(ctx *gin.Context, cfg *config.Config) (*MyClaims, authStatus) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return nil, EmptyHead
	}
	if !strings.HasPrefix(header, "Bearer ") {
		return nil, BadHead
	}
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	claim := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claim, func(t *jwt.Token) (any, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		return nil, InvalidToken
	}
	return claim, StatusOk
}

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claim, status := headerCheck(ctx, cfg)
		if status != StatusOk {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或登录已过期"})
			ctx.Abort()
			return
		}
		ctx.Set(CtxUserIDKey, claim.UserID)
		ctx.Next()
	}
}

func OptionalAuth(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claim, status := headerCheck(ctx, cfg)
		switch status {
		case StatusOk:
			ctx.Set(CtxUserIDKey, claim.UserID)
			ctx.Next()
		case EmptyHead:
			ctx.Next()
		default:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "身份校验失败"})
			ctx.Abort()
		}
	}
}

func RequireRole(role model.Role, cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claim, status := headerCheck(ctx, cfg)
		if status != StatusOk {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或登录已过期"})
			ctx.Abort()
			return
		}

		data := model.Users{}
		result := db.Model(model.Users{}).Where("id=?", claim.UserID).First(&data)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未查询到账户信息"})
				ctx.Abort()
				return
			}
			log.Println("[require_role]db err:", result.Error)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "查询账户信息失败"})
			ctx.Abort()
			return
		}

		if data.Status != model.UserStatusActive {
			ctx.JSON(http.StatusConflict, gin.H{"error": "账户已被封禁"})
			ctx.Abort()
			return
		}
		if data.Role != role {
			ctx.JSON(http.StatusConflict, gin.H{"error": "账户无权限进行该操作"})
			ctx.Abort()
			return
		}
		ctx.Set(CtxUserIDKey, claim.UserID)
		ctx.Next()
	}
}
