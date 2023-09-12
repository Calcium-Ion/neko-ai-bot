package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func DoRequest(baseUrl string, path string, token string) (map[string]interface{}, error) {
	uri := fmt.Sprintf(baseUrl + path)
	client := &http.Client{}
	// 创建请求
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("NewRequest Error:", err)
		return nil, err
	}

	// 添加请求头
	req.Header.Add("Authorization", "Bearer "+token)

	log.Println(req.Header)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Do Error:", err)
		return nil, err
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll Error:", err)
		return nil, err
	}
	// resp body {"data":{"quota":0},"message":"","success":true}
	respBody := make(map[string]interface{})
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		fmt.Println("Unmarshal Error:", err)
		return nil, err
	}
	return respBody, nil
}
