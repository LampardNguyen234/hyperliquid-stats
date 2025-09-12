package api

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/LampardNguyen234/hyperliquid-stats/pkg/common"
)

type DailyVolumeByUser struct {
	Time   time.Time `json:"time"`
	User   string    `json:"user"`
	Volume float64   `json:"daily_usd_volume"`
}

func (item *DailyVolumeByUser) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Time   string  `json:"time"`
		User   string  `json:"user"`
		Volume float64 `json:"daily_usd_volume"`
	}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	item.User = tmp.User
	item.Volume = tmp.Volume
	item.Time, err = time.Parse("2006-01-02T15:04:05", tmp.Time)
	if err != nil {
		return err
	}

	return nil
}

type DailyVolumeByUsers []DailyVolumeByUser

func (data DailyVolumeByUsers) FilterByDateRange(fromDate, toDate *time.Time) DailyVolumeByUsers {
	if fromDate == nil && toDate == nil {
		return data
	}

	var filtered DailyVolumeByUsers
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

// FilterByUser filters the data to include only entries for a specific user
func (data DailyVolumeByUsers) FilterByUser(username string) DailyVolumeByUsers {
	if username == "" {
		return data
	}

	var filtered DailyVolumeByUsers
	for _, item := range data {
		// Case-insensitive user matching
		if strings.EqualFold(item.User, username) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// SortByTime sorts the data by time (descending by default) and then by volume (descending) within the same date
func (data DailyVolumeByUsers) SortByTime(descending bool) DailyVolumeByUsers {
	sorted := make(DailyVolumeByUsers, len(data))
	copy(sorted, data)

	sort.Slice(sorted, func(i, j int) bool {
		// First, compare by date (truncated to day level)
		dateI := time.Date(sorted[i].Time.Year(), sorted[i].Time.Month(), sorted[i].Time.Day(), 0, 0, 0, 0, sorted[i].Time.Location())
		dateJ := time.Date(sorted[j].Time.Year(), sorted[j].Time.Month(), sorted[j].Time.Day(), 0, 0, 0, 0, sorted[j].Time.Location())

		if !dateI.Equal(dateJ) {
			if descending {
				return dateI.After(dateJ) // More recent dates first
			}
			return dateI.Before(dateJ) // Older dates first
		}

		// If same date, sort by volume descending (highest volume first)
		return sorted[i].Volume > sorted[j].Volume
	})

	return sorted
}

func (data DailyVolumeByUsers) FormatString(count int) string {
	ret := common.NewTableFormatter().WithHeader("Daily Volume By User")
	ret = ret.WithHeader("Date", "User", "Volume ($M USD)")

	if count > len(data) {
		count = len(data)
	}
	if count == 0 {
		count = len(data)
	}

	for i := 0; i < count; i++ {
		ret = ret.WithRow(
			data[i].Time.Format("2006-01-02"),
			data[i].User,
			fmt.Sprintf("%.4f", data[i].Volume/1000000),
		)
	}

	return ret.String()
}

type DailyVolumeByUserResponse struct {
	Data DailyVolumeByUsers `json:"chart_data"`
}
