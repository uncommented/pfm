package upbit

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const UPBIT_URL = "https://api.upbit.com"

func RequestToken(query string) (string, error) {
	var token *jwt.Token

	accesskey := os.Getenv("UPBIT_ACCESS_KEY")
	secretkey := os.Getenv("UPBIT_SECRET_KEY")

	if query == "" {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"access_key": accesskey,
				"nonce":      uuid.New().String(),
			})
	} else {
		sha_512 := sha512.New()
		sha_512.Write([]byte(query))
		query_hash := fmt.Sprintf("%x", sha_512.Sum(nil))
		token = jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"access_key":     accesskey,
				"nonce":          uuid.New().String(),
				"query_hash":     query_hash,
				"query_hash_alg": "SHA512",
			})
	}

	return token.SignedString([]byte(secretkey))
}

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
