package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CalculateRequest struct {
	Num1      float32 `json:"num1"`
	Num2      float32 `json:"num2"`
	Operation string  `json:"operation"`
}

type CompareRequest struct {
	Num1 float32 `json:"num1"`
	Num2 float32 `json:"num2"`
}

func main() {
	r := gin.Default()

	baseURL := "http://localhost:8000"

	// 路由1：计算操作
	r.POST("/calculate", func(ctx *gin.Context) {
		var req CalculateRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp, err := http.Post(baseURL+"/calculate", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		var result map[string]interface{}
		json.Unmarshal(body, &result)

		ctx.JSON(http.StatusOK, result)
	})

	// 路由2：比较两个参数的大小并返回较大值
	r.POST("/compare", func(ctx *gin.Context) {
		var req CompareRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp, err := http.Post(baseURL+"/compare", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		var result map[string]interface{}
		json.Unmarshal(body, &result)

		ctx.JSON(http.StatusOK, result)
	})

	r.Run(":8080")
}
