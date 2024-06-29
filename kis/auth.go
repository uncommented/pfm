package kis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func RequestToken(appkey string, appsecret string) string {
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
