package middleware

import (
	"Campus-errand-platform/config"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
