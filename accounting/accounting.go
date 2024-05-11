package accounting

import "time"

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
	jikoshihonRatio float64
}

type CashFlow struct {
	eigyoCF int64
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

func NewBalanceSheet(jikoshihonRatio float64) *BalanceSheet {
	return &BalanceSheet{
		jikoshihonRatio: jikoshihonRatio,
	}
}

func NewCashFlow(eigyoCF int64) *CashFlow {
	return &CashFlow{
		eigyoCF: eigyoCF,
	}
}
