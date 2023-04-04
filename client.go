package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/dyjh/grpc_calculator/calculator"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorClient(conn)

	r := gin.Default()

	r.POST("/calculate", func(c *gin.Context) {
		num1Str := c.PostForm("num1")
		num2Str := c.PostForm("num2")
		operation := c.PostForm("operation")

		num1, err := strconv.ParseFloat(num1Str, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid num1"})
			return
		}

		num2, err := strconv.ParseFloat(num2Str, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid num2"})
			return
		}

		op := pb.Operation_ADD
		switch operation {
		case "add":
			op = pb.Operation_ADD
		case "subtract":
			op = pb.Operation_SUBTRACT
		case "multiply":
			op = pb.Operation_MULTIPLY
		case "divide":
			op = pb.Operation_DIVIDE
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation"})
			return
		}

		req := &pb.CalculateRequest{Num1: float32(num1), Num2: float32(num2), Operation: op}
		res, err := client.Calculate(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to call Calculate: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": res.Result})
	})

	r.Run(":8080")
}
