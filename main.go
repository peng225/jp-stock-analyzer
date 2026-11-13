package main

import (
	"flag"
	"log"
	"sort"

	"github.com/peng225/jp-stock-analyzer/accounting"
	"github.com/peng225/jp-stock-analyzer/loader"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var code string
	var verbose bool
	flag.StringVar(&code, "c", "", "stock code")
	flag.BoolVar(&verbose, "v", false, "enable verbose mode")
	flag.Parse()

	fileNames := []struct {
		pl string
		bs string
		cf string
		sd string
	}{
		{
			pl: "data/2025/fy-profit-and-loss.csv",
			bs: "data/2025/fy-balance-sheet.csv",
			cf: "data/2025/fy-cash-flow-statement.csv",
			sd: "data/2025/fy-stock-dividend.csv",
		},
		{
			pl: "data/2024/fy-profit-and-loss.csv",
			bs: "data/2024/fy-balance-sheet.csv",
			cf: "data/2024/fy-cash-flow-statement.csv",
			sd: "data/2024/fy-stock-dividend.csv",
		},
		{
			pl: "data/2023/fy-profit-and-loss.csv",
			bs: "data/2023/fy-balance-sheet.csv",
			cf: "data/2023/fy-cash-flow-statement.csv",
			sd: "data/2023/fy-stock-dividend.csv",
		},
	}

	acList := make(map[string][]*accounting.Accounting, 0)
	for _, fileName := range fileNames {
		loader.LoadPL(fileName.pl, acList)
		loader.LoadBS(fileName.bs, acList)
		loader.LoadCF(fileName.cf, acList)
		loader.LoadStockDividend(fileName.sd, acList)

	}

	if len(acList) == 0 {
		log.Fatal("acList should not be empty.")
	}

	if code != "" {
		acs := acList[code]
		if accounting.Risky(acs, verbose) {
			log.Printf("code %s is going bankrupt or does not have enough data.", code)
		} else {
			log.Printf("code %s is safe.", code)
		}
		accounting.Evaluate(acs)
		return
	}

	candidateCode := make([]string, 0, 128)
	for code, acs := range acList {
		if accounting.Risky(acs, verbose) {
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
