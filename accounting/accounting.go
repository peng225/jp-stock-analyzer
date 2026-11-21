package accounting

import (
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Accounting struct {
	Date time.Time
	PL   *ProfitLoss
	BS   *BalanceSheet
	CF   *CashFlow
	SD   *StockDividend
}

type ProfitLoss struct {
	revenue         int64
	operatingProfit int64
	ordinaryProfit  int64
	netProfit       int64
	roe             float64
	roa             float64
}

type BalanceSheet struct {
	retainedEarnings       int64
	shortTermBorrowedMoney int64
	longTermBorrowedMoney  int64
	equityRatio            float64
}

type CashFlow struct {
	operatingCF    int64
	cashEquivalent int64
}

type StockDividend struct {
	dividendPayoutRatio float64
}

func NewAccounting(date time.Time, pl *ProfitLoss,
	bs *BalanceSheet, cf *CashFlow, sd *StockDividend) *Accounting {
	return &Accounting{
		Date: date,
		PL:   pl,
		BS:   bs,
		CF:   cf,
		SD:   sd,
	}
}

func NewProfitLoss(revenue, operatingProfit, ordinaryProfit,
	netProfit int64, roe, roa float64) *ProfitLoss {
	return &ProfitLoss{
		revenue:         revenue,
		operatingProfit: operatingProfit,
		ordinaryProfit:  ordinaryProfit,
		netProfit:       netProfit,
		roe:             roe,
		roa:             roa,
	}
}

func NewBalanceSheet(retainedEarnings, stkMoney, ltkMoney int64,
	equityRatio float64) *BalanceSheet {
	return &BalanceSheet{
		retainedEarnings:       retainedEarnings,
		shortTermBorrowedMoney: stkMoney,
		longTermBorrowedMoney:  ltkMoney,
		equityRatio:            equityRatio,
	}
}

func NewCashFlow(operatingCF, gdb int64) *CashFlow {
	return &CashFlow{
		operatingCF:    operatingCF,
		cashEquivalent: gdb,
	}
}

func NewStockDividend(hs float64) *StockDividend {
	return &StockDividend{
		dividendPayoutRatio: hs,
	}
}

func Risky(acs []*Accounting) (risky bool, unknown bool) {
	if len(acs) == 0 {
		slog.Info("no data")
		return false, true
	}

	sort.Slice(acs, func(i, j int) bool {
		return acs[i].Date.Before(acs[j].Date)
	})

	for _, ac := range acs {
		if ac.PL == nil {
			slog.Info("PL is nil")
			unknown = true
		} else {
			if ac.PL.netProfit <= 0 {
				slog.Info("Net profit is negative.")
				risky = true
			}
			if ac.PL.roe/ac.PL.roa >= 3.0 {
				slog.Info("ROE/ROA ratio is too large.")
				risky = true
			}
		}
		if ac.BS == nil {
			slog.Info("BS is nil")
			unknown = true
		} else {
			if ac.BS.retainedEarnings <= 0 {
				slog.Info("Retained earnings is negative.")
				risky = true
			}
			if ac.BS.equityRatio < 20 {
				slog.Info("Equity ratio is too small.")
				risky = true
			}
		}
		if ac.CF == nil {
			slog.Info("CF is nil")
			unknown = true
		} else {
			if ac.CF.operatingCF <= 0 {
				slog.Info("Operating cache flow is negative.")
				risky = true
			}
			if ac.BS != nil {
				if ac.CF.cashEquivalent < ac.BS.shortTermBorrowedMoney+ac.BS.longTermBorrowedMoney {
					slog.Info("Borrowed money is larger than cash equivalent.")
					risky = true
				}
			}
		}
		if ac.SD == nil {
			slog.Info("SD is nil")
			unknown = true
		} else {
			if ac.SD.dividendPayoutRatio <= 20.0 || 50.0 <= ac.SD.dividendPayoutRatio {
				slog.Info("Dividend payout ratio is not in the proper range.", "dividendPayoutRatio", ac.SD.dividendPayoutRatio)
				risky = true
			}
		}
	}
	return
}

func IsGrowing(acs []*Accounting, strict bool) bool {
	sort.Slice(acs, func(i, j int) bool {
		return acs[i].Date.Before(acs[j].Date)
	})
	if len(acs) < 2 {
		return false
	}
	if strict && len(acs) < 3 {
		return false
	}
	for i := 0; i < len(acs)-1; i++ {
		// Check the growth of the ordinary profit.
		if acs[i].PL == nil || acs[i+1].PL == nil {
			return false
		}
		if acs[i].PL.ordinaryProfit >= acs[i+1].PL.ordinaryProfit {
			return false
		}
		if acs[i].PL.ordinaryProfit < 0 || acs[i+1].PL.ordinaryProfit < 0 {
			return false
		}
		if float64(acs[i+1].PL.ordinaryProfit)/float64(acs[i].PL.ordinaryProfit) < 1.05 {
			return false
		}

		// Check the growth of the revenue.
		if acs[i].PL == nil || acs[i+1].PL == nil {
			return false
		}
		if acs[i].PL.revenue >= acs[i+1].PL.revenue {
			return false
		}
		if acs[i].PL.revenue < 0 || acs[i+1].PL.revenue < 0 {
			return false
		}
		if float64(acs[i+1].PL.revenue)/float64(acs[i].PL.revenue) < 1.05 {
			return false
		}
	}

	for _, ac := range acs {
		if ac.PL.roe < 10 {
			return false
		}
		if ac.PL.roa < 5 {
			return false
		}
	}
	return true
}

func Evaluate(acs []*Accounting) {
	sort.Slice(acs, func(i, j int) bool {
		return acs[i].Date.Before(acs[j].Date)
	})
	if len(acs) < 2 {
		slog.Info("no enough data.")
		return
	}

	for i, ac := range acs {
		if ac.PL == nil {
			slog.Info(fmt.Sprintf("ac[%d].PL is nil", i))
			continue
		}
		if ac.BS == nil {
			slog.Info(fmt.Sprintf("ac[%d].BS is nil", i))
			continue
		}
		if ac.CF == nil {
			slog.Info(fmt.Sprintf("ac[%d].CF is nil", i))
			continue
		}
		p := message.NewPrinter(language.English)
		p.Printf("== %s ==\n", strconv.Itoa(ac.Date.Year()))
		p.Printf("PL.netProfit: %d\n", ac.PL.netProfit)
		p.Printf("PL.ordinaryProfit: %d\n", ac.PL.ordinaryProfit)
		p.Printf("PL.roe: %f\n", ac.PL.roe)
		p.Printf("PL.roa: %f\n", ac.PL.roa)
		p.Printf("BS.retainedEarnings: %d\n", ac.BS.retainedEarnings)
		p.Printf("BS.equityRatio: %f\n", ac.BS.equityRatio)
		p.Printf("CF.operatingCF: %d\n", ac.CF.operatingCF)
	}
}
