package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/peng225/jp-stock-analyzer/accounting"
	"github.com/peng225/jp-stock-analyzer/loader"
)

var logLevel = new(slog.LevelVar)

func init() {
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
				}
			case slog.SourceKey:
				if src, ok := a.Value.Any().(*slog.Source); ok {
					if rel, err := filepath.Rel(rootDir, src.File); err == nil {
						src.File = rel
					}
				}
			}
			return a
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func main() {
	var code string
	var verbose bool
	var strict bool
	flag.StringVar(&code, "c", "", "stock code")
	flag.BoolVar(&verbose, "v", false, "enable verbose mode")
	flag.BoolVar(&strict, "s", false, "enable strict mode")
	flag.Parse()

	if !verbose {
		logLevel.Set(slog.LevelWarn)
	}

	fileNames := []struct {
		pl string
		bs string
		cf string
		sd string
	}{
		{
			pl: "data/2025/fy-profit-and-loss.csv",
			bs: "data/2025/fy-balance-sheet.csv",
			cf: "data/2025/fy-cash-flow-statement.csv",
			sd: "data/2025/fy-stock-dividend.csv",
		},
		{
			pl: "data/2024/fy-profit-and-loss.csv",
			bs: "data/2024/fy-balance-sheet.csv",
			cf: "data/2024/fy-cash-flow-statement.csv",
			sd: "data/2024/fy-stock-dividend.csv",
		},
		{
			pl: "data/2023/fy-profit-and-loss.csv",
			bs: "data/2023/fy-balance-sheet.csv",
			cf: "data/2023/fy-cash-flow-statement.csv",
			sd: "data/2023/fy-stock-dividend.csv",
		},
	}

	acList := make(map[string][]*accounting.Accounting, 0)
	for _, fileName := range fileNames {
		loader.LoadPL(fileName.pl, acList)
		loader.LoadBS(fileName.bs, acList)
		loader.LoadCF(fileName.cf, acList)
		loader.LoadStockDividend(fileName.sd, acList)
	}

	if len(acList) == 0 {
		slog.Error("acList should not be empty.")
		os.Exit(1)
	}

	if code != "" {
		acs := acList[code]
		if risky, unknown := accounting.Risky(acs); risky || unknown {
			fmt.Printf("code %s is risky or does not have enough data.\n", code)
		} else {
			fmt.Printf("code %s is safe.\n", code)
		}
		accounting.Evaluate(acs)
		return
	}

	candidateCode := make([]string, 0, 128)
	for code, acs := range acList {
		slog.Info("Start check.", "code", code)
		if risky, unknown := accounting.Risky(acs); risky || unknown {
			slog.Info(fmt.Sprintf("code %s is risky or does not have enough data.", code))
			continue
		}
		if !accounting.IsGrowing(acs, strict) {
			slog.Info(fmt.Sprintf("code %s is not growing.", code))
			continue
		}
		candidateCode = append(candidateCode, code)
	}
	sort.Strings(candidateCode)
	fmt.Println(candidateCode)
}
