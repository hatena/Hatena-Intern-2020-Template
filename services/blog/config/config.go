package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// Config は各種設定をまとめたもの
type Config struct {
	Mode                  string
	Port                  int
	DatabaseDSN           string
	AccountAddr           string
	AccountECDSAPublicKey *ecdsa.PublicKey
	RendererAddr          string
	GracefulStopTimeout   time.Duration
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	conf := &Config{
		Mode:                "production",
		Port:                8080,
		GracefulStopTimeout: 10 * time.Second,
	}

	// Mode
	mode := os.Getenv("MODE")
	if mode != "" {
		conf.Mode = mode
	}

	portStr := os.Getenv("PORT")
	if portStr != "" {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return nil, fmt.Errorf("Invalid PORT: %v", err)
		}
		conf.Port = port
	}

	// DatabaseDSN
	databaseDSN := os.Getenv("DATABASE_DSN")
	if databaseDSN == "" {
		return nil, errors.New("DATABASE_DSN is not set")
	}
	conf.DatabaseDSN = databaseDSN

	// AccountAddr
	accountAddr := os.Getenv("ACCOUNT_ADDR")
	if accountAddr == "" {
		return nil, errors.New("ACCOUNT_ADDR is not set")
	}
	conf.AccountAddr = accountAddr

	// AccountECDSAPublicKey
	accountECDSAPublicKeyFile := os.Getenv("ACCOUNT_ECDSA_PUBLIC_KEY_FILE")
	if accountECDSAPublicKeyFile == "" {
		return nil, errors.New("ECDSA_PRIVATE_KEY_FILE is not set")
	}
	accountECDSAPublicKey, err := loadECDSAPublicKey(accountECDSAPublicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key from %s", accountECDSAPublicKeyFile)
	}
	conf.AccountECDSAPublicKey = accountECDSAPublicKey

	// RendererAddr
	rendererAddr := os.Getenv("RENDERER_ADDR")
	if rendererAddr == "" {
		return nil, errors.New("RENDERER_ADDR is not set")
	}
	conf.RendererAddr = rendererAddr

	// GracefulStopTimeout
	gracefulStopTimeout := os.Getenv("GRACEFUL_STOP_TIMEOUT")
	if gracefulStopTimeout != "" {
		d, err := time.ParseDuration(gracefulStopTimeout)
		if err != nil {
			return nil, fmt.Errorf("GRACEFUL_STOP_TIMEOUT is invalid: %v", err)
		}
		conf.GracefulStopTimeout = d
	}

	return conf, nil
}

func loadECDSAPublicKey(file string) (*ecdsa.PublicKey, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to load public key")
	}
	rawKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := rawKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to load public key")
	}
	return key, nil
}
