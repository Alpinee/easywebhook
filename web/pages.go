/**
* @author: gongquanlin
* @since: 2024/3/23
* @desc:
 */

package web

import (
	"easywebhook/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	query := map[string]int{}
	err := c.BindQuery(query)
	if err != nil {
		// c.AbortWithStatus(http.StatusInternalServerError, gin.H{"error":"Invaliad "})
	}

	pageSize := query["pageSize"]
	pageNum := query["pageSize"]

	if pageSize == 0 {
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	var pageOffset int
	pageOffset = (pageNum - 1) * pageSize

	var results []domain.TokenScript

	userId := c.GetInt("userId")
	domain.Db.Where(domain.TokenScript{UserId: userId}).Offset(pageOffset).Limit(pageSize).Find(&results)

	var count int64
	domain.Db.Model(&domain.TokenScript{}).Where(domain.TokenScript{UserId: userId}).Count(&count)

	c.HTML(http.StatusOK, "index", gin.H{
		"list":  results,
		"total": count,
	})
}
