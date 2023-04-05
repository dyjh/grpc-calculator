package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CalculateRequest struct {
	Num1      float32 `json:"num1"`
	Num2      float32 `json:"num2"`
	Operation string  `json:"operation"`
}

type CalculateResponse struct {
	Result float32 `json:"result"`
}

func main() {
	// 调用 Calculate 方法
	calculateReq := CalculateRequest{
		Num1:      10,
		Num2:      5,
		Operation: "add",
	}

	jsonData, err := json.Marshal(calculateReq)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// 请根据您的 Kong 配置修改 URL
	resp, err := http.Post("http://localhost:8000/calculate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var calculateRes CalculateResponse
	err = json.Unmarshal(body, &calculateRes)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	fmt.Printf("Result of calculation: %v\n", calculateRes.Result)
}
