package main

import (
	"flag"
	"log"
	"sort"

	"github.com/jp-stock-analyzer/accounting"
	"github.com/jp-stock-analyzer/loader"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "enable verbose mode")
	flag.Parse()

	acList := make(map[string][]*accounting.Accounting, 0)
	loader.LoadPL("data/2023/fy-profit-and-loss.csv", acList)
	loader.LoadPL("data/2022/fy-profit-and-loss.csv", acList)
	loader.LoadPL("data/2021/fy-profit-and-loss.csv", acList)
	loader.LoadBS("data/2023/fy-balance-sheet.csv", acList)
	loader.LoadBS("data/2022/fy-balance-sheet.csv", acList)
	loader.LoadBS("data/2021/fy-balance-sheet.csv", acList)
	loader.LoadCF("data/2023/fy-cash-flow-statement.csv", acList)
	loader.LoadCF("data/2022/fy-cash-flow-statement.csv", acList)
	loader.LoadCF("data/2021/fy-cash-flow-statement.csv", acList)

	if len(acList) == 0 {
		log.Fatal("acList should not be empty.")
	}

	candidateCode := make([]string, 0, 128)
	for code, acs := range acList {
		if accounting.GoingBankrupt(acs) {
			if verbose {
				log.Printf("code %s is going bankrupt or does not have enough data.", code)
			}
			continue
		}
		if !accounting.IsGrowing(acs) {
			if verbose {
				log.Printf("code %s is not growing.", code)
			}
			continue
		}
		candidateCode = append(candidateCode, code)
	}
	sort.Strings(candidateCode)
	log.Println(candidateCode)
}
