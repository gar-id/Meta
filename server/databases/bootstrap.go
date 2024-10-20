package databases

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"MetaHandler/server/config/caches"
	"MetaHandler/server/databases/types"

	"gorm.io/gorm"
)

func switchDB() *gorm.DB {
	// Create DB variable
	var DB *gorm.DB

	// Check database connection type
	switch caches.MetaHandlerServer.MetaHandlerServer.Database.Driver {
	case "sqlite":
		DB = sqliteConnection()
	case "postgresql":
		DB = postgresConnection()
	case "mysql":
		DB = mysqlConnection()
	}

	// Return DB
	return DB
}

func migration(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(
		&types.UserData{},
		&types.GroupData{},
		&types.HostData{},
		&types.ClientAccess{},
	)

	// Add root user if not available
	var rootUser types.UserData
	userBytes := []byte("admin")
	userId := fmt.Sprintf("%x", md5.Sum(userBytes))
	rowEffected, _ := GetUser(&rootUser, userId)
	if rowEffected == 0 {
		// Generate root account
		if rootUser.UserID == "" {
			InsertUser("admin", "")
		}
		var rootData = types.UserData{
			UserID:           userId,
			Username:         "admin",
			Role:             "admin",
			PublicKeyExpired: time.Now().Add(24 * time.Hour * 30 * 12 * 100),
		}
		UpdateUser(rootData)
	}

	// Add agent user if not available
	var agentUser types.UserData
	agentBytes := []byte("centralissh-agent")
	agentId := fmt.Sprintf("%x", md5.Sum(agentBytes))
	rowEffected, _ = GetUser(&agentUser, agentId)
	if rowEffected == 0 {
		// Generate root account
		if rootUser.UserID == "" {
			InsertUser("centralissh-agent", "")
		}
		var agentData = types.UserData{
			UserID:           userId,
			Username:         "centralissh-agent",
			Role:             "agent",
			Token:            caches.MetaHandlerServer.MetaHandlerServer.API.AgentToken,
			PublicKeyExpired: time.Now().Add(24 * time.Hour * 30 * 12 * 100),
		}
		UpdateUser(agentData)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()
}

func createSQLiteDB() {
	// Parse SQLite location
	dbPath := strings.Split(caches.MetaHandlerServer.MetaHandlerServer.Database.Sqlite.Location, "/")
	var newPath string
	for _, splitDBPath := range dbPath {
		if !strings.Contains(splitDBPath, ".db") {
			newPath = filepath.Join(newPath, splitDBPath)
		}
	}

	// Create directory if doesn't exist
	os.MkdirAll(newPath, os.ModePerm)

	// Create SQLite DB file
	filelocation := caches.MetaHandlerServer.MetaHandlerServer.Database.Sqlite.Location
	_, err := os.OpenFile(filelocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Fatalf("error opening SQLite file: %v", err)
	}
}

func Bootstrap() {
	// Reformat with lower string
	caches.MetaHandlerServer.MetaHandlerServer.Database.Driver = strings.ToLower(caches.MetaHandlerServer.MetaHandlerServer.Database.Driver)

	// Create SQLite file if driver is SQLite
	switch caches.MetaHandlerServer.MetaHandlerServer.Database.Driver {
	case "sqlite":
		createSQLiteDB()
	}

	// Create connection from bootstrap func
	DB := switchDB()

	// Run database migration
	migration(DB)

}
