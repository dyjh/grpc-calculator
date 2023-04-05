package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/dyjh/grpc_calculator/calculator"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := calculator.NewCalculatorClient(conn)

	r := gin.Default()

	// 路由1：计算操作
	r.GET("/calculate/:operation/:num1/:num2", func(ctx *gin.Context) {
		num1, err := strconv.ParseFloat(ctx.Param("num1"), 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter num1"})
			return
		}

		num2, err := strconv.ParseFloat(ctx.Param("num2"), 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter num2"})
			return
		}

		op := ctx.Param("operation")
		var operation calculator.Operation
		switch op {
		case "add":
			operation = calculator.Operation_ADD
		case "subtract":
			operation = calculator.Operation_SUBTRACT
		case "multiply":
			operation = calculator.Operation_MULTIPLY
		case "divide":
			operation = calculator.Operation_DIVIDE
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation"})
			return
		}

		calculateReq := &calculator.CalculateRequest{
			Num1:      float32(num1),
			Num2:      float32(num2),
			Operation: operation,
		}

		calculateRes, err := c.Calculate(context.Background(), calculateReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": calculateRes.GetResult()})
	})

	// 路由2：
	r.GET("/compare/:num1/:num2", func(ctx *gin.Context) {
		num1, err := strconv.ParseFloat(ctx.Param("num1"), 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter num1"})
			return
		}
		num2, err := strconv.ParseFloat(ctx.Param("num2"), 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter num2"})
			return
		}

		compareReq := &calculator.CompareRequest{
			Num1: float32(num1),
			Num2: float32(num2),
		}

		compareRes, err := c.Compare(context.Background(), compareReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"max": compareRes.GetMax()})
	})

	r.Run(":8080")
}
