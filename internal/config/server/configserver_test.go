package server

import (
	"flag"
	"os"
	"sync"
	"testing"

	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)

	EndpointServer = flag.String("a", "localhost:8080", "endpoint")
	StoreInterval = flag.Int("i", 300, "store interval")
	FileStorePath = flag.String("f", "filestore.out", "file store path")
	Restore = flag.Bool("r", true, "restore")
	DatabaseDSN = flag.String("d", "", "database_DSN")
	Key = flag.String("k", "", "key")
	IsDB = false
}

func setEnv(t *testing.T, key, value string) {
	require.NoError(t, os.Setenv(key, value))
}

func unsetEnv(t *testing.T, key string) {
	require.NoError(t, os.Unsetenv(key))
}

var loggerOnce sync.Once

func TestMain(m *testing.M) {
	loggerOnce.Do(func() {
		logger.Initialize("info")
	})

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestConfigServer_DefaultValues(t *testing.T) {
	resetFlags()
	unsetEnv(t, "ADDRESS")
	unsetEnv(t, "STORE_INTERVAL")
	unsetEnv(t, "FILE_STORAGE_PATH")
	unsetEnv(t, "RESTORE")
	unsetEnv(t, "DATABASE_DSN")
	unsetEnv(t, "KEY")

	os.Args = []string{"cmd"}

	ConfigServer()

	assert.Equal(t, "localhost:8080", *EndpointServer)
	assert.Equal(t, 300, *StoreInterval)
	assert.Equal(t, "filestore.out", *FileStorePath)
	assert.True(t, *Restore)
	assert.Empty(t, *DatabaseDSN)
	assert.False(t, IsDB)
	assert.Empty(t, *Key)
}
func TestConfigServer_EnvOverridesDefault(t *testing.T) {
	resetFlags()
	setEnv(t, "ADDRESS", "example.com:9090")
	setEnv(t, "STORE_INTERVAL", "60")
	setEnv(t, "FILE_STORAGE_PATH", "/tmp/store.out")
	setEnv(t, "RESTORE", "false")
	setEnv(t, "DATABASE_DSN", "postgres://...")
	setEnv(t, "KEY", "secret_key")

	os.Args = []string{"cmd"}

	ConfigServer()

	assert.Equal(t, "example.com:9090", *EndpointServer)
	assert.Equal(t, 60, *StoreInterval)
	assert.Equal(t, "/tmp/store.out", *FileStorePath)
	assert.False(t, *Restore)
	assert.Equal(t, "postgres://...", *DatabaseDSN)
	assert.True(t, IsDB)
	assert.Equal(t, "secret_key", *Key)
}

func TestConfigServer_EnvOverridesFlag(t *testing.T) {
	resetFlags()
	setEnv(t, "ADDRESS", "env.example.com")
	setEnv(t, "STORE_INTERVAL", "100")
	setEnv(t, "FILE_STORAGE_PATH", "/tmp/env_store.out")
	setEnv(t, "RESTORE", "false")
	setEnv(t, "DATABASE_DSN", "env_postgres://...")
	setEnv(t, "KEY", "env_secret")

	os.Args = []string{
		"cmd",
		"-a=flag.example.com",
		"-i=10",
		"-f=/tmp/flag_store.out",
		"-r=false",
		"-d=flag_postgres://...",
		"-k=flag_secret",
	}

	ConfigServer()

	assert.Equal(t, "env.example.com", *EndpointServer)
	assert.Equal(t, 100, *StoreInterval)
	assert.Equal(t, "/tmp/env_store.out", *FileStorePath)
	assert.False(t, *Restore)
	assert.Equal(t, "env_postgres://...", *DatabaseDSN)
	assert.True(t, IsDB)
	assert.Equal(t, "env_secret", *Key)
}

func TestConfigServer_FlagOverridesDefault(t *testing.T) {
	resetFlags()
	unsetEnv(t, "ADDRESS")
	unsetEnv(t, "STORE_INTERVAL")
	unsetEnv(t, "FILE_STORAGE_PATH")
	unsetEnv(t, "RESTORE")
	unsetEnv(t, "DATABASE_DSN")
	unsetEnv(t, "KEY")

	os.Args = []string{
		"cmd",
		"-a=flag.example.com",
		"-i=10",
		"-f=/tmp/flag_store.out",
		"-r=false",
		"-d=flag_postgres://...",
		"-k=flag_secret",
	}

	ConfigServer()

	assert.Equal(t, "flag.example.com", *EndpointServer)
	assert.Equal(t, 10, *StoreInterval)
	assert.Equal(t, "/tmp/flag_store.out", *FileStorePath)
	assert.False(t, *Restore)
	assert.Equal(t, "flag_postgres://...", *DatabaseDSN)
	assert.True(t, IsDB)
	assert.Equal(t, "flag_secret", *Key)
}

func TestConfigServer_EmptyDatabaseDSN(t *testing.T) {
	resetFlags()
	unsetEnv(t, "DATABASE_DSN")

	os.Args = []string{"cmd"}
	ConfigServer()

	assert.Empty(t, *DatabaseDSN)
	assert.False(t, IsDB)
}
