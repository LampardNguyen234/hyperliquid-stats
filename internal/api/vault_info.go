package api

import (
	"encoding/json"
	"fmt"
	"github.com/LampardNguyen234/hyperliquid-stats/pkg/common"
	"strconv"
)

type VaultSummary struct {
	Name    string  `json:"name"`
	Address string  `json:"vaultAddress"`
	Leader  string  `json:"leader"`
	TVL     float64 `json:"tvl,omitempty"`
	Closed  bool    `json:"isClosed"`
}

func (item *VaultSummary) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Name    string `json:"name"`
		Address string `json:"vaultAddress"`
		Leader  string `json:"leader"`
		TVL     string `json:"tvl"`
		Closed  bool   `json:"isClosed"`
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	item.Name = tmp.Name
	item.Address = tmp.Address
	item.Leader = tmp.Leader
	item.TVL, _ = strconv.ParseFloat(tmp.TVL, 64)
	item.Closed = tmp.Closed

	return nil
}

type Vault struct {
	Data VaultSummary `json:"summary"`
}

func (v *Vault) IsHLP() bool {
	return v.Data.Leader == "0xdfc24b077bc1425ad1dea75bcb6f8158e10df303"
}

type Vaults []Vault

func (data Vaults) FormatString(count int) string {
	ret := common.NewTableFormatter().WithHeader("Open Vaults (HLP First)")
	ret = ret.WithHeader("Name", "Address", "TVL", "Type")

	if count > len(data) {
		count = len(data)
	}

	for i := 0; i < count; i++ {
		vault := data[i]
		vaultType := "Norm"
		if vault.IsHLP() {
			vaultType = "HLP"
		}
		name := vault.Data.Name
		if len(name) > 20 {
			name = fmt.Sprintf("%v....%v", name[:8], name[len(name)-8:])
		}

		ret = ret.WithRow(
			name,
			vault.Data.Address,
			fmt.Sprintf("%.2f", vault.Data.TVL),
			vaultType,
		)
	}

	return ret.String()
}

func (data Vaults) FilterByStatus(closed bool) Vaults {
	var filtered Vaults
	for _, vault := range data {
		if vault.Data.Closed == closed {
			filtered = append(filtered, vault)
		}
	}
	return filtered
}

func (data Vaults) FilterOpenVaults() Vaults {
	return data.FilterByStatus(false)
}

func (data Vaults) FilterByMinTVL(minTVL float64) Vaults {
	var filtered Vaults
	for _, vault := range data {
		if vault.Data.TVL >= minTVL {
			filtered = append(filtered, vault)
		}
	}
	return filtered
}

func (data Vaults) SortWithHLPPriority(ascending bool) Vaults {
	result := make(Vaults, len(data))
	copy(result, data)

	for i := 0; i < len(result)-1; i++ {
		for j := 0; j < len(result)-i-1; j++ {
			shouldSwap := false

			// First priority: HLP vaults come first
			if result[j].IsHLP() != result[j+1].IsHLP() {
				shouldSwap = !result[j].IsHLP() && result[j+1].IsHLP()
			} else {
				// Same HLP status, sort by TVL
				if ascending {
					shouldSwap = result[j].Data.TVL > result[j+1].Data.TVL
				} else {
					shouldSwap = result[j].Data.TVL < result[j+1].Data.TVL
				}
			}

			if shouldSwap {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

func (data Vaults) SortByTVL(ascending bool) Vaults {
	result := make(Vaults, len(data))
	copy(result, data)

	for i := 0; i < len(result)-1; i++ {
		for j := 0; j < len(result)-i-1; j++ {
			shouldSwap := false
			if ascending {
				shouldSwap = result[j].Data.TVL > result[j+1].Data.TVL
			} else {
				shouldSwap = result[j].Data.TVL < result[j+1].Data.TVL
			}

			if shouldSwap {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}
