package models

   import (
       "time"
       "gorm.io/gorm"
   )

   type MerchantBankDetails struct {
       gorm.Model
       MerchantID     string `gorm:"type:uuid;uniqueIndex"`
       BankName       string
       BankCode       string `gorm:"size:10"`
       AccountNumber  string `gorm:"size:255"`
       AccountName    string   `gorm:"size:255"`
       RecipientCode  string `gorm:"size:50"`
       Currency       string `gorm:"size:8;default:'NGN'"`
     //  Status         string `gorm:"default:'pending'"`
       CreatedAt      time.Time
       UpdatedAt      time.Time
	   Merchant   Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
   }