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
	StoreInterval      int    `json:"store_interval"`          // Интервал сохранения метрик в файл на диске
	FileStoragePath    string `json:"store_file"`              // Путь сохранения метрик в файл
	Restore            bool   `json:"restore"`                 // Загружать метрики из файла при старте приложения
	DBConnectionString string `json:"database_dsn"`            // Строка подключения а БД, формат 'host=%s port=%s user=%s password=%s dbname=%s sslmode=%s'
	Key                string `json:"key,omitempty"`           // Ключ для подписывания данных
	AddrPprof          string `json:"address_pprof,omitempty"` // Адрес и порт для профилировщика
	CryptoKey          string `json:"crypto_key"`              // Путь до файла с приватным ключом
	Config             string `json:"config,omitempty"`        // Путь до файла конфигурации
}

// Parse разбирает аргументы запуска приложения в переменнную Args.
func (f *FlagsServer) Parse() {
	var storeInterval, restore string

	flag.StringVar(&f.Config, "c", "./config.json", "path to config file")
	flag.StringVar(&f.AddrRun, "a", "", "address and port to run server")
	flag.StringVar(&storeInterval, "i", "", "interval for save current metrics data to disk")
	flag.StringVar(&f.FileStoragePath, "f", "", "path file to save current metrics")
	flag.StringVar(&f.DBConnectionString, "d", "", "string for connection to DB, format 'host=%s port=%s user=%s password=%s dbname=%s sslmode=%s'")
	flag.StringVar(&restore, "r", "", "restore metrics from file")
	flag.StringVar(&f.Key, "k", "", "key for signature request data")
	flag.StringVar(&f.AddrPprof, "p", ":8090", "address and port to pprof server")
	flag.StringVar(&f.CryptoKey, "crypto-key", "", "path to private key file")
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
}

// updateFromEnvironments обновляем настройки из переменных среды.
func (f *FlagsServer) updateFromEnvironments() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.AddrRun = envRunAddr
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
