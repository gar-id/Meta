package types

import (
	"time"

	"gorm.io/gorm"
)

type UserData struct {
	UserID     string `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Token      string
	Username   string
	MFA        bool
	TOTPSecret string
	Role       string
	Disabled   bool
}
