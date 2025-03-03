package configs

import (
	"flag"
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd",
		"-a", "127.0.0.1:8080", // address
		"-i", "5", // store interval
		"-f", "path", // storage path
		"-n", "fileName", // file name
		"-r", "true", // restore
	}

	// сбрасываем флаги
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	config, _ := NewConfig()

	assert.Equal(t, "127.0.0.1:8080", config.RunAddress)
	assert.Equal(t, int64(5), config.StoreInterval)
	assert.Equal(t, "path", config.FileStoragePath)
	assert.Equal(t, "fileName", config.FileStorageName)
	assert.True(t, config.Restore)
}
