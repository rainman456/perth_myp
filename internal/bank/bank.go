package bank

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
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
		file, err := os.ReadFile("bank.json")
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