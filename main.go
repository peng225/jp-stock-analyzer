package main

import (
	"log"

	"github.com/jp-stock-analyzer/accounting"
	"github.com/jp-stock-analyzer/loader"
)

func main() {
	log.SetFlags(log.Lshortfile)

	acList := make(map[string][]*accounting.Accounting, 0)
	bsLoader := loader.NewBSLoader("fy-balance-sheet.csv")
	bsLoader.Load(acList)

	plLoader := loader.NewPLLoader("fy-profit-and-loss.csv")
	plLoader.Load(acList)

	cfLoader := loader.NewCFLoader("fy-cash-flow-statement.csv")
	cfLoader.Load(acList)

	log.Println(acList["9914"][0].PL)
	log.Println(acList["9914"][0].BS)
	log.Println(acList["9914"][0].CF)
}
