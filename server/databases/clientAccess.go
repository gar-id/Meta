package databases

import (
	"time"

	clientType "MetaHandler/server/databases/types"
)

func ClientAccessCheck(clientData *clientType.ClientAccess, clientIP, userName string) (int64, error) {
	// Create connection from bootstrap func
	db := switchDB()
	query := db.Where("client_ip = ? AND user_name = ? ", clientIP, userName).Last(&clientData)

	// Execute then close
	sqlDB, _ := db.DB()
	sqlDB.Close()
	return query.RowsAffected, query.Error
}

func ClientAccessAdd(clientIP, userName string) (int64, error) {
	// Create connection from bootstrap func
	db := switchDB()
	var clientData = clientType.ClientAccess{
		ClientIP:  clientIP,
		UserName:  userName,
		ExpiredAt: time.Now().Add(1 * time.Hour),
	}

	// Create query
	query := db.Create(&clientData)

	// Execute then close
	sqlDB, _ := db.DB()
	sqlDB.Close()
	return query.RowsAffected, query.Error
}

func ClientAccessUpdate(clientData *clientType.ClientAccess, clientIP, userName string) (int64, error) {
	// Create connection from bootstrap func
	db := switchDB()
	query := db.Model(&clientData).Where("client_ip = ? AND user_name = ?", clientIP, userName).Update("expired_at", time.Now().Add(1*time.Hour))

	// Execute then close
	sqlDB, _ := db.DB()
	sqlDB.Close()
	return query.RowsAffected, query.Error
}

func ClientAccessDelete(clientData *clientType.ClientAccess, userName string) (int64, error) {
	// Create connection from bootstrap func
	db := switchDB()

	// Update database
	query := db.Unscoped().Where("user_name = ?", userName).Delete(&clientData)

	// Execute then close
	sqlDB, _ := db.DB()
	sqlDB.Close()
	return query.RowsAffected, query.Error
}
