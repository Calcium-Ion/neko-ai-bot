package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"neko-ai-bot/conf"
	"net/http"
)

type MidjourneyResponse struct {
	Code        int         `json:"code"`
	Description string      `json:"description"`
	Properties  interface{} `json:"properties"`
	Result      string      `json:"result"`
}

type Midjourney struct {
	MjId        string `json:"id"`
	Action      string `json:"action"`
	Prompt      string `json:"prompt"`
	PromptEn    string `json:"promptEn"`
	Description string `json:"description"`
	State       string `json:"state"`
	SubmitTime  int64  `json:"submitTime"`
	StartTime   int64  `json:"startTime"`
	FinishTime  int64  `json:"finishTime"`
	ImageUrl    string `json:"imageUrl"`
	Status      string `json:"status"`
	Progress    string `json:"progress"`
	FailReason  string `json:"failReason"`
}

func Fetch(taskId string) *Result {
	///mj/task/{id}/fetch
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": conf.Conf.ApiKey,
	}
	url := fmt.Sprintf("https://nekoapi.com/mj/task/%s/fetch", taskId)
	//log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return ErrorUtil("http request error")
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ErrorUtil("http request error")
	}
	defer resp.Body.Close()

	//if resp.StatusCode != 200 {
	//	return ErrorUtil("http status code error")
	//}

	var respData Midjourney
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		log.Println(err)
		return ErrorUtil("json decode error")
	}
	return &Result{
		Code:        200,
		Description: "success",
		Data:        respData,
	}
}

func Change(taskId string, action string, index int) *Result {
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": conf.Conf.ApiKey,
	}
	jsonMap := map[string]interface{}{
		"action": action,
		"index":  index,
		"taskId": taskId,
	}
	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		log.Println(err)
		return ErrorUtil("json marshal error")
	}

	req, err := http.NewRequest("POST", "https://nekoapi.com/mj/submit/change", bytes.NewBuffer(jsonData))
	for key, value := range header {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ErrorUtil("http request error")
	}
	defer resp.Body.Close()

	//if resp.StatusCode != 200 {
	//	return ErrorUtil("http status code error")
	//}

	var respData MidjourneyResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		log.Println(err)
		return ErrorUtil("json decode error")
	}
	return &Result{
		Code:        respData.Code,
		Description: respData.Description,
		Data:        respData.Result,
	}
}

func Imagine(base64 string, prompt string) *Result {
	//url https://nekoapi.com/mj/submit/imagine
	//method POST
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": conf.Conf.ApiKey,
	}

	jsonMap := map[string]interface{}{
		//"base64Array": map[string]interface{}{},
		"prompt": prompt,
	}
	if base64 != "" {
		jsonMap["base64Array"] = []string{base64}
	}

	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		log.Println(err)
		return ErrorUtil("json marshal error")
	}

	req, err := http.NewRequest("POST", "https://nekoapi.com/mj/submit/imagine", bytes.NewBuffer(jsonData))
	for key, value := range header {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ErrorUtil("http request error")
	}
	defer resp.Body.Close()

	//if resp.StatusCode != 200 {
	//	return ErrorUtil("http status code error")
	//}

	var respData MidjourneyResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		log.Println(err)
		return ErrorUtil("json decode error")
	}
	return &Result{
		Code:        respData.Code,
		Description: respData.Description,
		Data:        respData.Result,
	}
}
