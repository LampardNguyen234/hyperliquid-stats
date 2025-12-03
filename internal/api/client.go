package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/LampardNguyen234/go-rate-limiter"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	baseURL    string
	infoURL    string
	httpClient *http.Client
	limiter    rate.RateLimiter
}

func NewClient(baseURL, infoURL string) *Client {
	limiter, _ := rate.NewMultipleLimiter(
		rate.NewLimiter(time.Minute, 300),
		rate.NewLimiter(100*time.Millisecond, 3),
	)

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		infoURL: strings.TrimRight(infoURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		limiter: limiter,
	}
}

func (c *Client) BuildURL(path string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, strings.TrimLeft(path, "/"))
}

func (c *Client) FetchLargestUsers() (USDVolumeByUsers, error) {
	var result LargestVolumeResponse
	err := c.fetchData(c.BuildURL("largest_users_by_usd_volume"), &result)
	return result.Data, err
}

func (c *Client) FetchLargestTradeCounts() (LargestTradeCounts, error) {
	var result LargestTradeCountResponse
	err := c.fetchData(c.BuildURL("largest_users_by_trade_count"), &result)
	return result.Data, err
}

func (c *Client) FetchDailyVolumeByUser(fromDate, toDate *time.Time, username string) (DailyVolumeByUsers, error) {
	var result DailyVolumeByUserResponse

	// Build endpoint with optional parameters
	endpoint := "daily_usd_volume_by_user"
	params := url.Values{}

	if fromDate != nil {
		params.Add("from_date", fromDate.Format("2006-01-02"))
	}

	if toDate != nil {
		params.Add("to_date", toDate.Format("2006-01-02"))
	}

	if username != "" {
		params.Add("user", username)
	}

	if len(params) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	err := c.fetchData(c.BuildURL(endpoint), &result)
	if err != nil {
		return nil, err
	}

	// Apply client-side filtering if needed
	if fromDate != nil || toDate != nil {
		result.Data = result.Data.FilterByDateRange(fromDate, toDate)
	}

	// Apply client-side user filtering if needed (fallback if API doesn't support it)
	if username != "" {
		result.Data = result.Data.FilterByUser(username)
	}

	return result.Data, nil
}

func (c *Client) FetchDailyVolume(fromDate, toDate *time.Time) (DailyVolumes, error) {
	var result DailyVolumeResponse

	// Build endpoint with optional parameters
	endpoint := "daily_usd_volume"
	params := url.Values{}

	if fromDate != nil {
		params.Add("from_date", fromDate.Format("2006-01-02"))
	}

	if toDate != nil {
		params.Add("to_date", toDate.Format("2006-01-02"))
	}

	if len(params) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	err := c.fetchData(c.BuildURL(endpoint), &result)
	if err != nil {
		return nil, err
	}

	// Apply client-side filtering if needed
	if fromDate != nil || toDate != nil {
		result.Data = result.Data.FilterByDateRange(fromDate, toDate)
	}

	return result.Data, nil
}

func (c *Client) FetchAllVault() (Vaults, error) {
	result := Vaults{}
	err := c.fetchData("https://stats-data.hyperliquid.xyz/Mainnet/vaults", &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) FetchVaultVolume(vaultAddress string) (VaultVolume, error) {
	if !c.limiter.Allow() {
		return VaultVolume{}, errors.New("429: rate limit exceeded")
	}

	var result VaultVolumeResponse

	payload := VaultVolumeRequest{
		Type:    "vaultDetails",
		Address: vaultAddress,
		User:    "",
	}

	err := c.PostRequest(c.infoURL, payload, &result)
	if err != nil {
		return VaultVolume{}, errors.Wrapf(err, "failed to fetch vault volume for address %s", vaultAddress)
	}

	return result.Portfolio, nil
}

func (c *Client) FetchAllVaultVolumes(hlpOnly bool, count int) (VaultVolumesInfo, error) {
	return c.FetchAllVaultVolumesConcurrent(hlpOnly, count, 1)
}

func (c *Client) FetchAllVaultVolumesConcurrent(hlpOnly bool, count int, workers int) (VaultVolumesInfo, error) {
	start := time.Now()
	// First get all vaults
	vaults, err := c.FetchAllVault()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch vaults")
	}

	vaults = vaults.FilterOpenVaults().FilterByMinTVL(10).SortWithHLPPriority(false)
	if count > 0 && len(vaults) > count {
		vaults = vaults[:count+1]
	}

	// Filter vaults to process
	var vaultsToProcess []Vault
	for _, vault := range vaults {
		if vault.Data.Address == "0xdfc24b077bc1425ad1dea75bcb6f8158e10df303" || vault.Data.Closed {
			continue
		}
		if hlpOnly && !vault.IsHLP() {
			continue
		}
		vaultsToProcess = append(vaultsToProcess, vault)
	}

	if len(vaultsToProcess) == 0 {
		return VaultVolumesInfo{}, nil
	}

	// Set up worker pool
	if workers <= 0 {
		workers = 1
	}
	if workers > len(vaultsToProcess) {
		workers = len(vaultsToProcess)
	}

	// Channels for work distribution and result collection
	vaultChan := make(chan Vault, len(vaultsToProcess))
	resultChan := make(chan VaultVolumeInfo, len(vaultsToProcess))
	errorChan := make(chan error)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for vault := range vaultChan {
				tmpCount := 0
				for {
					volume, err := c.FetchVaultVolume(vault.Data.Address)
					if err != nil {
						if strings.Contains(err.Error(), "429") && tmpCount < 20 {
							tmpCount++

							t := int64(math.Min(float64(workers/2+1), 5)) * int64(time.Second)
							rand.NewSource(time.Now().Unix())
							t += rand.Int63() % t / 2

							time.Sleep(time.Duration(t))
							continue
						}
						errorChan <- errors.Wrapf(err, "failed to fetch volume for vault %s", vault.Data.Address)
						break
					}

					resultChan <- VaultVolumeInfo{
						Address: vault.Data.Address,
						Name:    vault.Data.Name,
						Volume:  volume,
						TVL:     vault.Data.TVL,
						IsHLP:   vault.IsHLP(),
					}
					break
				}
			}
		}()
	}

	// Send work to workers
	go func() {
		defer close(vaultChan)
		for _, vault := range vaultsToProcess {
			vaultChan <- vault
		}
	}()

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results
	var result VaultVolumesInfo
	var fetchErrors []error

	for {
		select {
		case vaultInfo := <-resultChan:
			result = append(result, vaultInfo)
			fmt.Printf("DONE for vault: %v, timeElapsed: %v, count: %v/%v\n", vaultInfo.Address, time.Since(start).String(), len(result), len(vaultsToProcess))
		case err = <-errorChan:
			fmt.Println("new error:", err)
			fetchErrors = append(fetchErrors, err)
		}

		if len(result)+len(fetchErrors) == len(vaultsToProcess) {
			break
		}
	}

	// Return results even if some requests failed
	return result, nil
}

func (c *Client) FetchData(url string, result interface{}) error {
	return c.fetchData(url, result)
}

func (c *Client) PostRequest(url string, payload interface{}, result interface{}) error {
	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal request payload")
	}

	// Create the POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrapf(err, "failed to create request")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed to make POST request to %s", url)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read response")
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return errors.Wrapf(err, "failed to parse JSON response")
	}

	return nil
}

func (c *Client) fetchData(url string, result interface{}) error {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch data")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("api returned status %d for %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read response")
	}

	return json.Unmarshal(body, &result)
}
