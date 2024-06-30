package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// author:QQ1831921455
// 抓取微信小程序 或 顺丰app
// 请求的URL
const url = "https://mcs-mimp-web.sf-express.com/mcs-mimp/commonPost/~memberNonactivity~memberDayFreeService~freeCouponPurchase"

// Cookie
var cookies = []string{

	"",
	"",
	"",
	"",
}

// header
var headersTemplate = map[string]string{
	"Host":            "mcs-mimp-web.sf-express.com",
	"Accept":          "application/json, text/plain, */*",
	"channel":         "26memapp",
	"sysCode":         "MCS-MIMP-CORE",
	"sw8":             "1-YzAxNjg4ZWYtZjQ2ZC00YjRlLWE1NTAtYTVmMTQ2ZjU2MjA1-ZTNlOTY5MDgtM2ZlOS00NWQxLWJhMTEtMjU3YWUwYzFjNGJj-0-ZmI0MDgxNzA4NWJlNGUzOThlMGI2ZjRiMDgxNzc3NDY=-d2Vi-L29yaWdpbi9hL21pbXAtYWN0aXZpdHkvbWVtYmVyRGF5-L21jcy1taW1wL2NvbW1vblBvc3Qvfm1lbWJlck5vbmFjdGl2aXR5fm1lbWJlckRheUZyZWVTZXJ2aWNlfmZyZWVDb3Vwb25QdXJjaGFzZQ==",
	"timestamp":       "1719558000",
	"Accept-Language": "zh-CN,zh-Hans;q=0.9",
	"Accept-Encoding": "gzip, deflate, br",
	"platform":        "SFAPP",
	"signature":       "ad9840eb2e70fb1798d75da2851c5764",
	"Origin":          "https://mcs-mimp-web.sf-express.com",
	"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 16_7_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 mediaCode=SFEXPRESSAPP-iOS-ML",
	"Referer":         "https://mcs-mimp-web.sf-express.com/origin/a/mimp-activity/memberDay?mobile=187****0480&userId=9D2D05BA609249B0A810476634434723&path=/memberDayV2&linkCode=SFAC20230413163415517&supportShare=YES&from=26memapp",
	"Content-Length":  "21",
	"Connection":      "keep-alive",
	"deviceId":        "D2TVQg2KROlgQUwLOxTBTUDRoLJTxF8Y40oYqpdjCCUncXc4",
	"Sec-Fetch-Dest":  "empty",
	"Sec-Fetch-Site":  "same-origin",
	"Sec-Fetch-Mode":  "cors",
	"Content-Type":    "application/json",
}

// 抢卷时间
var payload = map[string]string{
	"roundTime": "15:00",
	//"roundTime": "12:00",
	//"roundTime": "09:00",
}

type Response struct {
	Cookie       string
	StatusCode   int
	ResponseBody map[string]interface{}
}

// 发送请求
func sendRequest(cookie string, wg *sync.WaitGroup, responses *[]Response, mu *sync.Mutex) {
	defer wg.Done()

	headers := make(http.Header)
	for key, value := range headersTemplate {
		headers.Set(key, value)
	}
	headers.Set("Cookie", cookie)

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	mu.Lock()
	*responses = append(*responses, Response{
		Cookie:       cookie,
		StatusCode:   resp.StatusCode,
		ResponseBody: responseBody,
	})
	mu.Unlock()
}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	responses := []Response{}

	// 并发次数
	concurrency := 20
	// 每次并发的请求数量
	requestsPerConcurrency := 40

	// 设定的目标时间，格式为小时:分钟:秒
	targetTime := "14:59:45"
	//targetTime := "08:59:45"
	//targetTime := "11:59:45"

	// 等待直到设定的时间
	for {
		currentTime := time.Now().Format("15:04:05")
		fmt.Printf("Current time is %s\n", currentTime)
		if currentTime == targetTime {
			fmt.Printf("Current time is %s. Starting requests.\n", currentTime)
			break
		}
		time.Sleep(time.Second)
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(len(cookies) * requestsPerConcurrency)
		for j := 0; j < requestsPerConcurrency; j++ {
			for _, cookie := range cookies {
				go sendRequest(cookie, &wg, &responses, &mu)
			}
		}
		wg.Wait()
		currentTime := time.Now().Format("15:04:05")
		fmt.Printf("发送请求 Current time is %s\n", currentTime)
	}

	// 打印抢卷结果
	for i, res := range responses {
		fmt.Printf("Request %d:\n", i+1)
		fmt.Println("Cookie:", res.Cookie)
		fmt.Println("Status Code:", res.StatusCode)
		fmt.Println("Response Body:", res.ResponseBody)
		fmt.Println(strings.Repeat("-", 60))
	}
}
