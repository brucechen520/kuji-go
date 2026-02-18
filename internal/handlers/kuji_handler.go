package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 這裡我們假設已經注入了 Service
func GetPrizeList(c *gin.Context) {
	boxID := c.Query("box_id")
	// 呼叫 Service 獲取清單...
	c.JSON(http.StatusOK, gin.H{
		"box_id": boxID,
		"prizes": []string{"A賞-火拳艾斯", "B賞-魯夫"},
	})
}

func Draw(c *gin.Context) {
	// 接收抽獎請求
	c.JSON(http.StatusOK, gin.H{
		"message": "恭喜抽中 A 賞！",
	})
}
