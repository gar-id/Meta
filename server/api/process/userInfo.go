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

func UserInfo(c *fiber.Ctx) error {
	// Check token
	isAccept, returnResult := tokenCheck(c)
	if !isAccept {
		return returnResult
	}

	// Parse POST request
	params := c.Queries()

	// Get data from database
	var usersData []userTypes.UserData
	var err error
	var rowEffected int64
	var userData userTypes.UserData
	if params["username"] == "" {
		rowEffected, _ = databases.GetAllUser(&usersData)
	} else if params["username"] != "" {
		// Transform userName into md5
		SectionByte := []byte(params["username"])
		userID := fmt.Sprintf("%x", md5.Sum(SectionByte))
		rowEffected, _ = databases.GetUser(&userData, userID)
		usersData = append(usersData, userData)
	}
	if rowEffected == 0 {
		var returnMessage = types.General{
			HTTP_Code: fiber.StatusNotFound,
			Status:    "not_found",
			ClientIP:  c.IP(),
			Data: struct {
				Date    string "json:\"date\""
				Message string "json:\"message\""
			}{
				Date:    time.Now().String(),
				Message: fmt.Sprintf("User %v not found. User id: %v", params["username"], userData.UserID),
			},
		}
		return c.Status(fiber.StatusNotFound).JSON(returnMessage)
	}

	// Loop for all groupData
	var usersDetails []types.UserDetails
	for _, user := range usersData {
		// Get JSON group
		var getGroup []types.UserGroupJSON
		err = json.Unmarshal(user.Group, &getGroup)
		if err != nil {
			tools.ZapLogger("file", "server").Error(fmt.Sprintf("Error unmarshal. %v", err.Error()))
		}

		var UserDetails = types.UserDetails{
			Username:         user.Username,
			MFAEnabled:       user.MFA,
			Role:             user.Role,
			Group:            getGroup,
			Disabled:         user.Disabled,
			PublicKey:        user.PublicKey,
			PublicKeyExpired: user.PublicKeyExpired,
		}
		usersDetails = append(usersDetails, UserDetails)
	}

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
