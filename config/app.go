package config

import (
	"cfxWorld/lib/crawler"
	"cfxWorld/lib/moonswap"
	conflux "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var C App
var ConfPath string

type App struct {
	NodeURL string `yaml:"nodeURL"`
	ScanURL string `yaml:"scanURL"`
	Client  Client `yaml:"client"`
	Wallet  Wallet `yaml:"service"`
	Tx      Tx     `yaml:"tx"`
}

type Wallet struct {
	KeyStorePath string `yaml:"keyStorePath"`
}

type Client struct {
	RequestTimeout       time.Duration `yaml:"requestTimeout"`
	RequestRetryInterval time.Duration `yaml:"requestRetryInterval"`
	RequestMaxRetry      int           `yaml:"requestMaxRetry"`
}

type Tx struct {
	MaxWait time.Duration `yaml:"maxWait"`
}

func LoadConfig(file string, cfgPtr interface{}, dirs ...string) error {
	if cf, err := os.Stat(file); err == nil && !cf.IsDir() {
		log.Println("load cfg file:", file)
		viper.SetConfigFile(file)
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		return viper.Unmarshal(cfgPtr)
	}

	viper.SetConfigName(file)
	for _, d := range dirs {
		viper.AddConfigPath(d)
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(cfgPtr)
}

func MustInitConfig() {
	if err := LoadConfig(ConfPath, &C, ".", "./config"); err != nil {
		log.Println(err)
		log.Fatal("load config failed, please check your config file")
	}
	//init LazyLoad object
	crawler.LazyLoad(C.NodeURL, conflux.ClientOption{
		RetryCount:     C.Client.RequestMaxRetry,
		RetryInterval:  C.Client.RequestRetryInterval,
		RequestTimeout: C.Client.RequestTimeout,
	})
	moonswap.LazyLoad(C.NodeURL, conflux.ClientOption{
		RetryCount:     C.Client.RequestMaxRetry,
		RetryInterval:  C.Client.RequestRetryInterval,
		RequestTimeout: C.Client.RequestTimeout,
	})
}
