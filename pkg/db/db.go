package db

import (
	"nbim/configs"
	"nbim/pkg/logger"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	DB       *gorm.DB
	RedisCli *redis.Client
)

func init() {
	InitMysql(configs.GlobalConfig.Mysql)
	InitRedis(configs.GlobalConfig.RedisHost, configs.GlobalConfig.RedisPassword)
}

func InitMysql(dataSource string) {
	logger.Logger.Info("init mysql begin")
	var err error
	//"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	DB.SingularTable(true)
	DB.LogMode(true)
	logger.Logger.Info("init mysql done!")
}

func InitRedis(addr, password string) {
	logger.Logger.Info("init redis begin")
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       0,
		Password: password,
	})

	_, err := RedisCli.Ping().Result()
	if err != nil {
		panic(err)
	}
	logger.Logger.Info("init redis done!")
}

// 使用json的序列化方式对redis客户端set操作
func SetRedisByJson(key string, value interface{}, duration time.Duration) error {
	byte, err := sonic.Marshal(value)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	err = RedisCli.Set(key, byte, duration).Err()
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// 使用json的序列化方式对redis客户端get操作
func GetRedisByJson(key string, value interface{}) error {
	byte, err := RedisCli.Get(key).Bytes()
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	err = sonic.Unmarshal(byte, &value)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}
