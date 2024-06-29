package upbit

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func RequestBalance() []map[string]interface{} {
	client := &http.Client{}

	// Request balance
	req, err := http.NewRequest("GET", "https://api.upbit.com/v1/accounts", nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	token, err := RequestToken("")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

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

	var jsonRes []map[string]interface{}
	err = json.Unmarshal(data, &jsonRes)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return jsonRes
}
