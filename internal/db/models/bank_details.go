package models

   import (
       "time"
       "gorm.io/gorm"
   )

   type MerchantBankDetails struct {
       gorm.Model
       MerchantID     string `gorm:"type:uuid;uniqueIndex"`
       BankName       string
       BankCode       string `gorm:"size:5"`
       AccountNumber  string `gorm:"size:15"`
       AccountName    string
       RecipientCode  string `gorm:"size:50"`
       Currency       string `gorm:"size:3;default:'NGN'"`
       Status         string `gorm:"default:'pending'"`
       CreatedAt      time.Time
       UpdatedAt      time.Time
	   Merchant   Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
   }