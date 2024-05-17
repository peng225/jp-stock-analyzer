package accounting

import (
	"log"
	"sort"
	"time"
)

type Accounting struct {
	Date time.Time
	PL   *ProfitLoss
	BS   *BalanceSheet
	CF   *CashFlow
}

type ProfitLoss struct {
	eigyoProfit int64
	keijoProfit int64
	netProfit   int64
	roe         float64
	roa         float64
}

type BalanceSheet struct {
	profitJouyoMoney      int64
	shortTermKariireMoeny int64
	longTermKariireMoeny  int64
	jikoshihonRatio       float64
}

type CashFlow struct {
	eigyoCF           int64
	genkinDoutouButsu int64
}

func NewAccounting(date time.Time, pl *ProfitLoss,
	bs *BalanceSheet, cf *CashFlow) *Accounting {
	return &Accounting{
		Date: date,
		PL:   pl,
		BS:   bs,
		CF:   cf,
	}
}

func NewProfitLoss(eigyoProfit, keijoProfit,
	netProfit int64, roe, roa float64) *ProfitLoss {
	return &ProfitLoss{
		eigyoProfit: eigyoProfit,
		keijoProfit: keijoProfit,
		netProfit:   netProfit,
		roe:         roe,
		roa:         roa,
	}
}

func NewBalanceSheet(profitJouyoMoney, stkMoey, ltkMoney int64,
	jikoshihonRatio float64) *BalanceSheet {
	return &BalanceSheet{
		profitJouyoMoney:      profitJouyoMoney,
		shortTermKariireMoeny: stkMoey,
		longTermKariireMoeny:  ltkMoney,
		jikoshihonRatio:       jikoshihonRatio,
	}
}

func NewCashFlow(eigyoCF, gdb int64) *CashFlow {
	return &CashFlow{
		eigyoCF:           eigyoCF,
		genkinDoutouButsu: gdb,
	}
}

func GoingBankrupt(acs []*Accounting) bool {
	if len(acs) == 0 {
		log.Println("no data")
	}
	bankrupt := true
	for _, ac := range acs {
		if ac.PL == nil {
			continue
		}
		if ac.PL.netProfit > 0 {
			bankrupt = false
			break
		}
	}
	if bankrupt {
		return true
	}

	bankrupt = true
	for _, ac := range acs {
		if ac.BS == nil {
			continue
		}
		if ac.BS.profitJouyoMoney > 0 && ac.BS.jikoshihonRatio >= 20 {
			bankrupt = false
			break
		}
		log.Println(ac.BS.profitJouyoMoney, ac.BS.jikoshihonRatio)
	}
	if bankrupt {
		return true
	}

	bankrupt = true
	for _, ac := range acs {
		if ac.CF == nil {
			continue
		}
		if ac.CF.eigyoCF > 0 {
			bankrupt = false
			break
		}
	}
	if bankrupt {
		return true
	}

	bankrupt = true
	for _, ac := range acs {
		if ac.BS == nil || ac.CF == nil {
			continue
		}
		if ac.CF.genkinDoutouButsu > ac.BS.shortTermKariireMoeny+ac.BS.longTermKariireMoeny {
			bankrupt = false
			break
		}
	}
	return bankrupt
}

func IsGrowing(acs []*Accounting) bool {
	sort.Slice(acs, func(i, j int) bool {
		return acs[i].Date.Before(acs[j].Date)
	})
	for i := 0; i < len(acs)-1; i++ {
		if acs[i].PL == nil || acs[i+1].PL == nil {
			return false
		}
		if acs[i].PL.keijoProfit >= acs[i+1].PL.keijoProfit {
			return false
		}
	}
	return true
}
