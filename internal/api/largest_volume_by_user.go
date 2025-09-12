package api

import (
	"fmt"
	"github.com/LampardNguyen234/hyperliquid-stats/pkg/common"
)

type USDVolumeByUser struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type USDVolumeByUsers []USDVolumeByUser

func (data USDVolumeByUsers) FormatString(count int) string {
	ret := common.NewTableFormatter().WithHeader("Largest Volume By User")
	ret = ret.WithHeader("User", "Value ($M USD)")

	if count > len(data) {
		count = len(data)
	}

	for i := 0; i < count; i++ {
		ret = ret.WithRow(data[i].Name, fmt.Sprintf("%.4f", data[i].Value/1000000))
	}

	return ret.String()
}

type LargestVolumeResponse struct {
	Data USDVolumeByUsers `json:"table_data"`
}
