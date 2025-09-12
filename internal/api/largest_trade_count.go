package api

import (
	"encoding/json"
	"fmt"

	"github.com/LampardNguyen234/hyperliquid-stats/pkg/common"
)

type LargestTradeCount struct {
	Name  string `json:"name"`
	Value uint64 `json:"value"`
}

func (item *LargestTradeCount) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	item.Name = tmp.Name
	item.Value = uint64(tmp.Value)
	return nil
}

type LargestTradeCounts []LargestTradeCount

func (data LargestTradeCounts) FormatString(count int) string {
	ret := common.NewTableFormatter().WithHeader("Largest Trade Count By User")
	ret = ret.WithHeader("User", "Trade Count")

	if count > len(data) {
		count = len(data)
	}

	for i := 0; i < count; i++ {
		ret = ret.WithRow(data[i].Name, fmt.Sprintf("%d", data[i].Value))
	}

	return ret.String()
}

type LargestTradeCountResponse struct {
	Data LargestTradeCounts `json:"chart_data"`
}
