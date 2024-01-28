package flags

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var Args FlagsServer = FlagsServer{}

type FlagsServer struct {
	AddrRun            string
	StoreInterval      int
	FileStoragePath    string
	Restore            bool
	DBConnectionString string
	Key                string
}

func (f *FlagsServer) Parse() {
	flag.StringVar(&f.AddrRun, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&f.StoreInterval, "i", 300, "interval for save current metrics data to disk")
	flag.StringVar(&f.FileStoragePath, "f", "./metrics.json", "path file to save current metrics")
	flag.StringVar(&f.DBConnectionString, "d", "", "string for connection to DB, format 'host=%s port=%s user=%s password=%s dbname=%s sslmode=%s'")
	flag.BoolVar(&f.Restore, "r", true, "restore metrics from file")
	flag.StringVar(&f.Key, "k", "", "key for signature request data")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.AddrRun = envRunAddr
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
}
