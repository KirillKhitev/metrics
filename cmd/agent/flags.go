package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
)

var flags flagsAgent = flagsAgent{}

type flagsAgent struct {
	AddrRun        string `json:"address"`          // Адрес и порт, куда Агент будет слать данные
	PollInterval   int    `json:"poll_interval"`    // Интервал сбора метрик
	ReportInterval int    `json:"report_interval"`  // Интервал отправки метрик
	RateLimit      int    `json:"rate_limit"`       // Количество воркеров для отправки
	Key            string `json:"key,omitempty"`    // Ключ для подписи данных
	CryptoKey      string `json:"crypto_key"`       // Путь до файла с публичным ключом
	Config         string `json:"config,omitempty"` // Путь до файла конфигурации
}

// ParseFlags парсит аргументы запуска Агента в переменную flags.
func (f *flagsAgent) ParseFlags() {
	var pollInterval, reportInterval, rateLimit string

	flag.StringVar(&f.Config, "c", "./config.json", "path to config file")
	flag.StringVar(&f.AddrRun, "a", "", "address and port to run server")
	flag.StringVar(&pollInterval, "p", "", "poll metrics interval")
	flag.StringVar(&reportInterval, "r", "", "send metrics report interval")
	flag.StringVar(&rateLimit, "l", "", "request to server limit")
	flag.StringVar(&f.Key, "k", "", "key for signature request data")
	flag.StringVar(&f.CryptoKey, "crypto-key", "", "path to public key file")
	flag.Parse()

	f.updateFromConfig(map[string]string{
		"pollInterval":   pollInterval,
		"reportInterval": reportInterval,
		"rateLimit":      rateLimit,
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
