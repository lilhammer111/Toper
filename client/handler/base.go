package handler

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"os"
	"to-persist/client/constants"
	"to-persist/client/global"
	"to-persist/client/util"
)

func Ping(cmd *cobra.Command, args []string) {
	resp, err := util.Request(http.MethodGet, global.ClientConfig.Url.Root+global.ClientConfig.Url.Ping, nil, false)
	if err != nil {
		zap.S().Error("failed to get response, because ", err.Error())
		fmt.Println("nobody there")
		os.Exit(1)
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("pong")
		os.Exit(1)
	}

	fmt.Println(constants.InternalErrorReply)

	defer resp.Body.Close()

}

func ViewVersion(cmd *cobra.Command, args []string) {
	fmt.Println("v1.0.0")
}
