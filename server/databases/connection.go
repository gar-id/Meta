package databases

import (
	"fmt"
	"log"
	"os"
	"time"

	"MetaHandler/server/config/caches"

	"MetaHandler/tools"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func gormLogger() logger.Interface {
	var gormLogLevel logger.LogLevel
	switch caches.MetaHandlerServer.MetaHandlerServer.Log.Level {
	case "debug":
		gormLogLevel = logger.Info
	case "info":
		gormLogLevel = logger.Info
	case "warning":
		gormLogLevel = logger.Warn
	case "error":
		gormLogLevel = logger.Error
	case "panic":
		gormLogLevel = logger.Error
	case "fatal":
		gormLogLevel = logger.Error
	default:
		gormLogLevel = logger.Silent
	}
	filelocation := tools.DefaultString(caches.MetaHandlerServer.MetaHandlerServer.Log.Location, "/var/log")
	filelocation = fmt.Sprintf("%v/MetaHandler-Gorm.log", filelocation)
	logfile, _ := os.OpenFile(filelocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	newLogger := logger.New(
		log.New(logfile, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  gormLogLevel, // Log level
			IgnoreRecordNotFoundError: false,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,        // Don't include params in the SQL log
			Colorful:                  false,        // Disable color
		},
	)
	return newLogger
}
func postgresConnection() *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		caches.MetaHandlerServer.MetaHandlerServer.Database.Postgres.Host,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Postgres.User,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Postgres.Password,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Postgres.DB,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Postgres.Port,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Postgres.SSLMode,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Postgres.TimeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			ti, _ := time.LoadLocation(caches.MetaHandlerServer.MetaHandlerServer.Database.TimeZone)
			return time.Now().In(ti)
		},
		Logger: gormLogger(),
	})
	if err != nil {
		errmessage := fmt.Sprintf("Failed to connect database. Error details : %v", err)
		log.Fatal(errmessage)
	}

	return db
}

func mysqlConnection() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
		caches.MetaHandlerServer.MetaHandlerServer.Database.Mysql.User,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Mysql.Password,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Mysql.Host,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Mysql.Port,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Mysql.DB,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Mysql.Charset,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Mysql.ParseTime,
		caches.MetaHandlerServer.MetaHandlerServer.Database.Mysql.Loc)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		NowFunc: func() time.Time {
			ti, _ := time.LoadLocation(caches.MetaHandlerServer.MetaHandlerServer.Database.TimeZone)
			return time.Now().In(ti)
		},
		Logger: gormLogger(),
	})
	if err != nil {
		errmessage := fmt.Sprintf("Failed to connect database. Error details : %v", err)
		log.Fatal(errmessage)
	}

	return db
}

func sqliteConnection() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(caches.MetaHandlerServer.MetaHandlerServer.Database.Sqlite.Location), &gorm.Config{
		NowFunc: func() time.Time {
			ti, _ := time.LoadLocation(caches.MetaHandlerServer.MetaHandlerServer.Database.TimeZone)
			return time.Now().In(ti)
		},
		Logger: gormLogger(),
	})
	if err != nil {
		errmessage := fmt.Sprintf("Failed to connect database. Error details : %v", err)
		log.Fatal(errmessage)
	}

	return db
}
