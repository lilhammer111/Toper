package global

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"strconv"
	"to-persist/server/config"
)

const (
	ToperStatusUndone = "undone"
	ToperStatusDone   = "done"
)

var (
	MysqlDB     *gorm.DB
	RedisClient *redis.Client

	ServerConfig = &config.ServerConfig{}

	Debugging       bool
	AccessKeyId     string
	AccessKeySecret string
)

func init() {
	var err error
	Debugging, err = strconv.ParseBool(os.Getenv("TOPER_DEBUG"))
	if err != nil {
		zap.S().Panic("failed to convert string to bool type, because ", err.Error())
	}

	AccessKeyId = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")

	AccessKeySecret = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
}
