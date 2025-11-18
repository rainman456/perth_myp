package bank

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"api-customer-merchant/internal/api/dto"
)

var (
	ErrBankNotFound   = errors.New("invalid bank name")
	ErrBankFileNotFound = errors.New("bank.json file not found")
)

type BankService struct {
	banks map[string]string // name -> code
	mu    sync.RWMutex
	once  sync.Once
}





// Paystack API response minimal types
type paystackBank struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Slug string `json:"slug,omitempty"`
	// add more fields if you need them
}

type paystackResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    []paystackBank `json:"data"`
	Meta    any            `json:"meta,omitempty"`
}

type FetchBankService struct {
	apiURL     string
	secret     string
	client     *http.Client
	cacheFile  string
	mu         sync.RWMutex
	banks      map[string]string // normalized name -> code
	lastLoaded time.Time
}

func NewFetchBankService(opts ...func(*FetchBankService)) *FetchBankService {
	s := &FetchBankService{
		apiURL:    "https://api.paystack.co/bank",
		secret:    os.Getenv("PAYSTACK_SECRET_KEY"),
		client:    &http.Client{Timeout: 8 * time.Second},
		cacheFile: "banks.json",
		banks:     make(map[string]string),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// FetchAndCacheBanks calls Paystack, updates in-memory cache and writes raw response to cache file.
func (s *FetchBankService) FetchAndCacheBanks(ctx context.Context, country string) ([]dto.BankListItemDTO, error) {
	// Build request URL with optional query
	u, err := url.Parse(s.apiURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	if country != "" {
		q.Set("country", country)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if s.secret == "" {
		return nil, errors.New("missing PAYSTACK_SECRET (set env PAYSTACK_SECRET)")
	}
	req.Header.Set("Authorization", "Bearer "+s.secret)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("paystack request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// try to include response body in error (but avoid huge prints)
		short := string(body)
		if len(short) > 400 {
			short = short[:400] + "..."
		}
		return nil, fmt.Errorf("paystack returned status %d: %s", resp.StatusCode, short)
	}

	// Save raw response to cache file (atomic write)
	if err := writeFileAtomic(s.cacheFile, body, 0o644); err != nil {
		// Non-fatal: still proceed with parsed result but log error
		// If you have a logger, replace fmt.Printf
		fmt.Printf("warning: failed to write cache file: %v\n", err)
	}

	// Parse response
	var psResp paystackResponse
	if err := json.Unmarshal(body, &psResp); err != nil {
		return nil, fmt.Errorf("invalid paystack response: %w", err)
	}

	// Build dtos and update in-memory map
	items := make([]dto.BankListItemDTO, 0, len(psResp.Data))
	tempMap := make(map[string]string, len(psResp.Data))
	for _, b := range psResp.Data {
		name := strings.ToUpper(strings.TrimSpace(b.Name))
		code := b.Code
		tempMap[name] = code
		items = append(items, dto.BankListItemDTO{
			Name: b.Name,
			Code: code,
		})
	}

	// update cache
	s.mu.Lock()
	s.banks = tempMap
	s.lastLoaded = time.Now()
	s.mu.Unlock()

	// sort
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

// LoadBanksFromFile reads the cache file and returns the bank items.
// file is expected to be the raw Paystack response structure.
func (s *FetchBankService) LoadBanksFromFile() ([]dto.BankListItemDTO, error) {
	path := s.cacheFile
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrBankFileNotFound
		}
		return nil, err
	}

	var psResp paystackResponse
	if err := json.Unmarshal(data, &psResp); err != nil {
		return nil, err
	}

	items := make([]dto.BankListItemDTO, 0, len(psResp.Data))
	tempMap := make(map[string]string, len(psResp.Data))
	for _, b := range psResp.Data {
		name := strings.ToUpper(strings.TrimSpace(b.Name))
		code := b.Code
		tempMap[name] = code
		items = append(items, dto.BankListItemDTO{
			Name: b.Name,
			Code: code,
		})
	}

	// update in-memory cache
	s.mu.Lock()
	s.banks = tempMap
	s.lastLoaded = time.Now()
	s.mu.Unlock()

	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

// GetBanks tries API first, falls back to local file if API fails.
func (s *FetchBankService) GetBanks(ctx context.Context, country string) ([]dto.BankListItemDTO, error) {
	items, err := s.FetchAndCacheBanks(ctx, country)
	if err == nil {
		return items, nil
	}

	// API failed — fallback to file
	fmt.Printf("info: fetch failed, falling back to cache file: %v\n", err)
	items, fErr := s.LoadBanksFromFile()
	if fErr != nil {
		return nil, fmt.Errorf("both fetch and cache failed: fetch error: %w; cache error: %v", err, fErr)
	}
	return items, nil
}

// GetBankCode returns code for bank name from current cache (must have been loaded)
func (s *FetchBankService) GetBankCode(bankName string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(bankName))
	s.mu.RLock()
	defer s.mu.RUnlock()
	code, ok := s.banks[normalized]
	if !ok {
		return "", ErrBankNotFound
	}
	return code, nil
}

// helper: atomic write file
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}




var instance *BankService

// GetBankService returns singleton instance
func GetBankService() *BankService {
	if instance == nil {
		instance = &BankService{
			banks: make(map[string]string),
		}
	}
	return instance
}

// LoadBanks loads bank.json once
func (bs *BankService) LoadBanks() error {
	var loadErr error
	bs.once.Do(func() {
		file, err := os.ReadFile("banks.json")
		if err != nil {
			loadErr = ErrBankFileNotFound
			return
		}

		var data map[string]string
		if err := json.Unmarshal(file, &data); err != nil {
			loadErr = err
			return
		}

		bs.mu.Lock()
		bs.banks = data
		bs.mu.Unlock()
	})
	return loadErr
}

// GetBankCode validates and returns bank code
func (bs *BankService) GetBankCode(bankName string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(bankName))
	
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	code, ok := bs.banks[normalized]
	if !ok {
		return "", ErrBankNotFound
	}
	return code, nil
}

// GetAllBanks returns all banks
func (bs *BankService) GetAllBanks() map[string]string {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Return copy
	result := make(map[string]string, len(bs.banks))
	for k, v := range bs.banks {
		result[k] = v
	}
	return result
}