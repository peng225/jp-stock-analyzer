package accounting

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeAndGrowing(t *testing.T) {
	acList := make(map[string][]*Accounting, 0)
	const code string = "100a"
	pl := NewProfitLoss(100, 80, 60, 40, 12, 7)
	require.NotNil(t, pl)
	bs := NewBalanceSheet(100, 300, 600, 25)
	require.NotNil(t, bs)
	cf := NewCashFlow(100, 1000)
	require.NotNil(t, cf)
	sd := NewStockDividend(30)
	require.NotNil(t, sd)
	acList[code] = []*Accounting{
		NewAccounting(time.Now(), pl, bs, cf, sd),
	}

	pl = NewProfitLoss(120, 100, 80, 60, 11, 6)
	require.NotNil(t, pl)
	bs = NewBalanceSheet(200, 200, 500, 35)
	require.NotNil(t, bs)
	cf = NewCashFlow(10, 800)
	require.NotNil(t, cf)
	sd = NewStockDividend(40)
	require.NotNil(t, sd)
	acList[code] = append(acList[code],
		NewAccounting(time.Now().AddDate(1, 0, 0), pl, bs, cf, sd))

	r, u := Risky(acList[code], true)
	assert.False(t, r)
	assert.False(t, u)
	assert.True(t, IsGrowing(acList[code], false))
}

func TestRiskyAndNotGrowing(t *testing.T) {
	acList := make(map[string][]*Accounting, 0)
	const code string = "100a"
	pl := NewProfitLoss(100, 80, 60, 40, 12, 7)
	require.NotNil(t, pl)
	bs := NewBalanceSheet(100, 500, 600, 25)
	require.NotNil(t, bs)
	cf := NewCashFlow(100, 1000)
	require.NotNil(t, cf)
	sd := NewStockDividend(30)
	require.NotNil(t, sd)
	acList[code] = []*Accounting{
		NewAccounting(time.Now(), pl, bs, cf, sd),
	}

	pl = NewProfitLoss(104, 100, 80, 60, 11, 6)
	require.NotNil(t, pl)
	bs = NewBalanceSheet(200, 200, 500, 35)
	require.NotNil(t, bs)
	cf = NewCashFlow(10, 800)
	require.NotNil(t, cf)
	sd = NewStockDividend(40)
	require.NotNil(t, sd)
	acList[code] = append(acList[code],
		NewAccounting(time.Now().AddDate(1, 0, 0), pl, bs, cf, sd))

	r, u := Risky(acList[code], true)
	assert.True(t, r)
	assert.False(t, u)
	assert.False(t, IsGrowing(acList[code], false))
}
