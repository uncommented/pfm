package kis

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const TR_ID = "TTTS3012R"

func RequestBalance() map[string]interface{} {
	PrepareToken()

	client := &http.Client{}

	cano := os.Getenv("CANO")
	acnt_prdt_cd := os.Getenv("ACNT_PRDT_CD")
	token := os.Getenv("KIS_TOKEN")
	appkey := os.Getenv("KIS_APPKEY")
	appsecret := os.Getenv("KIS_APPSECRET")

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

	return jsonRes
}
