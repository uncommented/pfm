package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jiseongg/pfm/kis"
	"github.com/jiseongg/pfm/upbit"
	"github.com/joho/godotenv"
)

var vendor string

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	flag.StringVar(&vendor, "vendor", "kis", "use 'kis' or 'upbit'")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if vendor == "kis" {
		jsonRes := kis.RequestBalance()

		output1 := jsonRes["output1"].([]interface{})[0]
		clear(jsonRes)

		fmt.Printf("Item: %s\n", output1.(map[string]interface{})["ovrs_item_name"])
		fmt.Printf("Holdings: %s\n", output1.(map[string]interface{})["ord_psbl_qty"])
		fmt.Printf("Avg. purchasing price: %s\n", output1.(map[string]interface{})["pchs_avg_pric"])
		fmt.Printf("Current price: %s\n", output1.(map[string]interface{})["now_pric2"])
		fmt.Printf("Profit(%%): %s\n", output1.(map[string]interface{})["evlu_pfls_rt"])
	} else {
		jsonRes := upbit.RequestBalance()
		for _, currency := range jsonRes {
			fmt.Println("---")
			fmt.Printf("Currency: %s\n", currency["currency"].(string))
			fmt.Printf("Balance: %s\n", currency["balance"].(string))
			fmt.Printf("Avg. purchasing price: %s\n", currency["avg_buy_price"].(string))
		}
		clear(jsonRes)
	}
}
