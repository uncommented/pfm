package upbit

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const UPBIT_URL = "https://api.upbit.com"

func RequestBalance() []map[string]interface{} {
	req, err := http.NewRequest("GET", UPBIT_URL+"/v1/accounts", nil)
	if err != nil {
		log.Fatal(err)
		return []map[string]interface{}{}
	}

	token, err := RequestToken("")
	if err != nil {
		log.Fatal(err)
		return []map[string]interface{}{}
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return []map[string]interface{}{}
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return []map[string]interface{}{}
	}

	var jsonRes []map[string]interface{}
	err = json.Unmarshal(data, &jsonRes)
	if err != nil {
		log.Fatal(err)
		return []map[string]interface{}{}
	}
	return jsonRes
}

func RequestMarketSnapshot(ticker string) map[string]interface{} {
	req, err := http.NewRequest("GET", UPBIT_URL+"/v1/ticker", nil)
	if err != nil {
		log.Fatal(err)
		return make(map[string]interface{})
	}

	query := req.URL.Query()
	marketCode := "KRW-" + ticker
	query.Add("markets", marketCode)
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return make(map[string]interface{})
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return make(map[string]interface{})
	}

	var jsonRes []map[string]interface{}
	err = json.Unmarshal(data, &jsonRes)
	if err != nil {
		log.Fatal(err)
		return make(map[string]interface{})
	} else if len(jsonRes) != 1 {
		return make(map[string]interface{})
	}
	return jsonRes[0]
}

func RequestMarketInfo(ticker string) map[string]interface{} {
	req, err := http.NewRequest("GET", UPBIT_URL+"/v1/market/all", nil)
	if err != nil {
		log.Fatal(err)
		return make(map[string]interface{})
	}

	marketCode := "KRW-" + ticker

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return make(map[string]interface{})
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return make(map[string]interface{})
	}

	var jsonRes []map[string]interface{}
	err = json.Unmarshal(data, &jsonRes)
	if err != nil {
		log.Fatal(err)
		return make(map[string]interface{})
	}

	for _, marketInfo := range jsonRes {
		if marketInfo["market"].(string) == marketCode {
			return marketInfo
		}
	}

	return make(map[string]interface{})
}
