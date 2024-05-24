package config

import (
	"booker/internal/infrastructure/logging"
	"os"
	"strconv"
	"time"
)

type Storage string

const (
	Memory   Storage = "memory"
	External Storage = "external"
)

type Config struct {
	Debug                         bool
	IsPrepareAvailabilityRequired bool
	ServerPort                    int
	ShutdownTimeoutSec            time.Duration //nolint:stylecheck
	StorageType                   Storage
}

func NewConfig() Config {
	shutdownTimeoutRaw := os.Getenv("SHUTDOWN_TIMEOUT_SEC")
	shutdownTimeout, err := strconv.ParseInt(shutdownTimeoutRaw, 10, 64)
	if err != nil {
		shutdownTimeout = 5 // fallback
		logging.LogErrorf("parse shurdown timeout failed, fallback to %d", shutdownTimeout)
	}

	storageType := Storage(os.Getenv("STORAGE_TYPE"))
	switch storageType {
	case Memory, External:
	default:
		storageType = Memory // fallback
		logging.LogErrorf("invalid storage type, fallback to %s", storageType)
	}

	portRaw := os.Getenv("SERVER_PORT")
	port, err := strconv.Atoi(portRaw)
	if err != nil {
		port = 8080
		logging.LogErrorf("parse server port failed, fallback to %d", port)
	}

	return Config{
		Debug:                         os.Getenv("DEBUG") == "true",
		IsPrepareAvailabilityRequired: os.Getenv("IS_PREPARE_AVAILABILITY_REQUIRED") == "true",
		ServerPort:                    port,
		ShutdownTimeoutSec:            time.Duration(shutdownTimeout) * time.Second,
		StorageType:                   storageType,
	}
}
