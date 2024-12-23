package config

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// Config構造体
type Config struct {
	// MySQL設定
	DatabaseUser     string `json:"mysql_user"`
	DatabasePassword string `json:"mysql_password"`
	DatabaseHost     string `json:"mysql_host"`
	DatabasePort     string `json:"mysql_port"`
	DatabaseName     string `json:"mysql_name"`
	DatabaseLoc      string `json:"mysql_loc"`
	DatabaseDsn      string `json:"mysql_dsn"`

	WebSocketPort     string `json:"websocket_port"`
	RedisAddr         string `json:"redis_addr"`
	RedisPassword     string `json:"redis_password"`
	FuelPricePerLiter string `json:"fuel_price_per_liter"`
	FuelEfficiency    string `json:"fuel_efficiency"`
}

// LoadConfig関数
func LoadConfig() (*Config, error) {
	var config Config

	// 環境によって設定を読み込む
	environment := os.Getenv("APP_ENV") // "local" または "production"
	// AWS環境の場合、Secrets Managerから読み込む
	if environment != "" && environment != "local" {
		if err := loadSecrets(&config); err != nil {
			return nil, err
		}
	}

	// 環境変数から設定を一括取得し、Configに設定
	envVars := map[string]*string{
		"MYSQL_USER":           &config.DatabaseUser,
		"MYSQL_PASSWORD":       &config.DatabasePassword,
		"MYSQL_HOST":           &config.DatabaseHost,
		"MYSQL_PORT":           &config.DatabasePort,
		"MYSQL_DATABASE":       &config.DatabaseName,
		"MYSQL_LOC":            &config.DatabaseLoc,
		"WEBSOCKET_PORT":       &config.WebSocketPort,
		"REDIS_ADDR":           &config.RedisAddr,
		"REDIS_PASSWORD":       &config.RedisPassword,
		"FUEL_PRICE_PER_LITER": &config.FuelPricePerLiter,
		"FUEL_EFFICIENCY":      &config.FuelEfficiency,
	}

	for key, value := range envVars {
		*value = getEnv(key, *value)
	}

	// DSNを生成
	config.DatabaseDsn = config.getDSN()

	return &config, nil
}

// loadSecrets関数
func loadSecrets(config *Config) error {
	// AWSセッションの作成
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"), // リージョン指定
	})
	if err != nil {
		return err
	}

	svc := secretsmanager.New(sess)

	// 各シークレット名を指定して取得
	secretNames := []string{
		"MYSQL_USER",
		"MYSQL_PASSWORD",
		"MYSQL_HOST",
		"MYSQL_PORT",
		"MYSQL_DATABASE",
		"MYSQL_LOC",
		"WEBSOCKET_PORT",
		"REDIS_ADDR",
		"REDIS_PASSWORD",
		"FUEL_PRICE_PER_LITER",
		"FUEL_EFFICIENCY",
	}

	for _, secretName := range secretNames {
		value, err := getSecret(svc, secretName)
		if err != nil {
			return err
		}
		switch secretName {
		case "MYSQL_USER":
			config.DatabaseUser = value
		case "MYSQL_PASSWORD":
			config.DatabasePassword = value
		case "MYSQL_HOST":
			config.DatabaseHost = value
		case "MYSQL_PORT":
			config.DatabasePort = value
		case "MYSQL_DATABASE":
			config.DatabaseName = value
		case "MYSQL_LOC":
			config.DatabaseLoc = value
		case "WEBSOCKET_PORT":
			config.WebSocketPort = value
		case "REDIS_ADDR":
			config.RedisAddr = value
		case "REDIS_PASSWORD":
			config.RedisPassword = value
		case "FUEL_PRICE_PER_LITER":
			config.FuelPricePerLiter = value
		case "FUEL_EFFICIENCY":
			config.FuelEfficiency = value
		}
	}

	return nil
}

// getSecret関数
func getSecret(svc *secretsmanager.SecretsManager, secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	if result.SecretString != nil {
		return *result.SecretString, nil
	}
	return "", fmt.Errorf("secret %s not found", secretName)
}

// getDSNメソッド
func (c *Config) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName,
	)
}

// getEnv関数
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
