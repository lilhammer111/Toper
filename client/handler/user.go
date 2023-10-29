package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"to-persist/client/constant"
	"to-persist/client/form"
	"to-persist/client/global"
	"to-persist/client/util"
)

// toper register lilhammer111 -m 12312313212

var (
	user = &form.User{}

	sendSmsOK bool
)

type tokenResponse struct {
	Token string `json:"Token"`
}

func RequestToSendSms(cmd *cobra.Command, args []string) {
	var err error
	user.Mobile, err = cmd.Flags().GetString("mobile")

	var isValid bool
	if isValid = validateMobile(user.Mobile); !isValid {
		fmt.Println("Incorrectly formatted cell phone number.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	resp, err := util.Request(http.MethodGet, global.Config.Url.Base+global.Config.Url.Sms+fmt.Sprintf("?mobile=%s", user.Mobile), nil, false)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	defer resp.Body.Close()

	sendSmsOK = true
}

func Register(cmd *cobra.Command, args []string) {
	if !sendSmsOK {
		os.Exit(1)
	}

	var isValid bool

	if user.Name, isValid = validateAndModifyUsername(args[0]); !isValid {
		fmt.Println("You can't use spaces as usernames.")
		os.Exit(1)
	}

	fmt.Print("ENTER PASSWORD:")

enterPassword:
	pwd, err := gopass.GetPasswd() // This will not echo the password while typing
	if err != nil {
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	if isValid = validatePassword(pwd); !isValid {
		fmt.Println("Passwords should be between 8 and 16 characters long.")
		fmt.Print("ENTER PASSWORD AGAIN:")
		goto enterPassword
	}

	fmt.Print("ENTER SMS Verification Code:")

	_, err = fmt.Scanln(&user.SMSCode)

	userRequest, err := json.Marshal(user)
	if err != nil {
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	resp, err := util.Request(http.MethodPost, global.Config.Url.Base+global.Config.Url.Register, bytes.NewReader(userRequest), false)
	if err != nil {
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	defer resp.Body.Close()

	// Set Token
	respJson, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}
	var tokenResp tokenResponse
	err = json.Unmarshal(respJson, &tokenResp)
	if err != nil {
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}
	viper.Set("token", tokenResp.Token)
	// 将更改写回原配置文件
	err = viper.WriteConfig()
	if err != nil {
		fmt.Println(constant.InternalError)
		zap.S().Error("Error writing to config file:", err)
		os.Exit(1)
	}

	// answer user by status code
	switch resp.StatusCode {
	case http.StatusCreated:
		fmt.Println("Congratulations on your registration! You can now log in to your account using the login command.")
	case http.StatusBadRequest:
		fmt.Println("Bad request. Please check the information you provided.")
	case http.StatusConflict:
		fmt.Println("Conflict. The username or phone number might be already in use.")
	case http.StatusInternalServerError:
		fmt.Println("Server error. Please try again later.")
	default:
		fmt.Println("An unexpected error occurred. Please try again later.")
	}
}

func Login(cmd *cobra.Command, args []string) {
	fmt.Print("Enter password: ")
	password, err := gopass.GetPasswd() // This will not echo the password while typing
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}

	// todo Use the password, e.g., authenticate against a server

	// For this demo, just print a message
	fmt.Println("Password received. (Not displaying it for security reasons.)", password)
}

func Logout(cmd *cobra.Command, args []string) {

}

func validateMobile(mobile string) bool {
	matched, err := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	if err != nil {
		fmt.Println("Can't validate mobile, because", err.Error())
		return false
	}
	return matched
}

func validateAndModifyUsername(username string) (string, bool) {
	res := strings.TrimSpace(username)
	if res == "" {
		return "", false
	}

	return res, true
}

func validatePassword(pwd []byte) bool {
	p := string(pwd)
	if p == "" {
		return false
	}

	if len(p) < 8 {
		return false
	}

	if len(p) > 16 {
		return false
	}

	return true
}