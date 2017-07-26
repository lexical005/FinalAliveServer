package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func get(url string) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Connection", "keep-alive")
	response, _ := client.Do(request)
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
	}
	fmt.Printf("%v", response)
}

func postJSON(url, content string) {
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(content)))
	request.Close = true
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	} else if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
		fmt.Printf("%v", response)
	}
}

func main() {
	// url := "http://115.159.73.243:21040/Activity"
	url := "http://127.0.0.1:21040/Activity"
	content := `{"TimeRangeLimit":{"ActiveTime":{"TimeRange":[],"ActiveType":"EveryWeek"},"TimeZone":-8,"EndTime":"2017-03-31 23:59:59","ResetEveryDay":"true","StartTime":"2017-03-01 00:00:00"},"ServerRoute":"2-Editor","Detail":{"GameRangeLimit":[{"GameSubType":1,"GameMainType":"PVE"},{"GameSubType":2,"GameMainType":"PVE"},{"GameSubType":3,"GameMainType":"PVE"},{"GameSubType":1,"GameMainType":"PVP"}],"NormalShow":{"Content":"Content","ShowImage":"xxxxxx","ButtonParam":"","ButtonType":"Go"},"MultiValue":2},"Description":{"MainTitle":"MainTitle Description","LeftType":"hot","LeftTitle":"LeftTitle Description"},"Activity":{"Action":"Add","Type":"MultiExp","uuid":"MultiExp1"}}`
	postJSON(url, content)
}
