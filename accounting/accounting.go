package accounting

import (
	"log"
	"sort"
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
	revenue     int64
	eigyoProfit int64
	keijoProfit int64
	netProfit   int64
	roe         float64
	roa         float64
}

type BalanceSheet struct {
	profitJouyoMoney      int64
	shortTermKariireMoney int64
	longTermKariireMoney  int64
	jikoshihonRatio       float64
}

type CashFlow struct {
	eigyoCF           int64
	genkinDoutouButsu int64
}

type StockDividend struct {
	haitouSeikou float64
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

func NewProfitLoss(revenue, eigyoProfit, keijoProfit,
	netProfit int64, roe, roa float64) *ProfitLoss {
	return &ProfitLoss{
		revenue:     revenue,
		eigyoProfit: eigyoProfit,
		keijoProfit: keijoProfit,
		netProfit:   netProfit,
		roe:         roe,
		roa:         roa,
	}
}

func NewBalanceSheet(profitJouyoMoney, stkMoney, ltkMoney int64,
	jikoshihonRatio float64) *BalanceSheet {
	return &BalanceSheet{
		profitJouyoMoney:      profitJouyoMoney,
		shortTermKariireMoney: stkMoney,
		longTermKariireMoney:  ltkMoney,
		jikoshihonRatio:       jikoshihonRatio,
	}
}

func NewCashFlow(eigyoCF, gdb int64) *CashFlow {
	return &CashFlow{
		eigyoCF:           eigyoCF,
		genkinDoutouButsu: gdb,
	}
}

func NewStockDividend(hs float64) *StockDividend {
	return &StockDividend{
		haitouSeikou: hs,
	}
}

func Risky(acs []*Accounting) bool {
	if len(acs) == 0 {
		log.Println("no data")
	}

	sort.Slice(acs, func(i, j int) bool {
		return acs[i].Date.Before(acs[j].Date)
	})

	noDataConditions := []func(ac *Accounting) bool{
		func(ac *Accounting) bool {
			return ac.PL == nil
		},
		func(ac *Accounting) bool {
			return ac.BS == nil
		},
		func(ac *Accounting) bool {
			return ac.CF == nil
		},
		func(ac *Accounting) bool {
			return ac.SD == nil
		},
	}

	riskyConditions := []func(ac *Accounting) bool{
		// #0
		func(ac *Accounting) bool {
			return ac.PL == nil || ac.PL.netProfit <= 0
		},
		// #1
		func(ac *Accounting) bool {
			return ac.BS == nil || ac.BS.profitJouyoMoney <= 0 || ac.BS.jikoshihonRatio < 20
		},
		// #2
		func(ac *Accounting) bool {
			return ac.CF == nil || ac.CF.eigyoCF <= 0
		},
		// #3
		func(ac *Accounting) bool {
			return ac.BS == nil || ac.CF == nil ||
				ac.CF.genkinDoutouButsu <= ac.BS.shortTermKariireMoney+ac.BS.longTermKariireMoney
		},
		// #4
		func(ac *Accounting) bool {
			return ac.PL == nil || ac.PL.roe/ac.PL.roa >= 3.0
		},
		// #5
		func(ac *Accounting) bool {
			return ac.SD == nil || ac.SD.haitouSeikou <= 20.0 || 50.0 <= ac.SD.haitouSeikou
		},
	}

	for _, noDataCondition := range noDataConditions {
		for _, ac := range acs {
			if noDataCondition(ac) {
				log.Println("no data")
				return true
			}
		}
	}

	risky := false
	for i, riskyCondition := range riskyConditions {
		for j, ac := range acs {
			if riskyCondition(ac) {
				log.Printf("riskyConditions[%d] met for ac[%d].", i, j)
				risky = true
			}
		}
	}
	return risky
}

func IsGrowing(acs []*Accounting) bool {
	sort.Slice(acs, func(i, j int) bool {
		return acs[i].Date.Before(acs[j].Date)
	})
	if len(acs) < 2 {
		return false
	}
	for i := 0; i < len(acs)-1; i++ {
		// Check the growth of the keijo profit.
		if acs[i].PL == nil || acs[i+1].PL == nil {
			return false
		}
		if acs[i].PL.keijoProfit >= acs[i+1].PL.keijoProfit {
			return false
		}
		if acs[i].PL.keijoProfit < 0 || acs[i+1].PL.keijoProfit < 0 {
			return false
		}
		if float64(acs[i+1].PL.keijoProfit)/float64(acs[i].PL.keijoProfit) < 1.08 {
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
		if float64(acs[i+1].PL.revenue)/float64(acs[i].PL.revenue) < 1.08 {
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
		log.Println("no enough data.")
		return
	}

	for i, ac := range acs {
		if ac.PL == nil {
			log.Printf("ac[%d].PL is nil", i)
			continue
		}
		p := message.NewPrinter(language.English)
		p.Printf("ac[%d].PL.keijoProfit: %d\n", i, ac.PL.keijoProfit)
		p.Printf("ac[%d].PL.roe: %f\n", i, ac.PL.roe)
		p.Printf("ac[%d].PL.roa: %f\n", i, ac.PL.roa)
	}
}
