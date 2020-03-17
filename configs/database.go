package configs

import (
	"log"

	"github.com/go-pg/pg/v9"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

var instanceRedis *redis.Client

type DatabaseConfig struct {
	Username string
	Password string
	Database string
	Address  string
}

func NewDatatabase(databaseConfig *DatabaseConfig) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     databaseConfig.Address,
		User:     databaseConfig.Username,
		Password: databaseConfig.Password,
		Database: databaseConfig.Database,
	})

	_, err := db.Exec("SELECT 1")
	if err != nil {
		log.Fatal("Postgres down", err)
	}

	return db
}

func ConfigDatabase() *DatabaseConfig {
	mode := viper.GetString("mode")
	if mode == "develop" {
		return &DatabaseConfig{
			Username: viper.GetString("postgres-develop.username"),
			Password: viper.GetString("postgres-develop.password"),
			Database: viper.GetString("postgres-develop.database"),
			Address:  viper.GetString("postgres-develop.address"),
		}
	}
	return &DatabaseConfig{
		Username: viper.GetString("postgres-production.username"),
		Password: viper.GetString("postgres-production.password"),
		Database: viper.GetString("postgres-production.database"),
		Address:  viper.GetString("postgres-production.address"),
	}
}

func ConnectCacheDatabase() *redis.Client {
	if instanceRedis == nil {
		client := redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis_address"),
			Password: "",
			DB:       0,
		})

		if _, err := client.Ping().Result(); err != nil {
			log.Fatal("rediserror", err)
		}

		instanceRedis = client
	}
	return instanceRedis
}
