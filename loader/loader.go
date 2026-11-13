package loader

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/peng225/jp-stock-analyzer/accounting"
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
	revenueIndex     = 2
	eigyoProfitIndex = 3
	keijoProfitIndex = 4
	netProfitIndex   = 5
	roeIndex         = 7
	roaIndex         = 8

	// CF
	eigyoCFIndex           = 2
	genkinDoutouButsuIndex = 6

	// Stock dividend
	haitouSeikouIndex = 5
)

func isBlank(data string) bool {
	return data == "-"
}

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

		if isBlank(row[revenueIndex]) {
			continue
		}
		revenue, err := strconv.ParseInt(row[revenueIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		if isBlank(row[eigyoProfitIndex]) {
			continue
		}
		eigyoProfit, err := strconv.ParseInt(row[eigyoProfitIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		if isBlank(row[keijoProfitIndex]) {
			continue
		}
		keijoProfit, err := strconv.ParseInt(row[keijoProfitIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		if isBlank(row[netProfitIndex]) {
			continue
		}
		netProfit, err := strconv.ParseInt(row[netProfitIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		if isBlank(row[roeIndex]) {
			continue
		}
		roe, err := strconv.ParseFloat(row[roeIndex], 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		if isBlank(row[roaIndex]) {
			continue
		}
		roa, err := strconv.ParseFloat(row[roaIndex], 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		newPL := accounting.NewProfitLoss(revenue, eigyoProfit, keijoProfit,
			netProfit, roe, roa)
		if _, ok := acList[code]; !ok {
			newAc := accounting.NewAccounting(date, newPL, nil, nil, nil)
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
				newAc := accounting.NewAccounting(date, newPL, nil, nil, nil)
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

		if isBlank(row[profitJouyoMoneyIndex]) {
			continue
		}
		pjMoney, err := strconv.ParseInt(row[profitJouyoMoneyIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		stkMoney, err := strconv.ParseInt(row[shortTermKariireMoenyIndex], 10, 64)
		if err != nil {
			if isBlank(row[shortTermKariireMoenyIndex]) {
				stkMoney = 0
			} else {
				log.Fatal(err.Error())
			}
		}

		ltkMoney, err := strconv.ParseInt(row[longTermKariireMoenyIndex], 10, 64)
		if err != nil {
			if isBlank(row[longTermKariireMoenyIndex]) {
				ltkMoney = 0
			} else {
				log.Fatal(err.Error())
			}
		}

		if isBlank(row[jikoshihonRatioIndex]) {
			continue
		}
		jsRatio, err := strconv.ParseFloat(row[jikoshihonRatioIndex], 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		newBS := accounting.NewBalanceSheet(pjMoney, stkMoney, ltkMoney, jsRatio)
		if _, ok := acList[code]; !ok {
			newAc := accounting.NewAccounting(date, nil, newBS, nil, nil)
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
				newAc := accounting.NewAccounting(date, nil, newBS, nil, nil)
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

		if isBlank(row[eigyoCFIndex]) {
			continue
		}
		eigyoCF, err := strconv.ParseInt(row[eigyoCFIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		if isBlank(row[genkinDoutouButsuIndex]) {
			continue
		}
		gdb, err := strconv.ParseInt(row[genkinDoutouButsuIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		newCF := accounting.NewCashFlow(eigyoCF, gdb)
		if _, ok := acList[code]; !ok {
			newAc := accounting.NewAccounting(date, nil, nil, newCF, nil)
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
				newAc := accounting.NewAccounting(date, nil, nil, newCF, nil)
				acList[code] = append(acList[code], newAc)
			}
		}
	}
}

func LoadStockDividend(fileName string, acList map[string][]*accounting.Accounting) {
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

		if isBlank(row[haitouSeikouIndex]) {
			continue
		}
		haitouSeikou, err := strconv.ParseFloat(row[haitouSeikouIndex], 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		newSD := accounting.NewStockDividend(haitouSeikou)
		if _, ok := acList[code]; !ok {
			newAc := accounting.NewAccounting(date, nil, nil, nil, newSD)
			acList[code] = []*accounting.Accounting{
				newAc,
			}
		} else {
			found := false
			for _, v := range acList[code] {
				if v.Date.Year() == date.Year() &&
					v.Date.Month() == date.Month() {
					v.SD = newSD
					found = true
					break
				}
			}
			if !found {
				newAc := accounting.NewAccounting(date, nil, nil, nil, newSD)
				acList[code] = append(acList[code], newAc)
			}
		}
	}
}
