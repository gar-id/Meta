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

func UserUpdate(c *fiber.Ctx) error {
	// Check token
	isAccept, returnResult := tokenCheck(c)
	if !isAccept {
		return returnResult
	}

	// Parse POST request
	var userPayload types.UserDetails
	err := c.BodyParser(&userPayload)
	if err != nil {
		tools.ZapLogger("file", "server").Info(err.Error())
		var returnMessage = types.General{
			HTTP_Code: fiber.StatusBadRequest,
			Status:    "failed",
			ClientIP:  c.IP(),
			Data: struct {
				Date    string "json:\"date\""
				Message string "json:\"message\""
			}{
				Date:    time.Now().String(),
				Message: "Invalid payload",
			},
		}
		return c.Status(fiber.StatusBadRequest).JSON(returnMessage)
	}

	// Get data from database
	var userData userTypes.UserData
	// Generate userId
	SectionByte := []byte(userPayload.Username)
	userID := fmt.Sprintf("%x", md5.Sum(SectionByte))
	rowEffected, err := databases.GetUser(&userData, userID)
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
				Message: "User not found",
			},
		}
		return c.Status(fiber.StatusNotFound).JSON(returnMessage)
	}
	// Get JSON group
	var getGroup []types.UserGroupJSON
	err = json.Unmarshal(userData.Group, &getGroup)
	if err != nil {
		tools.ZapLogger("file", "server").Error(fmt.Sprintf("Error unmarshal. %v", err.Error()))
	}

	// Update user to database
	userData.MFA = userPayload.MFAEnabled
	if userData.PublicKey != userPayload.PublicKey && userPayload.PublicKey != "" {
		userData.PublicKeyExpired = time.Now().Add(24 * time.Hour * 30 * 3)
		userData.PublicKey = userPayload.PublicKey
	}
	userData.Role = userPayload.Role
	userData.Disabled = userPayload.Disabled
	var userGroups []types.UserGroupJSON
	if len(userPayload.Group) > 0 {
		jsonGroup, err := json.Marshal(userPayload.Group)
		if err != nil {
			tools.ZapLogger("file", "server").Info(err.Error())
		}
		userGroups = append(userGroups, userPayload.Group...)
		userData.Group = jsonGroup
	} else {
		userGroups = getGroup
	}
	staticToken, _ := databases.UpdateUser(userData)
	var usersDetails []types.UserDetails
	if userData.Role != "admin" || userData.Role != "root" {
		var clientAccess userTypes.ClientAccess
		databases.ClientAccessDelete(&clientAccess, userData.Username)
	}

	// Set variable for output
	userDetails := types.UserDetails{
		Username:         userData.Username,
		MFAEnabled:       userData.MFA,
		Role:             userData.Role,
		Group:            userGroups,
		Disabled:         userData.Disabled,
		PublicKey:        userData.PublicKey,
		PublicKeyExpired: userData.PublicKeyExpired,
		Token:            staticToken,
	}
	usersDetails = append(usersDetails, userDetails)
	var returnMessage = types.UserResponse{
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
	return c.Status(fiber.StatusOK).JSON(returnMessage)
}
