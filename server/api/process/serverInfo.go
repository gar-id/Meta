package process

import (
	"crypto/md5"
	"fmt"
	"time"

	"MetaHandler/server/api/process/types"
	"MetaHandler/server/databases"
	dataTypes "MetaHandler/server/databases/types"
	"MetaHandler/tools"

	"github.com/gofiber/fiber/v2"
)

func GetServerInfo(c *fiber.Ctx) error {
	// Parse POST request
	params := c.Queries()

	// Get data from database
	var serversData []dataTypes.MetaServer
	var serverData dataTypes.MetaServer
	var rowEffected int64
	var err error
	if params["host_ip"] == "" {
		// Get all hosts
		rowEffected, err = databases.GetAllServerInfo(&serversData)
	} else if params["host_ip"] != "" {
		// Transform userName into md5
		sectionByte := []byte(params["public_ip"])
		serverID := fmt.Sprintf("%x", md5.Sum(sectionByte))
		rowEffected, err = databases.GetServerInfo(&serverData, serverID)
	}

	// Check result
	if err != nil || rowEffected == 0 {
		if err != nil {
			tools.ZapLogger("file", "server").Error(err.Error())
		}
		var returnMessage = types.General{
			HTTP_Code: fiber.StatusNotFound,
			Status:    "not_found",
			ClientIP:  c.IP(),
			Data: struct {
				Date    string "json:\"date\""
				Message string "json:\"message\""
			}{
				Date:    time.Now().String(),
				Message: fmt.Sprintf("Meta server data with IP public %v not found", tools.DefaultString(params["public_ip"], "all host")),
			},
		}
		return c.Status(fiber.StatusNotFound).JSON(returnMessage)
	}

	// Store data to variable
	serversData = append(serversData, serverData)
	var serversDetails []types.ServerDetails
	for _, server := range serversData {
		if server.ServerID != "" {
			var serverDetails = types.ServerDetails{
				Hostname:         server.Hostname,
				PublicIP:         server.PublicIP,
				Environment:      server.Environment,
				Status:           server.Status,
				ActiveConnection: server.ActiveConnection,
			}

			serversDetails = append(serversDetails, serverDetails)
		}
	}

	// Setup result
	var response = types.ServerResponse{
		HTTP_Code: fiber.StatusOK,
		Status:    "success",
		ClientIP:  c.IP(),
		Data: struct {
			Date    string                "json:\"date\""
			Servers []types.ServerDetails "json:\"servers\""
		}{
			Date:    time.Now().String(),
			Servers: serversDetails,
		},
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
