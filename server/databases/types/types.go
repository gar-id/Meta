package types

import (
	"gorm.io/gorm"
)

type MetaServer struct {
	gorm.Model
	ServerID         string
	Hostname         string
	PublicIP         string
	Environment      string
	Status           string
	ActiveConnection string
}

type Service struct {
	gorm.Model
	ServerID    string
	ServiceName string
	Status      string
}

type Stunnel struct {
	gorm.Model
	ServerID       string
	ConnectionName string
	Config         string
}
