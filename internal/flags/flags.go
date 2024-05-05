// Пакет для конфигурирования сервера
package flags

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
)

// Args хранит параметры запуска приложения.
var Args FlagsServer = FlagsServer{}

type FlagsServer struct {
	AddrRun            string `json:"address"`                 // Адрес и порт сервера
	AddrRunGRPC        string `json:"address_grpc"`            // Адрес и порт gRPC-сервера
	StoreInterval      int    `json:"store_interval"`          // Интервал сохранения метрик в файл на диске
	FileStoragePath    string `json:"store_file"`              // Путь сохранения метрик в файл
	Restore            bool   `json:"restore"`                 // Загружать метрики из файла при старте приложения
	DBConnectionString string `json:"database_dsn"`            // Строка подключения а БД, формат 'host=%s port=%s user=%s password=%s dbname=%s sslmode=%s'
	Key                string `json:"key,omitempty"`           // Ключ для подписывания данных
	AddrPprof          string `json:"address_pprof,omitempty"` // Адрес и порт для профилировщика
	CryptoKey          string `json:"crypto_key"`              // Путь до файла с приватным ключом
	Config             string `json:"config,omitempty"`        // Путь до файла конфигурации
	TrustedSubnet      string `json:"trusted_subnet"`          // Доверенная подсеть в виде CIDR
}

// Parse разбирает аргументы запуска приложения в переменнную Args.
func (f *FlagsServer) Parse() {
	var storeInterval, restore string

	flag.StringVar(&f.Config, "c", "./config.json", "path to config file")
	flag.StringVar(&f.AddrRun, "a", "", "address and port to run server")
	flag.StringVar(&f.AddrRunGRPC, "ga", "", "address and port to run gRPC-server")
	flag.StringVar(&storeInterval, "i", "", "interval for save current metrics data to disk")
	flag.StringVar(&f.FileStoragePath, "f", "", "path file to save current metrics")
	flag.StringVar(&f.DBConnectionString, "d", "", "string for connection to DB, format 'host=%s port=%s user=%s password=%s dbname=%s sslmode=%s'")
	flag.StringVar(&restore, "r", "", "restore metrics from file")
	flag.StringVar(&f.Key, "k", "", "key for signature request data")
	flag.StringVar(&f.AddrPprof, "p", ":8090", "address and port to pprof server")
	flag.StringVar(&f.CryptoKey, "crypto-key", "", "path to private key file")
	flag.StringVar(&f.TrustedSubnet, "t", "", "trusted subnet (CIDR)")
	flag.Parse()

	f.updateFromConfig(map[string]string{
		"storeInterval": storeInterval,
		"restore":       restore,
	})
	f.updateFromEnvironments()
}

// updateFromConfig обновляет настройки приложения настройками из файла.
func (f *FlagsServer) updateFromConfig(d map[string]string) {
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		f.Config = envConfig
	}

	config := f.getConfigFromFile()

	if f.AddrRun == "" {
		f.AddrRun = config.AddrRun
	}

	if f.AddrRunGRPC == "" {
		f.AddrRunGRPC = config.AddrRunGRPC
	}

	if f.AddrPprof == "" {
		f.AddrPprof = config.AddrPprof
	}

	if d["storeInterval"] == "" {
		f.StoreInterval = config.StoreInterval
	} else {
		v, err := strconv.Atoi(d["storeInterval"])
		if err != nil {
			log.Printf("wrong value flag StoreInterval: %s", d["storeInterval"])
		}

		f.StoreInterval = v
	}

	if f.FileStoragePath == "" {
		f.FileStoragePath = config.FileStoragePath
	}

	if f.DBConnectionString == "" {
		f.DBConnectionString = config.DBConnectionString
	}

	if d["restore"] == "" {
		f.Restore = config.Restore
	} else {
		v, err := strconv.ParseBool(d["restore"])
		if err != nil {
			log.Printf("wrong value flag Restore: %s", d["restore"])
		}

		f.Restore = v
	}

	if f.Key == "" {
		f.Key = config.Key
	}

	if f.CryptoKey == "" {
		f.CryptoKey = config.CryptoKey
	}

	if f.TrustedSubnet == "" {
		f.TrustedSubnet = config.TrustedSubnet
	}
}

// updateFromEnvironments обновляем настройки из переменных среды.
func (f *FlagsServer) updateFromEnvironments() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.AddrRun = envRunAddr
	}

	if envRunAddrGRPC := os.Getenv("ADDRESS_GRPC"); envRunAddrGRPC != "" {
		f.AddrRunGRPC = envRunAddrGRPC
	}

	if envAddrPprof := os.Getenv("ADDRESS_PPROF"); envAddrPprof != "" {
		f.AddrPprof = envAddrPprof
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		f.FileStoragePath = envFileStoragePath
	}

	if envDBConnectionString := os.Getenv("DATABASE_DSN"); envDBConnectionString != "" {
		f.DBConnectionString = envDBConnectionString
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if val, err := strconv.ParseBool(envRestore); err == nil {
			f.Restore = val
		} else {
			log.Printf("wrong value environment RESTORE: %s", envRestore)
		}
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if val, err := strconv.Atoi(envStoreInterval); err == nil {
			f.StoreInterval = val
		} else {
			log.Printf("wrong value environment STORE_INTERVAL: %s", envStoreInterval)
		}
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		f.Key = envKey
	}

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		f.CryptoKey = envCryptoKey
	}

	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		f.TrustedSubnet = envTrustedSubnet
	}
}

// getConfigFromFile загружает файл настроек.
func (f *FlagsServer) getConfigFromFile() *FlagsServer {
	result := &FlagsServer{}

	if f.Config == "" {
		return result
	}

	data, err := os.ReadFile(f.Config)

	if err != nil {
		log.Printf("error opening config file: %s", err)
		return result
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		log.Printf("error in config file: %s", err)
	}

	return result
}

var ArgsClient flagsAgent = flagsAgent{}

type flagsAgent struct {
	AddrRun        string `json:"address"`          // Адрес и порт, куда Агент будет слать данные
	PollInterval   int    `json:"poll_interval"`    // Интервал сбора метрик
	ReportInterval int    `json:"report_interval"`  // Интервал отправки метрик
	RateLimit      int    `json:"rate_limit"`       // Количество воркеров для отправки
	Key            string `json:"key,omitempty"`    // Ключ для подписи данных
	CryptoKey      string `json:"crypto_key"`       // Путь до файла с публичным ключом
	Config         string `json:"config,omitempty"` // Путь до файла конфигурации
	GRPC           bool   `json:"grpc,omitempty"`   // Использовать протокол gRPC
}

// ParseFlags парсит аргументы запуска Агента в переменную flags.
func (f *flagsAgent) ParseFlags() {
	var pollInterval, reportInterval, rateLimit, grpc string

	flag.StringVar(&f.Config, "c", "./config.json", "path to config file")
	flag.StringVar(&f.AddrRun, "a", "", "address and port to run server")
	flag.StringVar(&pollInterval, "p", "", "poll metrics interval")
	flag.StringVar(&reportInterval, "r", "", "send metrics report interval")
	flag.StringVar(&rateLimit, "l", "", "request to server limit")
	flag.StringVar(&f.Key, "k", "", "key for signature request data")
	flag.StringVar(&f.CryptoKey, "crypto-key", "", "path to public key file")
	flag.StringVar(&grpc, "grpc", "", "use gRPC")
	flag.Parse()

	f.updateFromConfig(map[string]string{
		"pollInterval":   pollInterval,
		"reportInterval": reportInterval,
		"rateLimit":      rateLimit,
		"grpc":           grpc,
	})
	f.updateFromEnvironments()
}

// updateFromConfig обновляет настройки приложения настройками из файла.
func (f *flagsAgent) updateFromConfig(d map[string]string) {
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		f.Config = envConfig
	}

	config := f.getConfigFromFile()

	if f.AddrRun == "" {
		f.AddrRun = config.AddrRun
	}

	if d["pollInterval"] == "" {
		f.PollInterval = config.PollInterval
	} else {
		v, err := strconv.Atoi(d["pollInterval"])
		if err != nil {
			log.Printf("wrong value flag PollInterval: %s", d["pollInterval"])
		}

		f.PollInterval = v
	}

	if d["reportInterval"] == "" {
		f.ReportInterval = config.ReportInterval
	} else {
		v, err := strconv.Atoi(d["reportInterval"])
		if err != nil {
			log.Printf("wrong value flag ReportInterval: %s", d["reportInterval"])
		}

		f.ReportInterval = v
	}

	if d["rateLimit"] == "" {
		f.RateLimit = config.RateLimit
	} else {
		v, err := strconv.Atoi(d["rateLimit"])
		if err != nil {
			log.Printf("wrong value flag RateLimit: %s", d["rateLimit"])
		}

		f.RateLimit = v
	}

	if f.Key == "" {
		f.Key = config.Key
	}

	if f.CryptoKey == "" {
		f.CryptoKey = config.CryptoKey
	}

	if d["grpc"] == "" {
		f.GRPC = config.GRPC
	} else {
		v, err := strconv.ParseBool(d["grpc"])
		if err != nil {
			log.Printf("wrong value flag grpc: %s", d["grpc"])
		}

		f.GRPC = v
	}
}

// updateFromEnvironments обновляем настройки из переменных среды.
func (f *flagsAgent) updateFromEnvironments() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.AddrRun = envRunAddr
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if val, err := strconv.Atoi(envPollInterval); err == nil {
			f.PollInterval = val
		} else {
			log.Printf("wrong value environment POLL_INTERVAL: %s", envPollInterval)
		}
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		if val, err := strconv.Atoi(envRateLimit); err == nil {
			f.RateLimit = val
		} else {
			log.Printf("wrong value environment RATE_LIMIT: %s", envRateLimit)
		}
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if val, err := strconv.Atoi(envReportInterval); err == nil {
			f.ReportInterval = val
		} else {
			log.Printf("wrong value environment REPORT_INTERVAL: %s", envReportInterval)
		}
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		f.Key = envKey
	}

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		f.CryptoKey = envCryptoKey
	}

	if envGRPC := os.Getenv("GRPC"); envGRPC != "" {
		if val, err := strconv.ParseBool(envGRPC); err == nil {
			f.GRPC = val
		} else {
			log.Printf("wrong value environment GRPC: %s", envGRPC)
		}
	}
}

// getConfigFromFile загружает настройки из файла
func (f *flagsAgent) getConfigFromFile() *flagsAgent {
	result := &flagsAgent{}

	if f.Config == "" {
		return result
	}

	data, err := os.ReadFile(f.Config)

	if err != nil {
		log.Printf("error opening config file: %s", err)
		return result
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		log.Printf("error in config file: %s", err)
	}

	return result
}
