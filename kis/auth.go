package kis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func requestToken() {
	client := &http.Client{}

	appkey := os.Getenv("KIS_APPKEY")
	appsecret := os.Getenv("KIS_APPSECRET")

	body := []byte(fmt.Sprintf(`{
		"appkey": "%s",
		"appsecret": "%s",
		"grant_type": "client_credentials"
	}`, appkey, appsecret))

	req, err := http.NewRequest("POST", "https://openapi.koreainvestment.com:9443/oauth2/tokenP", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	var jsonRes map[string]interface{}
	err = json.Unmarshal(data, &jsonRes)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	token := jsonRes["access_token"].(string)
	token_expired := jsonRes["access_token_token_expired"].(string)
	clear(jsonRes)

	os.Setenv("KIS_TOKEN", token)
	os.Setenv("KIS_TOKEN_EXPIRED", token_expired)

	log.Println("Update .env with following new token information")
	log.Printf("KIS_TOKEN=%s\n", token)
	log.Printf("KIS_TOKEN_EXPIRED=%s\n", token_expired)
}

func PrepareToken() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	token_expired, err := time.ParseInLocation(time.DateTime, os.Getenv("KIS_TOKEN_EXPIRED"), loc)
	token := os.Getenv("KIS_TOKEN")

	if err != nil || time.Now().After(token_expired) || token == "" {
		log.Println("Token is expired! Request another one!")
		requestToken()
	}
}
