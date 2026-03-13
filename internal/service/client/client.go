package client

import (
	"bytes"
	"encoding/json"
	"festival_greeting/internal/service/config"
	"fmt"
	"io"
	"net/http"
)

type APIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

type RequestBody struct {
	Model       string              `json:"model"`
	Messages    []map[string]string `json:"messages"`
	Temperature float64             `json:"temperature"`
}

type ResponseBody struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

type APIClient struct {
	Client *http.Client
}

func NewClient(model config.Model) *APIClient {
	return &APIClient{
		Client: &http.Client{},
	}
}

func (c *APIClient) GetResponse(prompt string, model config.Model) (string, error) {
	requestBody := RequestBody{
		Model: model.ModelName,
		Messages: []map[string]string{
			{"role": "user", "content": prompt},
		},
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化为json数据失败 %w", err)
	}

	request, err := http.NewRequest("POST", model.BaseUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+model.ApiKey)
	fmt.Printf("发送请求到模型: %s\n", model.BaseUrl)
	fmt.Printf("请求内容: %s\n", string(jsonData))
	fmt.Printf("请求头: Content-Type=%s, Authorization=%s\n", request.Header.Get("Content-Type"), "Bearer ***")
	fmt.Printf("请求体: %s\n", string(jsonData))
	response, err := c.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("请求失败,请检查模型链接和apikey %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败 %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var apiErr APIErrorResponse
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error.Message != "" {
			return "", fmt.Errorf("模型接口返回错误(%d): %s", response.StatusCode, apiErr.Error.Message)
		}
		return "", fmt.Errorf("模型接口返回错误(%d): %s", response.StatusCode, string(body))
	}

	var responseBody ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
		fmt.Printf("HTTP状态码: %d\n", response.StatusCode)
		fmt.Printf("HTTP状态: %s\n", response.Status)
		fmt.Printf("JSON解析失败，原始响应: %s\n", string(body))
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(responseBody.Choices) == 0 {
		fmt.Printf("模型响应无choices，原始响应: %s\n", string(body))
		return "", fmt.Errorf("模型响应格式异常: choices为空")
	}
	return responseBody.Choices[0].Message.Content, nil

}
