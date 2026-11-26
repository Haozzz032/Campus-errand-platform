package handler

import (
	"Campus-errand-platform/Http"
	"Campus-errand-platform/model"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOrdercategories(ctx *gin.Context, db *gorm.DB) {
	var req Http.GetOrdercategoriesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data []model.OrderCategories
	result := db.Model(&model.OrderCategories{}).Where("school_id=? and IsActive=?", req.SchoolID, req.ActiveOnly).Order("sort_order ASC").Find(&data)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "当前学校暂无开放订单类型"})
			return
		}
		log.Println("[GetOrdercategories]db err:", result.Error)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "查询订单类型失败"})
		return
	}

	res := Http.GetOrdercategoriesResponse{
		Data: data,
	}

	ctx.JSON(http.StatusOK, res)
}
