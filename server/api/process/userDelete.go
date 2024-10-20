package process

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"MetaHandler/server/databases"
	userTypes "MetaHandler/server/databases/types"

	types "MetaHandler/server/api/process/types"

	"MetaHandler/tools"

	"github.com/gofiber/fiber/v2"
)

func UserDelete(c *fiber.Ctx) error {
	// Check token
	isAccept, returnResult := tokenCheck(c)
	if !isAccept {
		return returnResult
	}

	// Init variable
	var returnMessage = types.General{
		HTTP_Code: fiber.StatusOK,
		Status:    "success",
		ClientIP:  c.IP(),
		Data: struct {
			Date    string "json:\"date\""
			Message string "json:\"message\""
		}{
			Date:    time.Now().String(),
			Message: "All good",
		},
	}

	// Parse POST request
	var userPayload types.UserDelete
	err := c.BodyParser(&userPayload)
	if err != nil || userPayload.Username == "" {
		if err != nil {
			tools.ZapLogger("file", "server").Info(err.Error())
		}
		returnMessage.HTTP_Code = fiber.StatusBadRequest
		returnMessage.Status = "failed"
		returnMessage.Data.Date = time.Now().String()
		returnMessage.Data.Message = "Invalid payload"
		return c.Status(fiber.StatusBadRequest).JSON(returnMessage)
	} else if !userPayload.Confirm {
		returnMessage.HTTP_Code = fiber.StatusAccepted
		returnMessage.Status = "not_confirmed"
		returnMessage.Data.Date = time.Now().String()
		returnMessage.Data.Message = "Confirm payload is false"
		return c.Status(fiber.StatusAccepted).JSON(returnMessage)
	}

	// Transform userName into md5
	SectionByte := []byte(userPayload.Username)
	UserID := fmt.Sprintf("%x", md5.Sum(SectionByte))
	var userData userTypes.UserData
	databases.GetUser(&userData, UserID)
	rowEffected, err := databases.DeleteUser(userData, UserID)
	if err != nil || rowEffected == 0 {
		var returnMessage = types.General{
			HTTP_Code: fiber.StatusNotFound,
			Status:    "not_found",
			ClientIP:  c.IP(),
			Data: struct {
				Date    string "json:\"date\""
				Message string "json:\"message\""
			}{
				Date:    time.Now().String(),
				Message: fmt.Sprintf("User %v not found", userPayload.Username),
			},
		}
		return c.Status(fiber.StatusNotFound).JSON(returnMessage)
	}
	databases.PermanentDeleteUser(&userData)
	var clientAccess userTypes.ClientAccess
	databases.ClientAccessDelete(&clientAccess, userData.Username)

	// Setup result
	// Get JSON group
	var getGroup []types.UserGroupJSON
	err = json.Unmarshal(userData.Group, &getGroup)
	if err != nil {
		tools.ZapLogger("file", "server").Error(fmt.Sprintf("Error unmarshal. %v", err.Error()))
	}
	var usersDetails = []types.UserDetails{
		{
			Username:         userData.Username,
			MFAEnabled:       userData.MFA,
			Role:             userData.Role,
			Group:            getGroup,
			Disabled:         userData.Disabled,
			PublicKey:        userData.PublicKey,
			PublicKeyExpired: userData.PublicKeyExpired,
		},
	}

	var setReturn = types.UserResponse{
		HTTP_Code: fiber.StatusOK,
		Status:    "success",
		ClientIP:  c.IP(),
		Data: struct {
			Date  string              "json:\"date\""
			Users []types.UserDetails "json:\"users\""
		}{
			Date:  time.Now().String(),
			Users: usersDetails,
		},
	}
	return c.Status(fiber.StatusOK).JSON(setReturn)
}
