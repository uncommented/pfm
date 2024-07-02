package kis

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	BASE_URL = "https://openapi.koreainvestment.com:9443"
	TR_ID    = "TTTS3012R"
)

func RequestBalance(accountNumber string, currency string) map[string]interface{} {
	PrepareToken()

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

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("appkey", appkey)
	req.Header.Add("appsecret", appsecret)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("tr_id", TR_ID)
	q := req.URL.Query()
	q.Add("CANO", cano)
	q.Add("ACNT_PRDT_CD", acnt_prdt_cd)
	q.Add("OVRS_EXCG_CD", "NASD")
	q.Add("TR_CRCY_CD", currency)
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
