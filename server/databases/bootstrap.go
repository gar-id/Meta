package databases

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"MetaHandler/server/config/caches"
	"MetaHandler/server/databases/types"
	"MetaHandler/tools"

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
		&types.MetaServer{},
		&types.Service{},
		&types.Stunnel{},
	)

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
		tools.ZapLogger("both", "server").Error(fmt.Sprintf("error opening SQLite file: %v", err))
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
