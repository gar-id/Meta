package databases

import (
	"crypto/md5"
	"fmt"

	"MetaHandler/server/config/caches"
	userType "MetaHandler/server/databases/types"

	"MetaHandler/tools"

	"golang.org/x/crypto/bcrypt"
)

func GetAllUser(userData *[]userType.UserData) (int64, error) {
	// Create connection from bootstrap func
	db := switchDB()
	query := db.Find(&userData)

	// Execute then close
	sqlDB, _ := db.DB()
	sqlDB.Close()
	return query.RowsAffected, query.Error
}

func GetUser(userData *userType.UserData, userID string) (int64, error) {
	// Create connection from bootstrap func
	db := switchDB()
	query := db.Where("user_id = ?", userID).Last(&userData)

	// Execute then close
	sqlDB, _ := db.DB()
	sqlDB.Close()
	return query.RowsAffected, query.Error
}

func InsertUser(userName, TOTPSecret string) error {
	// Create connection from bootstrap func
	db := switchDB()

	// Setup variable to store to DB
	SectionByte := []byte(userName)
	userID := fmt.Sprintf("%x", md5.Sum(SectionByte))
	var dataStore userType.UserData
	if userName == "root" {
		tokenBytes, _ := bcrypt.GenerateFromPassword([]byte(caches.MetaHandlerServer.MetaHandlerServer.API.MainToken), 14)
		dataStore.Role = "root"
		dataStore.Token = string(tokenBytes)
	} else if userName == "centralissh-agent" {
		tokenBytes, _ := bcrypt.GenerateFromPassword([]byte(caches.MetaHandlerServer.MetaHandlerServer.API.AgentToken), 14)
		dataStore.Role = "agent"
		dataStore.Token = string(tokenBytes)
	} else {
		dataStore.Role = "user"
	}
	dataStore.UserID = userID
	dataStore.Username = userName
	dataStore.Disabled = false
	dataStore.TOTPSecret = TOTPSecret

	// Store data to DB
	db.Create(&dataStore)

	// Execute then close
	sqlDB, _ := db.DB()
	sqlDB.Close()
	return db.Error
}

func UpdateUser(userData userType.UserData) (string, error) {
	// Create connection from bootstrap func
	db := switchDB()
	var staticToken string

	// Get data from DB
	var getUser userType.UserData
	GetUser(&getUser, userData.UserID)
	if userData.Role != getUser.Role && userData.Role == "admin" {
		staticToken = fmt.Sprintf("cssh-%v", tools.RandomString(14))
		tokenBytes, _ := bcrypt.GenerateFromPassword([]byte(staticToken), 14)
		userData.Token = string(tokenBytes)
	} else if userData.Username == "root" && userData.Role == "root" {
		staticToken = caches.MetaHandlerServer.MetaHandlerServer.API.MainToken
		tokenBytes, _ := bcrypt.GenerateFromPassword([]byte(staticToken), 14)
		userData.Token = string(tokenBytes)
	} else if userData.Username == "centralissh-agent" && userData.Role == "agent" {
		staticToken = caches.MetaHandlerServer.MetaHandlerServer.API.AgentToken
		tokenBytes, _ := bcrypt.GenerateFromPassword([]byte(staticToken), 14)
		userData.Token = string(tokenBytes)
	} else if userData.Role != "admin" {
		db.Model(&userData).Where("user_id = ?", userData.UserID).Update("token", nil)
		staticToken = "not authorized"
	} else {
		staticToken = "same as before"
	}

	// Update database
	db.Model(&userData).Where("user_id = ?", userData.UserID).Updates(userData)

	sqlDB, _ := db.DB()
	sqlDB.Close()
	return staticToken, db.Error
}

func DeleteUser(userData userType.UserData, userID string) (rowEffected int64, err error) {
	// Create connection from bootstrap func
	db := switchDB()
	sqlQuery := db.Where("user_id = ?", userID).Delete(&userData)

	// Execute then close
	sqlDB, _ := db.DB()
	if sqlQuery.Error != nil {
		tools.ZapLogger("file", "server").Info(sqlQuery.Error.Error())
		sqlDB.Close()
		return sqlQuery.RowsAffected, sqlQuery.Error
	}
	sqlDB.Close()
	return sqlQuery.RowsAffected, sqlQuery.Error
}

func PermanentDeleteUser(userData *userType.UserData) (err error) {
	// Create connection from bootstrap func
	db := switchDB()
	if userData.UserID == "" {
		return db.Error
	} else {
		db.Unscoped().Delete(&userData)
	}

	// Execute then close
	sqlDB, err := db.DB()
	if err != nil {
		tools.ZapLogger("file", "server").Info(err.Error())
		sqlDB.Close()
		return db.Error
	}
	sqlDB.Close()
	return db.Error
}
