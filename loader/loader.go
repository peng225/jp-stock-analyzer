package loader

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jp-stock-analyzer/accounting"
)

const (
	// Common
	codeIndex = 0
	dateIndex = 1

	// BL
	profitJouyoMoneyIndex      = 5
	shortTermKariireMoenyIndex = 6
	longTermKariireMoenyIndex  = 7
	jikoshihonRatioIndex       = 9

	// PS
	eigyoProfitIndex = 3
	keijoProfitIndex = 4
	netProftIndex    = 5
	roeIndex         = 7
	roaIndex         = 8

	// CF
	eigyoCFIndex           = 2
	genkinDoutouButsuIndex = 6
)

func LoadPL(fileName string, acList map[string][]*accounting.Accounting) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		code := row[codeIndex]
		dateStr := row[dateIndex]
		date, err := time.Parse("2006/01", dateStr)
		if err != nil {
			log.Fatal(err.Error())
		}
		eigyoProfit, err := strconv.ParseInt(row[eigyoProfitIndex], 10, 64)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		keijoProfit, err := strconv.ParseInt(row[keijoProfitIndex], 10, 64)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		netProfit, err := strconv.ParseInt(row[netProftIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		roe, err := strconv.ParseFloat(row[roeIndex], 64)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		roa, err := strconv.ParseFloat(row[roaIndex], 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		newPL := accounting.NewProfitLoss(eigyoProfit, keijoProfit,
			netProfit, roe, roa)
		if _, ok := acList[code]; !ok {
			newAc := accounting.NewAccounting(date, newPL, nil, nil)
			acList[code] = []*accounting.Accounting{
				newAc,
			}
		} else {
			found := false
			for _, v := range acList[code] {
				if v.Date.Year() == date.Year() &&
					v.Date.Month() == date.Month() {
					v.PL = newPL
					found = true
					break
				}
			}
			if !found {
				newAc := accounting.NewAccounting(date, newPL, nil, nil)
				acList[code] = append(acList[code], newAc)
			}
		}
	}
}

func LoadBS(fileName string, acList map[string][]*accounting.Accounting) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		code := row[codeIndex]
		dateStr := row[dateIndex]
		date, err := time.Parse("2006/01", dateStr)
		if err != nil {
			log.Fatal(err.Error())
		}
		pjMoney, err := strconv.ParseInt(row[profitJouyoMoneyIndex], 10, 64)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		stkMoney, err := strconv.ParseInt(row[shortTermKariireMoenyIndex], 10, 64)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		ltkMoney, err := strconv.ParseInt(row[longTermKariireMoenyIndex], 10, 64)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		jsRatio, err := strconv.ParseFloat(row[jikoshihonRatioIndex], 64)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		newBS := accounting.NewBalanceSheet(pjMoney, stkMoney, ltkMoney, jsRatio)
		if _, ok := acList[code]; !ok {
			newAc := accounting.NewAccounting(date, nil, newBS, nil)
			acList[code] = []*accounting.Accounting{
				newAc,
			}
		} else {
			found := false
			for _, v := range acList[code] {
				if v.Date.Year() == date.Year() &&
					v.Date.Month() == date.Month() {
					v.BS = newBS
					found = true
					break
				}
			}
			if !found {
				newAc := accounting.NewAccounting(date, nil, newBS, nil)
				acList[code] = append(acList[code], newAc)
			}
		}
	}
}

func LoadCF(fileName string, acList map[string][]*accounting.Accounting) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		code := row[codeIndex]
		dateStr := row[dateIndex]
		date, err := time.Parse("2006/01", dateStr)
		if err != nil {
			log.Fatal(err.Error())
		}
		eigyoCF, err := strconv.ParseInt(row[eigyoCFIndex], 10, 64)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		gdb, err := strconv.ParseInt(row[genkinDoutouButsuIndex], 10, 64)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		newCF := accounting.NewCashFlow(eigyoCF, gdb)
		if _, ok := acList[code]; !ok {
			newAc := accounting.NewAccounting(date, nil, nil, newCF)
			acList[code] = []*accounting.Accounting{
				newAc,
			}
		} else {
			found := false
			for _, v := range acList[code] {
				if v.Date.Year() == date.Year() &&
					v.Date.Month() == date.Month() {
					v.CF = newCF
					found = true
					break
				}
			}
			if !found {
				newAc := accounting.NewAccounting(date, nil, nil, newCF)
				acList[code] = append(acList[code], newAc)
			}
		}
	}
}
