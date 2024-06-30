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

		balancesPerItem := jsonRes["output1"].([]interface{})
		overallPerformance := jsonRes["output2"].(map[string]interface{})
		clear(jsonRes)
		fmt.Println("\n=== Balance ===")
		for i, _balance := range balancesPerItem {
			balance := _balance.(map[string]interface{})
			itemName := balance["ovrs_item_name"]
			ticker := balance["ovrs_pdno"]
			fmt.Printf("%d. %s (%s)\n", i+1, itemName, ticker)
			fmt.Printf(" Holdings: %s\n", balance["ord_psbl_qty"])
			fmt.Printf(" Avg. purchasing price: %s\n", balance["pchs_avg_pric"])
			fmt.Printf(" Current price: %s\n", balance["now_pric2"])
			fmt.Printf(" Profit(%%): %s\n\n", balance["evlu_pfls_rt"])
		}

		fmt.Println("\n=== Performance ===")
		fmt.Printf("Purchasing amount: %s\n", overallPerformance["frcr_pchs_amt1"])
		fmt.Printf("Evaluation amount: %s\n", overallPerformance["tot_evlu_pfls_amt"])
		fmt.Printf("Realized profit/loss ($): %s\n", overallPerformance["ovrs_rlzt_pfls_amt"])
		fmt.Printf("Realized profit/loss (%%): %s\n", overallPerformance["rlzt_erng_rt"])
		fmt.Printf("Total profit ($): %s\n", overallPerformance["ovrs_tot_pfls"])
		fmt.Printf("Total profit (%%): %s\n\n", overallPerformance["tot_pftrt"])
	} else {
		jsonRes := upbit.RequestBalance()
		for i, currency := range jsonRes {
			currencyName := currency["currency"]
			fmt.Printf("%d. %s\n", i+1, currencyName)
			fmt.Printf(" Balance: %s\n", currency["balance"].(string))
			fmt.Printf(" Avg. purchasing price: %s\n\n", currency["avg_buy_price"].(string))
		}
		clear(jsonRes)
	}
}
