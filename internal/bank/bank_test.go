package bank
import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBankCode_Valid(t *testing.T) {
	// Setup test bank.json
	content := `{"9 PAYMENT SERVICE BANK": "120001", "ABBEY MORTGAGE BANK": "070010"}`
	err := os.WriteFile("bank.json", []byte(content), 0644)
	assert.NoError(t, err)
	defer os.Remove("bank.json")

	svc := GetBankService()
	err = svc.LoadBanks()
	assert.NoError(t, err)

	code, err := svc.GetBankCode("9 PAYMENT SERVICE BANK")
	assert.NoError(t, err)
	assert.Equal(t, "120001", code)
}

func TestGetBankCode_Invalid(t *testing.T) {
	svc := GetBankService()
	_, err := svc.GetBankCode("NONEXISTENT BANK")
	assert.ErrorIs(t, err, ErrBankNotFound)
}

func TestGetBankCode_CaseInsensitive(t *testing.T) {
	content := `{"TEST BANK": "123456"}`
	err := os.WriteFile("../../bank.json", []byte(content), 0644)
	assert.NoError(t, err)
	defer os.Remove("../../bank.json")

	svc := GetBankService()
	svc.LoadBanks()

	code, err := svc.GetBankCode("test bank")
	assert.NoError(t, err)
	assert.Equal(t, "123456", code)
}