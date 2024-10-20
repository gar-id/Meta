package process

import (
	"crypto/md5"
	"fmt"
	"time"

	"MetaHandler/server/databases"
	userTypes "MetaHandler/server/databases/types"

	types "MetaHandler/server/api/process/types"

	"MetaHandler/tools"

	"github.com/gofiber/fiber/v2"
	"github.com/xlzd/gotp"
)

func TOTPInit(c *fiber.Ctx) error {
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
	var userPayload types.TOTP
	err := c.BodyParser(&userPayload)
	if err != nil {
		tools.ZapLogger("file", "server").Info(err.Error())
		returnMessage.HTTP_Code = fiber.StatusBadRequest
		returnMessage.Status = "failed"
		returnMessage.Data.Date = time.Now().String()
		returnMessage.Data.Message = "Invalid payload"
		return c.Status(fiber.StatusBadRequest).JSON(returnMessage)
	}

	// Get data from database
	var userData userTypes.UserData
	// Setup variable to store to DB
	SectionByte := []byte(userPayload.Username)
	userID := fmt.Sprintf("%x", md5.Sum(SectionByte))
	rowEffected, _ := databases.GetUser(&userData, userID)

	if rowEffected != 1 {
		// Process to get TOTP and store userData to database
		randomSecret := gotp.RandomSecret(32)
		totpUri := generateTOTPWithSecret(randomSecret, userPayload.Username)
		err := databases.InsertUser(userPayload.Username, randomSecret)
		if err != nil {
			tools.ZapLogger("file", "server").Error(err.Error())
		}

		// Return to user
		var returnTOTP = types.TOTPResponse{
			HTTP_Code: fiber.StatusAccepted,
			Status:    "accepted",
			ClientIP:  c.IP(),
			Data: struct {
				Date    string "json:\"date\""
				Secret  string "json:\"secret\""
				Message string "json:\"message\""
				TOTPURI string "json:\"totp_uri\""
			}{
				Date:    time.Now().String(),
				Secret:  randomSecret,
				Message: "User and TOTP created",
				TOTPURI: totpUri,
			},
		}
		return c.Status(fiber.StatusAccepted).JSON(returnTOTP)
	}

	// Check database result
	if userData.Disabled {
		// Return to user
		returnMessage.HTTP_Code = fiber.StatusForbidden
		returnMessage.Status = "forbidden"
		returnMessage.Data.Date = time.Now().String()
		returnMessage.Data.Message = fmt.Sprintf("%v is disabled. Contact MetaHandler admin to get more info", userPayload)
		return c.Status(fiber.StatusForbidden).JSON(returnMessage)
	} else if !userData.MFA {
		// Process to get TOTP and store userData to database
		randomSecret := gotp.RandomSecret(32)
		totpUri := generateTOTPWithSecret(randomSecret, userPayload.Username)
		userData.TOTPSecret = randomSecret
		databases.UpdateUser(userData)

		// Return to user
		var returnTOTP = types.TOTPResponse{
			HTTP_Code: fiber.StatusAccepted,
			Status:    "accepted",
			ClientIP:  c.IP(),
			Data: struct {
				Date    string "json:\"date\""
				Secret  string "json:\"secret\""
				Message string "json:\"message\""
				TOTPURI string "json:\"totp_uri\""
			}{
				Date:    time.Now().String(),
				Secret:  randomSecret,
				Message: "User and TOTP created",
				TOTPURI: totpUri,
			},
		}
		return c.Status(fiber.StatusAccepted).JSON(returnTOTP)
	}

	// Check TOTP Code
	if verifyOTP(userData.TOTPSecret, userPayload.TOTPCode, userPayload.Username) {
		// Return to user
		returnMessage.HTTP_Code = fiber.StatusOK
		returnMessage.Status = "success"
		returnMessage.Data.Date = time.Now().String()
		returnMessage.Data.Message = "Success to verify 2FA with TOTP"
	} else {
		// Return to user
		returnMessage.HTTP_Code = fiber.StatusForbidden
		returnMessage.Status = "forbidden"
		returnMessage.Data.Date = time.Now().String()
		returnMessage.Data.Message = "Wrong TOTP code"
	}

	return c.Status(fiber.StatusForbidden).JSON(returnMessage)
}

func generateTOTPWithSecret(randomSecret, userName string) (uri string) {
	// Generate URI for TOTP
	uri = gotp.NewDefaultTOTP(randomSecret).ProvisioningUri(fmt.Sprintf("%v@centralissh", userName), "MetaHandler")
	tools.ZapLogger("file", "server").Debug(fmt.Sprintf("TOTP has created for user %v with URL %v. TOTP code: %v", userName, uri, randomSecret))
	return uri
}

func verifyOTP(randomSecret, totpCode, userName string) bool {
	totp := gotp.NewDefaultTOTP(randomSecret)
	tools.ZapLogger("file", "server").Debug(fmt.Sprintf("%v real TOTP when %v is %v", userName, (time.Now().Unix() - 3), totp.Now()))
	// Validate the provided OTP
	if totp.Now() == totpCode {
		tools.ZapLogger("file", "server").Debug(fmt.Sprintf("Authentication for user %v success", userName))
		return true
	} else {
		tools.ZapLogger("file", "server").Debug(fmt.Sprintf("Authentication for user %v failed", userName))
		return false
	}
}
