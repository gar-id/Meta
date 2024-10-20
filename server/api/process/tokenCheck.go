package process

import (
	"crypto/md5"
	"fmt"
	"time"

	"MetaHandler/server/databases"
	userType "MetaHandler/server/databases/types"

	types "MetaHandler/server/api/process/types"

	"MetaHandler/tools"

	"MetaHandler/server/config/caches"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func tokenCheck(c *fiber.Ctx) (accept bool, jsonResult error) {
	// Init variable
	var returnMessage = types.General{
		HTTP_Code: fiber.StatusForbidden,
		Status:    "forbidden",
		ClientIP:  c.IP(),
		Data: struct {
			Date    string "json:\"date\""
			Message string "json:\"message\""
		}{
			Date:    time.Now().String(),
			Message: "Username is not register or static token is invalid",
		},
	}

	clientIP := c.IP()
	var userName string
	if clientIP != "127.0.0.1" {
		// Check from access DB
		var clientData userType.ClientAccess
		// Check system or not
		if c.Get("MetaHandler-User") == "" {
			userName = "system"
		} else if c.Get("MetaHandler-User") != "" {
			userName = c.Get("MetaHandler-User")
		}
		rowEffected, err := databases.ClientAccessCheck(&clientData, clientIP, userName)
		if err != nil {
			tools.ZapLogger("file", "server").Warn(err.Error())
		}
		if rowEffected == 1 {
			// If access is not expired, accept
			if clientData.ExpiredAt.After(time.Now()) {
				return true, nil
			} else {
				return validation(c, returnMessage, clientData, clientIP, userName, rowEffected)
			}
		} else if rowEffected == 0 {
			if c.Get("MetaHandler-Token") == "" {
				returnMessage.Data.Message = "Empty token header"
				return false, c.JSON(returnMessage)
			}
			return validation(c, returnMessage, clientData, clientIP, userName, rowEffected)
		} else if c.Get("MetaHandler-Token") == "" {
			returnMessage.Data.Message = "Empty token header"
			return false, c.JSON(returnMessage)
		}
	}

	// Return accept if from localhost
	return true, nil
}

func validation(c *fiber.Ctx, returnMessage types.General, clientData userType.ClientAccess, clientIP, userName string, rowEffected int64) (accept bool, jsonResult error) {
	// if system or agent
	var tokenCheck error
	if userName == "system" {
		tokenCheck = bcrypt.CompareHashAndPassword([]byte(caches.MetaHandlerServer.MetaHandlerServer.API.MainToken), []byte(c.Get("MetaHandler-Token")))
		if tokenCheck == nil {
			return true, nil
		}
	} else if userName == "centralissh-agent" {
		tokenCheck = bcrypt.CompareHashAndPassword([]byte(caches.MetaHandlerServer.MetaHandlerServer.API.AgentToken), []byte(c.Get("MetaHandler-Token")))
		if tokenCheck == nil {
			return true, nil
		}
	}

	// Get token from DB
	SectionByte := []byte(userName)
	userID := fmt.Sprintf("%x", md5.Sum(SectionByte))
	var dataRetrieve userType.UserData
	databases.GetUser(&dataRetrieve, userID)

	// Return if token is not available
	if dataRetrieve.Token == "" {
		return false, c.JSON(returnMessage)
	} else if dataRetrieve.Token != "" {
		tokenCheck := bcrypt.CompareHashAndPassword([]byte(dataRetrieve.Token), []byte(c.Get("MetaHandler-Token")))
		if tokenCheck == nil {
			if rowEffected == 1 {
				updatedRow, err := databases.ClientAccessUpdate(&clientData, clientIP, userName)
				if err != nil {
					tools.ZapLogger("file", "server").Error(fmt.Sprintf("Row effected is %v with error %v", updatedRow, err.Error()))
				}
			} else if rowEffected == 0 {
				updatedRow, err := databases.ClientAccessAdd(clientIP, userName)
				if err != nil {
					tools.ZapLogger("file", "server").Error(fmt.Sprintf("Row effected is %v with error %v", updatedRow, err.Error()))
				}
			}
			return true, nil
		} else {
			returnMessage.Data.Message = fmt.Sprintf("Invalid token. Err: %v", tokenCheck)
			return false, c.JSON(returnMessage)
		}
	} else {
		return false, c.JSON(returnMessage)
	}
}
