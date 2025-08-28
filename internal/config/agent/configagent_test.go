package agent

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	EndpointAgent = flag.String("a", "localhost:8080", "endpoint")
	ReportInterval = flag.Int("r", 10, "report interval")
	PollInterval = flag.Int("p", 2, "poll interval")
	Key = flag.String("k", "", "key")
	RateLimit = flag.Int("l", 1, "rate limit")
}

func setEnv(t *testing.T, key, value string) {
	err := os.Setenv(key, value)
	assert.NoError(t, err)
}

func unsetEnv(t *testing.T, key string) {
	err := os.Unsetenv(key)
	assert.NoError(t, err)
}

func TestConfigAgent_DefaultValues(t *testing.T) {
	resetFlags()
	unsetEnv(t, "ADDRESS")
	unsetEnv(t, "REPORT_INTERVAL")
	unsetEnv(t, "POLL_INTERVAL")
	unsetEnv(t, "KEY")
	unsetEnv(t, "RATE_LIMIT")
	unsetEnv(t, "CRYPTO_KEY")

	os.Args = []string{"cmd"}

	ConfigAgent()

	assert.Equal(t, "localhost:8080", *EndpointAgent)
	assert.Equal(t, 10, *ReportInterval)
	assert.Equal(t, 2, *PollInterval)
	assert.Equal(t, "", *Key)
	assert.Equal(t, 1, *RateLimit)
}

func TestConfigAgent_EnvOverridesDefault(t *testing.T) {
	resetFlags()
	setEnv(t, "ADDRESS", "localhost:9090")
	setEnv(t, "REPORT_INTERVAL", "30")
	setEnv(t, "POLL_INTERVAL", "5")
	setEnv(t, "KEY", "secret_key")
	setEnv(t, "RATE_LIMIT", "10")
	setEnv(t, "CRYPTO_KEY", "cert.pem")

	os.Args = []string{"cmd"}

	ConfigAgent()

	assert.Equal(t, "localhost:9090", *EndpointAgent)
	assert.Equal(t, 30, *ReportInterval)
	assert.Equal(t, 5, *PollInterval)
	assert.Equal(t, "secret_key", *Key)
	assert.Equal(t, 10, *RateLimit)
}

func TestConfigAgent_EnvOverridesFlag(t *testing.T) {
	resetFlags()
	setEnv(t, "ADDRESS", "env.example.com")
	setEnv(t, "REPORT_INTERVAL", "100")
	setEnv(t, "POLL_INTERVAL", "20")
	setEnv(t, "KEY", "env_key")
	setEnv(t, "RATE_LIMIT", "5")

	os.Args = []string{"cmd", "-a=flag.example.com", "-r=60", "-p=3", "-k=flag_key", "-l=2"}

	ConfigAgent()

	assert.Equal(t, "env.example.com", *EndpointAgent)
	assert.Equal(t, 100, *ReportInterval)
	assert.Equal(t, 20, *PollInterval)
	assert.Equal(t, "env_key", *Key)
	assert.Equal(t, 5, *RateLimit)
}

func TestConfigAgent_FlagOverridesDefault(t *testing.T) {
	resetFlags()
	unsetEnv(t, "ADDRESS")
	unsetEnv(t, "REPORT_INTERVAL")
	unsetEnv(t, "POLL_INTERVAL")
	unsetEnv(t, "KEY")
	unsetEnv(t, "RATE_LIMIT")
	unsetEnv(t, "CRYPTO_KEY")

	os.Args = []string{"cmd", "-a=flag.example.com", "-r=60", "-p=3", "-k=flag_key", "-l=2"}

	ConfigAgent()

	assert.Equal(t, "flag.example.com", *EndpointAgent)
	assert.Equal(t, 60, *ReportInterval)
	assert.Equal(t, 3, *PollInterval)
	assert.Equal(t, "flag_key", *Key)
	assert.Equal(t, 2, *RateLimit)
}

func TestConfigAgent_Config(t *testing.T) {
	resetFlags()
	unsetEnv(t, "ADDRESS")
	unsetEnv(t, "REPORT_INTERVAL")
	unsetEnv(t, "POLL_INTERVAL")
	unsetEnv(t, "KEY")
	unsetEnv(t, "RATE_LIMIT")
	unsetEnv(t, "CRYPTO_KEY")

	os.Args = []string{"cmd", "-config=../../../configAgent.json"}

	err := ConfigAgent()
	if err != nil {
		assert.Fail(t, "ConfigAgent() error", err.Error())
	}

	assert.Equal(t, "localhost:8080", *EndpointAgent)
	assert.Equal(t, 1, *ReportInterval)
	assert.Equal(t, 1, *PollInterval)
}

func TestConfigAgent_FlagOverridesConfig(t *testing.T) {
	resetFlags()
	unsetEnv(t, "ADDRESS")
	unsetEnv(t, "REPORT_INTERVAL")
	unsetEnv(t, "POLL_INTERVAL")
	unsetEnv(t, "KEY")
	unsetEnv(t, "RATE_LIMIT")
	unsetEnv(t, "CRYPTO_KEY")

	os.Args = []string{"cmd", "-config=../../../configAgent.json", "-a=flag.example.com", "-r=60", "-p=3", "-k=flag_key", "-l=2", "-crypto-key=asd"}

	err := ConfigAgent()
	if err != nil {
		assert.Fail(t, "ConfigAgent() error", err.Error())
	}

	assert.Equal(t, "flag.example.com", *EndpointAgent)
	assert.Equal(t, 60, *ReportInterval)
	assert.Equal(t, 3, *PollInterval)
	assert.Equal(t, "flag_key", *Key)
	assert.Equal(t, "asd", *CryptoKey)
	assert.Equal(t, 2, *RateLimit)
}

func TestConfigAgent_BadConfig1(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "-config=../../../bad.json"}

	err := ConfigAgent()
	if err == nil {
		assert.Fail(t, "ConfigAgent() error shoud be not nil")
	}
}

func TestConfigAgent_BadConfig2(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "-config=../../../testconfigs/agentTest1.json"}

	err := ConfigAgent()
	if err == nil {
		assert.Fail(t, "ConfigAgent() error shoud be not nil")
	}
}
func TestConfigAgent_BadConfig3(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "-config=../../../testconfigs/agentTest2.json"}

	err := ConfigAgent()
	if err == nil {
		assert.Fail(t, "ConfigAgent() error shoud be not nil")
	}
}
func TestConfigAgent_BadConfig4(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "-config=../../../testconfigs/agentTest3.json"}

	err := ConfigAgent()
	if err == nil {
		assert.Fail(t, "ConfigAgent() error shoud be not nil")
	}
}
