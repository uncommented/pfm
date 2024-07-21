package kis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/uncommented/pfm/portfolio/utils"
)

const (
	BASE_URL = "https://openapi.koreainvestment.com:9443"
	TR_ID    = "TTTS3012R"
)

func requestToken() {
	log.Println("Request token...")

	client := &http.Client{}

	appkey := os.Getenv("KIS_APPKEY")
	appsecret := os.Getenv("KIS_APPSECRET")

	body := []byte(fmt.Sprintf(`{
		"appkey": "%s",
		"appsecret": "%s",
		"grant_type": "client_credentials"
	}`, appkey, appsecret))

	req, err := http.NewRequest("POST", BASE_URL+"/oauth2/tokenP", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to make POST request for token: %v", err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to request token: %v", err)
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		return
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal(data, &jsonRes)
	if err != nil {
		log.Printf("Failed to unmarshal data to json: %v", err)
		return
	}

	token := utils.UnmarshalToString(jsonRes, "access_token")
	token_expired := utils.UnmarshalToString(jsonRes, "access_token_token_expired")

	os.Setenv("KIS_TOKEN", token)
	os.Setenv("KIS_TOKEN_EXPIRED", token_expired)
}

func prepareToken() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Printf("Failed to load location: %v", err)
		return
	}

	token := os.Getenv("KIS_TOKEN")
	if token == "" {
		log.Println("Token is not prepared!")
		requestToken()
	} else {
		token_expired, err := time.ParseInLocation(time.DateTime, os.Getenv("KIS_TOKEN_EXPIRED"), loc)
		if err != nil {
			log.Println("Failed to parse token expiration date!")
			requestToken()
		} else if time.Now().After(token_expired) {
			log.Println("Token is expired!")
			requestToken()
		}
	}
}

func RequestBalance(currency Currency, marketCode MarketCode) map[string]interface{} {
	prepareToken()

	client := &http.Client{}

	accountNumber := os.Getenv("KIS_ACCOUNT_NUMBER")
	accountNumberSplits := strings.Split(accountNumber, "-")
	if len(accountNumberSplits) != 2 {
		log.Printf("Invalid account number format: %s", accountNumber)
		return make(map[string]interface{})
	}
	cano := accountNumberSplits[0]
	acnt_prdt_cd := accountNumberSplits[1]

	token := os.Getenv("KIS_TOKEN")
	appkey := os.Getenv("KIS_APPKEY")
	appsecret := os.Getenv("KIS_APPSECRET")

	req, err := http.NewRequest("GET", BASE_URL+"/uapi/overseas-stock/v1/trading/inquire-balance", nil)
	if err != nil {
		log.Printf("Failed to make GET request: %v", err)
		return make(map[string]interface{})
	}

	req.Header.Add("appkey", appkey)
	req.Header.Add("appsecret", appsecret)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("tr_id", TR_ID)
	q := req.URL.Query()
	q.Add("CANO", cano)
	q.Add("ACNT_PRDT_CD", acnt_prdt_cd)
	q.Add("OVRS_EXCG_CD", marketCode.String())
	q.Add("TR_CRCY_CD", currency.String())
	q.Add("CTX_AREA_FK200", "")
	q.Add("CTX_AREA_NK200", "")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to request: %v", err)
		return make(map[string]interface{})
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		return make(map[string]interface{})
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal(data, &jsonRes)
	if err != nil {
		log.Printf("Failed to unmarshal data to json: %v", err)
		return make(map[string]interface{})
	}

	return jsonRes
}
