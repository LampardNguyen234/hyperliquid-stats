package api

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/LampardNguyen234/hyperliquid-stats/pkg/common"
)

type DailyVolume struct {
	Time   time.Time `json:"time"`
	Volume float64   `json:"daily_usd_volume"`
}

func (item *DailyVolume) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Time   string  `json:"time"`
		Volume float64 `json:"daily_usd_volume"`
	}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	item.Volume = tmp.Volume
	item.Time, err = time.Parse("2006-01-02T15:04:05", tmp.Time)
	if err != nil {
		return err
	}

	return nil
}

type DailyVolumes []DailyVolume

// FilterByDateRange filters the data to include only entries within the specified date range
func (data DailyVolumes) FilterByDateRange(fromDate, toDate *time.Time) DailyVolumes {
	if fromDate == nil && toDate == nil {
		return data
	}

	var filtered DailyVolumes
	for _, item := range data {
		// Truncate to date only for comparison
		itemDate := time.Date(item.Time.Year(), item.Time.Month(), item.Time.Day(), 0, 0, 0, 0, item.Time.Location())

		if fromDate != nil && itemDate.Before(*fromDate) {
			continue
		}

		if toDate != nil && itemDate.After(*toDate) {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}

// SortByTime sorts the data by time (descending by default)
func (data DailyVolumes) SortByTime(descending bool) DailyVolumes {
	sorted := make(DailyVolumes, len(data))
	copy(sorted, data)

	sort.Slice(sorted, func(i, j int) bool {
		if descending {
			return sorted[i].Time.After(sorted[j].Time) // More recent dates first
		}
		return sorted[i].Time.Before(sorted[j].Time) // Older dates first
	})

	return sorted
}

// FormatString formats the daily volume data as a table string
func (data DailyVolumes) FormatString(count int) string {
	ret := common.NewTableFormatter().WithHeader("Daily Volume")
	ret = ret.WithHeader("Date", "Volume ($B)")

	if count > len(data) {
		count = len(data)
	}
	if count == 0 {
		count = len(data)
	}

	sum := 0.0

	for i := 0; i < count; i++ {
		sum += data[i].Volume
		ret = ret.WithRow(
			data[i].Time.Format("2006-01-02"),
			fmt.Sprintf("%.4f", data[i].Volume/1000000000),
		)
	}
	ret = ret.WithFooter("SUM", fmt.Sprintf("%.4f", sum/1000000000))

	return ret.String()
}

type DailyVolumeResponse struct {
	Data DailyVolumes `json:"chart_data"`
}
