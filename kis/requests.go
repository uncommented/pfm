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
)

const (
	BASE_URL = "https://openapi.koreainvestment.com:9443"
	TR_ID    = "TTTS3012R"
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

	req, err := http.NewRequest("POST", BASE_URL+"/oauth2/tokenP", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

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

	os.Setenv("KIS_TOKEN", token)
	os.Setenv("KIS_TOKEN_EXPIRED", token_expired)
}

func prepareToken() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	token := os.Getenv("KIS_TOKEN")
	token_expired, err := time.ParseInLocation(time.DateTime, os.Getenv("KIS_TOKEN_EXPIRED"), loc)

	if err != nil || time.Now().After(token_expired) || token == "" {
		log.Println("Token is expired! Request another one!")
		requestToken()
	}
}

func RequestBalance(accountNumber string, currency Currency, marketCode MarketCode) map[string]interface{} {
	prepareToken()

	client := &http.Client{}

	accountNumberSplits := strings.Split(accountNumber, "-")
	cano := accountNumberSplits[0]
	acnt_prdt_cd := accountNumberSplits[1]

	token := os.Getenv("KIS_TOKEN")
	appkey := os.Getenv("KIS_APPKEY")
	appsecret := os.Getenv("KIS_APPSECRET")

	req, err := http.NewRequest("GET", BASE_URL+"/uapi/overseas-stock/v1/trading/inquire-balance", nil)
	if err != nil {
		log.Fatal(err)
		return nil
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
		log.Fatal(err)
		return nil
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal(data, &jsonRes)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return jsonRes
}
