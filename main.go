package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jiseongg/pfm/kis"

	"github.com/joho/godotenv"
)

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
		token = kis.RequestToken(appkey, appsecret)
	}

	// account info
	cano := os.Getenv("CANO")
	acnt_prdt_cd := os.Getenv("ACNT_PRDT_CD")

	jsonRes := kis.RequestBalance(cano, acnt_prdt_cd, appkey, appsecret, token)

	output1 := jsonRes["output1"].([]interface{})[0]
	clear(jsonRes)

	fmt.Printf("Item: %s\n", output1.(map[string]interface{})["ovrs_item_name"])
	fmt.Printf("Holdings: %s\n", output1.(map[string]interface{})["ord_psbl_qty"])
	fmt.Printf("Avg. purchasing price: %s\n", output1.(map[string]interface{})["pchs_avg_pric"])
	fmt.Printf("Current price: %s\n", output1.(map[string]interface{})["now_pric2"])
	fmt.Printf("Profit(%%): %s\n", output1.(map[string]interface{})["evlu_pfls_rt"])
}
