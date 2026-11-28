package handler

import (
	"Campus-errand-platform/Http"
	"Campus-errand-platform/model"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddOrderCategories(ctx *gin.Context, db *gorm.DB) {
	var req Http.AddOrderCategoriesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := model.OrderCategories{
		SchoolID:    req.SchoolID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconUrl,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActivite,
	}

	var check_code model.OrderCategories
	codechecker := 0
	_result := db.Model(&model.OrderCategories{}).Where("school_id=? and code=?", data.SchoolID, data.Code).First(&check_code)
	if _result.Error != nil {
		if errors.Is(_result.Error, gorm.ErrRecordNotFound) {
			codechecker = 1
		}
		log.Println("[AddOrderCategories]db err:", _result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单分类失败"})
		return
	}

	if codechecker != 1 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "该学校内订单分类code重复"})
		return
	}

	result := db.Model(&model.OrderCategories{}).Create(&data)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "添加订单类型失败"})
		return
	}

	res := Http.AddOrderCategoriesResponse{
		Data: data,
	}
	ctx.JSON(http.StatusOK, res)
}

func ModOrderCategories(ctx *gin.Context, db *gorm.DB) {
	var req Http.ModOrderCategoriesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idstr := ctx.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单分类id"})
		return
	}

	data := model.OrderCategories{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconUrl,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActivite,
	}

	result := db.Model(&model.OrderCategories{}).Where("id=?", id).Select("name,description,icon_url,sort_order,is_active").Updates(&data)
	if result.Error != nil {
		log.Println("[mod_article]db:", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "修改订单类型失败"})
		return
	}

	res := Http.ModOrderCategoriesResponse{
		Data: data,
	}
	ctx.JSON(http.StatusOK, res)
}
