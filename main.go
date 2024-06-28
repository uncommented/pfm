package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const TR_ID = "TTTS3012R"

const CANO = ""
const ACNT_PRDT_CD = ""
const APP_KEY = ""
const APP_SECRET = ""

func main() {
	client := &http.Client{}

	// Request token
	body := []byte(fmt.Sprintf(`{
		"appkey": "%s",
		"appsecret": "%s",
		"grant_type": "client_credentials"
	}`, APP_KEY, APP_SECRET))

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
	accessToken := jsonRes["access_token"]
	clear(jsonRes)

	// Request balance
	req, err = http.NewRequest("GET", "https://openapi.koreainvestment.com:9443/uapi/overseas-stock/v1/trading/inquire-balance", nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("appkey", APP_KEY)
	req.Header.Add("appsecret", APP_SECRET)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("tr_id", TR_ID)
	q := req.URL.Query()
	q.Add("CANO", CANO)
	q.Add("ACNT_PRDT_CD", ACNT_PRDT_CD)
	q.Add("OVRS_EXCG_CD", "NASD")
	q.Add("TR_CRCY_CD", "USD")
	q.Add("CTX_AREA_FK200", "")
	q.Add("CTX_AREA_NK200", "")
	req.URL.RawQuery = q.Encode()

	res, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &jsonRes)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	output1 := jsonRes["output1"].([]interface{})[0]

	fmt.Printf("Item: %s\n", output1.(map[string]interface{})["ovrs_item_name"])
	fmt.Printf("Holdings: %s\n", output1.(map[string]interface{})["ord_psbl_qty"])
	fmt.Printf("Avg. purchasing price: %s\n", output1.(map[string]interface{})["pchs_avg_pric"])
	fmt.Printf("Current price: %s\n", output1.(map[string]interface{})["now_pric2"])
	fmt.Printf("Profit(%%): %s\n", output1.(map[string]interface{})["evlu_pfls_rt"])
	clear(jsonRes)
}
