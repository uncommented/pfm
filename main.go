package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const TR_ID = "TTTS3012R"

func requestToken(appkey string, appsecret string) string {
	client := &http.Client{}

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

	new_token_info := make(map[string]string)
	new_token_info["TOKEN"] = token
	new_token_info["TOKEN_EXPIRED"] = token_expired

	new_token_info_str, err := godotenv.Marshal(new_token_info)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("Update .env with following new token information")
	fmt.Printf("\n%s\n\n", new_token_info_str)

	return token
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// secrets
	token_expired, err := time.Parse(time.DateTime, os.Getenv("TOKEN_EXPIRED"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	appkey := os.Getenv("APPKEY")
	appsecret := os.Getenv("APPSECRET")

	var token string
	if time.Now().Before(token_expired) {
		token = os.Getenv("TOKEN")
	} else {
		log.Println("Token is expired! Request another one!")
		token = requestToken(appkey, appsecret)
	}

	// account info
	cano := os.Getenv("CANO")
	acnt_prdt_cd := os.Getenv("ACNT_PRDT_CD")

	client := &http.Client{}

	// Request balance
	req, err := http.NewRequest("GET", "https://openapi.koreainvestment.com:9443/uapi/overseas-stock/v1/trading/inquire-balance", nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("appkey", appkey)
	req.Header.Add("appsecret", appsecret)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("tr_id", TR_ID)
	q := req.URL.Query()
	q.Add("CANO", cano)
	q.Add("ACNT_PRDT_CD", acnt_prdt_cd)
	q.Add("OVRS_EXCG_CD", "NASD")
	q.Add("TR_CRCY_CD", "USD")
	q.Add("CTX_AREA_FK200", "")
	q.Add("CTX_AREA_NK200", "")
	req.URL.RawQuery = q.Encode()

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
	output1 := jsonRes["output1"].([]interface{})[0]
	clear(jsonRes)

	fmt.Printf("Item: %s\n", output1.(map[string]interface{})["ovrs_item_name"])
	fmt.Printf("Holdings: %s\n", output1.(map[string]interface{})["ord_psbl_qty"])
	fmt.Printf("Avg. purchasing price: %s\n", output1.(map[string]interface{})["pchs_avg_pric"])
	fmt.Printf("Current price: %s\n", output1.(map[string]interface{})["now_pric2"])
	fmt.Printf("Profit(%%): %s\n", output1.(map[string]interface{})["evlu_pfls_rt"])
}
