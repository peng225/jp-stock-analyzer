package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"accounting"
)

const (
	// Common
	CodeIndex = 0
	DateIndex = 1

	// BL
	JikoshihonRatioIndex = 9

	// PS
	EigyoProfitIndex = 3
	KeijoProfitIndex = 4
	NetProftIndex    = 5
	ROEIndex         = 7
	ROAIndex         = 8
)

func main() {
	log.SetFlags(log.Lshortfile)

	acList := make(map[string][]*accounting.Accounting, 0)

	// Read BS
	bsFile, err := os.Open("fy-balance-sheet.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer bsFile.Close()

	r := csv.NewReader(bsFile)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		code := row[CodeIndex]
		dateStr := row[DateIndex]
		date, err := time.Parse("1234/56", dateStr)
		if err != nil {
			log.Fatal(err.Error())
		}
		jsRatio, err := strconv.ParseFloat(row[JikoshihonRatioIndex], 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		newBS := NewBalanceSheet(jsRatio)
		newAc := NewAccounting(date, nil, newBS, nil)
		if _, ok := acList[code]; ok {
			acList[code] = []*Accounting{
				newAc,
			}
		} else {
			acList[code] = append(acList[code], newAc)
		}
	}

	// Read PS
	plFile, err := os.Open("fy-profit-and-loss.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer plFile.Close()

	r = csv.NewReader(plFile)
	rows, err = r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		code := row[CodeIndex]
		dateStr := row[DateIndex]
		date, err := time.Parse("1234/56", dateStr)
		if err != nil {
			log.Fatal(err.Error())
		}
		eigyoProfit, err := strconv.ParseInt(row[EigyoProfitIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		keijoProfit, err := strconv.ParseInt(row[KeijoProfitIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		netProfit, err := strconv.ParseInt(row[NetProftIndex], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		roe, err := strconv.ParseFloat(row[ROEIndex], 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		roa, err := strconv.ParseFloat(row[ROAIndex], 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		newPL := NewProfitLoss(eigyoProfit, keijoProfit,
			netProfit, roe, roa)
		if _, ok := acList[code]; ok {
			log.Fatalf("code %s not found.", code)
		}
		found := false
		foundIndex := -1
		for i, ac := range acList[code] {
			if ac.Date.Year() == date.Year() &&
				ac.Date.Month() == date.Month() {
				foundIndex = i
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("date %s not found.", date)
		}
		acList[code][foundIndex].PL = newPL
	}
}
