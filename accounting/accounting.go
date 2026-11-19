package accounting

import (
	"log"
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

func Risky(acs []*Accounting, verbose bool) (risky bool, unknown bool) {
	if !verbose {
		slog.SetLogLoggerLevel(slog.LevelWarn)
		defer slog.SetLogLoggerLevel(slog.LevelInfo)
	}
	if len(acs) == 0 {
		slog.Info("no data")
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
			if ac.BS.profitJouyoMoney <= 0 {
				slog.Info("Profit jouyo money is negative.")
				risky = true
			}
			if ac.BS.jikoshihonRatio < 20 {
				slog.Info("Jikoshihon ratio is too small.")
				risky = true
			}
		}
		if ac.CF == nil {
			slog.Info("CF is nil")
			unknown = true
		} else {
			if ac.CF.eigyoCF <= 0 {
				slog.Info("Eigyo cache flow is negative.")
				risky = true
			}
			if ac.BS != nil {
				if ac.CF.genkinDoutouButsu < ac.BS.shortTermKariireMoney+ac.BS.longTermKariireMoney {
					slog.Info("Kariire is larger than genkin doutou butsu.")
					risky = true
				}
			}
		}
		if ac.SD == nil {
			slog.Info("SD is nil")
			unknown = true
		} else {
			if ac.SD.haitouSeikou <= 20.0 || 50.0 <= ac.SD.haitouSeikou {
				slog.Info("Haitou seikou is not in the proper range.", "haitouSeikou", ac.SD.haitouSeikou)
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
		if float64(acs[i+1].PL.keijoProfit)/float64(acs[i].PL.keijoProfit) < 1.05 {
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
		log.Println("no enough data.")
		return
	}

	for i, ac := range acs {
		if ac.PL == nil {
			log.Printf("ac[%d].PL is nil", i)
			continue
		}
		if ac.BS == nil {
			log.Printf("ac[%d].BS is nil", i)
			continue
		}
		if ac.CF == nil {
			log.Printf("ac[%d].CF is nil", i)
			continue
		}
		p := message.NewPrinter(language.English)
		p.Printf("== %s ==\n", strconv.Itoa(ac.Date.Year()))
		p.Printf("PL.netProfit: %d\n", ac.PL.netProfit)
		p.Printf("PL.keijoProfit: %d\n", ac.PL.keijoProfit)
		p.Printf("PL.roe: %f\n", ac.PL.roe)
		p.Printf("PL.roa: %f\n", ac.PL.roa)
		p.Printf("BS.profitJouyoMoney: %d\n", ac.BS.profitJouyoMoney)
		p.Printf("BS.jikoshihonRatio: %f\n", ac.BS.jikoshihonRatio)
		p.Printf("CF.eigyoCF: %d\n", ac.CF.eigyoCF)
	}
}
