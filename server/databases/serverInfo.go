package databases

import (
	"MetaHandler/server/databases/types"
	"MetaHandler/tools"
)

func GetAllServerInfo(serversData *[]types.MetaServer) (rowEffected int64, err error) {
	// Create connection from bootstrap func
	db := switchDB()
	runQuery := db.Find(&serversData)

	// Execute then close
	sqlDB, _ := db.DB()
	if runQuery.Error != nil {
		tools.ZapLogger("file", "server").Info(runQuery.Error.Error())
		sqlDB.Close()
		return runQuery.RowsAffected, db.Error
	}
	sqlDB.Close()
	return runQuery.RowsAffected, runQuery.Error
}

func GetServerInfo(serverData *types.MetaServer, ServerID string) (rowEffected int64, err error) {
	// Create connection from bootstrap func
	db := switchDB()
	runQuery := db.Where("server_id = ?", ServerID).Last(&serverData)

	// Execute then close
	sqlDB, _ := db.DB()
	if runQuery.Error != nil {
		tools.ZapLogger("file", "server").Info(runQuery.Error.Error())
		sqlDB.Close()
		return runQuery.RowsAffected, db.Error
	}
	sqlDB.Close()
	return runQuery.RowsAffected, runQuery.Error
}
