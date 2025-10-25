package config

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	// use tmp dir for testing
	originalFolder := configFolder
	originalFile := configFile
	configFolder = "/tmp/repo/config_test"
	configFile = "/tmp/repo/config_test/config.json"
	defer func() {
		// other tests shouldn't parse real config files, but just in case reset the paths
		configFolder = originalFolder
		configFile = originalFile
	}()
	// ensure no existing tmp file
	err := os.Remove(configFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("removing existing tmp cfg: %s", err.Error())
	}

	testCases := []struct {
		name      string
		cfgFile   string
		expected  Config
		expectErr bool
	}{
		{
			name:     "default",
			expected: defaultConfig,
		},
		{
			name: "custom config",
			cfgFile: `{
				"local": {
					"root": "/my/custom/root"
				}
			}`,
			expected: Config{
				Local: LocalConfig{
					Root: "/my/custom/root",
				},
			},
		},
		{
			name:      "invalid config",
			cfgFile:   `I AM NOT JSON`,
			expectErr: true,
		},
	}

	mu := new(sync.Mutex) // prevent file race conditions causing test flakes
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mu.Lock()
			defer mu.Unlock()

			if tc.cfgFile != "" {
				require.NoError(t, os.WriteFile(configFile, []byte(tc.cfgFile), 0644), "cfg file setup")
				// cleanup to avoid affecting other tests
				defer func(t *testing.T) {
					assert.NoError(t, os.Remove(configFile))
				}(t)
			}

			cfg, err := Parse()
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, cfg)
		})
	}
}
