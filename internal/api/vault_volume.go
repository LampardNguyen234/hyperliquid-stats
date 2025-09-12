package api

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/LampardNguyen234/hyperliquid-stats/pkg/common"
	"github.com/pkg/errors"
)

type VaultVolume struct {
	Day     float64
	Week    float64
	Month   float64
	AllTime float64

	PerpDay     float64
	PerpWeek    float64
	PerpMonth   float64
	PerpAllTime float64
}

func (v *VaultVolume) UnmarshalJSON(data []byte) error {
	var tmp [][]interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "failed to unmarshal to []interface{}")
	}

	for _, d := range tmp {
		if len(d) < 2 {
			return errors.New("invalid vault volume")
		}

		var tmpMap map[string]interface{}
		jsb, _ := json.Marshal(d[1])
		if err := json.Unmarshal(jsb, &tmpMap); err != nil {
			return errors.Wrap(err, "failed to unmarshal to map")
		}

		switch d[0].(string) {
		case "day":
			v.Day, _ = strconv.ParseFloat(tmpMap["vlm"].(string), 64)
		case "week":
			v.Week, _ = strconv.ParseFloat(tmpMap["vlm"].(string), 64)
		case "month":
			v.Month, _ = strconv.ParseFloat(tmpMap["vlm"].(string), 64)
		case "allTime":
			v.AllTime, _ = strconv.ParseFloat(tmpMap["vlm"].(string), 64)
		case "perpDay":
			v.PerpDay, _ = strconv.ParseFloat(tmpMap["vlm"].(string), 64)
		case "perpWeek":
			v.PerpWeek, _ = strconv.ParseFloat(tmpMap["vlm"].(string), 64)
		case "perpMonth":
			v.PerpMonth, _ = strconv.ParseFloat(tmpMap["vlm"].(string), 64)
		case "perpAllTime":
			v.PerpAllTime, _ = strconv.ParseFloat(tmpMap["vlm"].(string), 64)
		default:
			return errors.New("invalid vault volume")
		}
	}

	return nil
}

type VaultVolumeRequest struct {
	Type    string `json:"type"`
	Address string `json:"vaultAddress"`
	User    string `json:"user,omitempty"`
}

type VaultVolumeResponse struct {
	Portfolio VaultVolume `json:"portfolio"`
}

type VaultVolumeInfo struct {
	Address string
	Name    string
	Volume  VaultVolume
	TVL     float64
	IsHLP   bool
}

type VaultVolumesInfo []VaultVolumeInfo

func (data VaultVolumesInfo) SortByField(field string) VaultVolumesInfo {
	result := make(VaultVolumesInfo, len(data))
	copy(result, data)

	switch strings.ToLower(field) {
	case "day":
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Volume.Day > result[j].Volume.Day
		})
	case "week":
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Volume.Week > result[j].Volume.Week
		})
	case "month":
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Volume.Month > result[j].Volume.Month
		})
	case "all-time":
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Volume.AllTime > result[j].Volume.AllTime
		})

	}

	return result
}

func (data VaultVolumesInfo) FormatSummary() string {
	// Aggregate volumes by HLP status
	hlpTotals := struct {
		Day     float64
		Week    float64
		Month   float64
		AllTime float64
		TVL     float64
		Count   int
	}{}

	nonHLPTotals := struct {
		Day     float64
		Week    float64
		Month   float64
		AllTime float64
		TVL     float64
		Count   int
	}{}

	// Calculate totals
	for _, vault := range data {
		if vault.IsHLP {
			hlpTotals.Day += vault.Volume.PerpDay
			hlpTotals.Week += vault.Volume.PerpWeek
			hlpTotals.Month += vault.Volume.PerpMonth
			hlpTotals.AllTime += vault.Volume.PerpAllTime
			hlpTotals.TVL += vault.TVL
			hlpTotals.Count++
		} else {
			nonHLPTotals.Day += vault.Volume.PerpDay
			nonHLPTotals.Week += vault.Volume.PerpWeek
			nonHLPTotals.Month += vault.Volume.PerpMonth
			nonHLPTotals.AllTime += vault.Volume.PerpAllTime
			nonHLPTotals.TVL += vault.TVL
			nonHLPTotals.Count++
		}
	}

	var result strings.Builder

	// HLP Totals Section
	result.WriteString("=== VAULT VOLUME SUMMARY ===\n\n")

	hlpTable := common.NewTableFormatter().WithHeader("HLP Vaults Summary")
	hlpTable = hlpTable.WithHeader("Metric", "Day", "Week", "Month", "All Time", "TVL")
	hlpTable = hlpTable.WithRow(
		fmt.Sprintf("Total (%d vaults)", hlpTotals.Count),
		fmt.Sprintf("%.3f", hlpTotals.Day/1000000),
		fmt.Sprintf("%.3f", hlpTotals.Week/1000000),
		fmt.Sprintf("%.3f", hlpTotals.Month/1000000),
		fmt.Sprintf("%.3f", hlpTotals.AllTime/1000000),
		fmt.Sprintf("%.3f", hlpTotals.TVL/1000000),
	)
	hlpTable = hlpTable.WithCaption("HLP Volume Summary (Values are in $M)")
	result.WriteString(hlpTable.String())
	result.WriteString("\n")

	// Non-HLP Totals Section
	nonHLPTable := common.NewTableFormatter().WithHeader("Non-HLP Vaults Summary")
	nonHLPTable = nonHLPTable.WithHeader("Metric", "Day", "Week", "Month", "All Time", "TVL")
	nonHLPTable = nonHLPTable.WithRow(
		fmt.Sprintf("Total (%d vaults)", nonHLPTotals.Count),
		fmt.Sprintf("%.3f", nonHLPTotals.Day/1000000),
		fmt.Sprintf("%.3f", nonHLPTotals.Week/1000000),
		fmt.Sprintf("%.3f", nonHLPTotals.Month/1000000),
		fmt.Sprintf("%.3f", nonHLPTotals.AllTime/1000000),
		fmt.Sprintf("%.3f", nonHLPTotals.TVL/1000000),
	)
	nonHLPTable = nonHLPTable.WithCaption("Non-HLP Volume Summary (Values are in $M)")
	result.WriteString(nonHLPTable.String())
	result.WriteString("\n")

	// Top 10 TVL Section
	topTVL := make(VaultVolumesInfo, len(data))
	copy(topTVL, data)

	// Sort by TVL descending (HLP first)
	sort.SliceStable(topTVL, func(i, j int) bool {
		if topTVL[i].IsHLP != topTVL[j].IsHLP {
			return topTVL[i].IsHLP && !topTVL[j].IsHLP
		}
		return topTVL[i].TVL > topTVL[j].TVL
	})

	// Limit to top 10
	displayCount := 10
	if len(topTVL) < 10 {
		displayCount = len(topTVL)
	}

	topTable := common.NewTableFormatter().WithHeader("Top 10 Vaults by TVL")
	topTable = topTable.WithHeader("Rank", "Address", "Type", "TVL", "Day", "Week", "Month", "All Time")

	for i := 0; i < displayCount; i++ {
		vault := topTVL[i]
		vaultType := "Vault"
		if vault.IsHLP {
			vaultType = "HLP"
		}
		topTable = topTable.WithRow(
			fmt.Sprintf("#%d", i+1),
			vault.Address,
			vaultType,
			fmt.Sprintf("%.3f", vault.TVL/1000000),
			fmt.Sprintf("%.3f", vault.Volume.Day/1000000),
			fmt.Sprintf("%.3f", vault.Volume.Week/1000000),
			fmt.Sprintf("%.3f", vault.Volume.Month/1000000),
			fmt.Sprintf("%.3f", vault.Volume.AllTime/1000000),
		)
	}
	topTable = topTable.WithCaption("Values are in $M")
	result.WriteString(topTable.String())

	return result.String()
}

func (data VaultVolumesInfo) FormatString() string {
	ret := common.NewTableFormatter().WithHeader("Vault Volumes")
	ret = ret.WithHeader("Address", "Type", "TVL", "Day", "Week", "Month", "All Time")

	for i := 0; i < len(data); i++ {
		vault := data[i]
		t := "Norm"
		if vault.IsHLP {
			t = "HLP"
		}
		ret = ret.WithRow(
			vault.Address,
			t,
			fmt.Sprintf("%.3f", vault.TVL/1000000),
			fmt.Sprintf("%.3f", vault.Volume.PerpDay/1000000),
			fmt.Sprintf("%.3f", vault.Volume.PerpWeek/1000000),
			fmt.Sprintf("%.3f", vault.Volume.PerpMonth/1000000),
			fmt.Sprintf("%.3f", vault.Volume.PerpAllTime/1000000),
		)
	}
	ret = ret.WithCaption("Values are in $M")

	return ret.String()
}

func (v VaultVolume) FormatSingle(name, address string) string {
	ret := common.NewTableFormatter().WithHeader(fmt.Sprintf("Vault Volume: %s", name))
	ret = ret.WithHeader("Period", "Volume", "Perp Volume")

	ret = ret.WithRow("Day", fmt.Sprintf("%.2f", v.Day), fmt.Sprintf("%.2f", v.PerpDay))
	ret = ret.WithRow("Week", fmt.Sprintf("%.2f", v.Week), fmt.Sprintf("%.2f", v.PerpWeek))
	ret = ret.WithRow("Month", fmt.Sprintf("%.2f", v.Month), fmt.Sprintf("%.2f", v.PerpMonth))
	ret = ret.WithRow("All Time", fmt.Sprintf("%.2f", v.AllTime), fmt.Sprintf("%.2f", v.PerpAllTime))

	return ret.String()
}
