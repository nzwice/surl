package config

import (
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tcs := []struct {
		raw         string
		expectedCfg AppConfig
		expectedErr bool
	}{
		{
			raw: `listener_addr: ":8080"`,
			expectedCfg: AppConfig{
				HttpAddr: ":8080",
			},
		},
	}
	for idx, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(idx), func(t *testing.T) {
			vp := viper.New()
			vp.SetConfigType("yaml")
			vp.ReadConfig(strings.NewReader(tc.raw))
			var cfg AppConfig
			err := vp.Unmarshal(&cfg)
			if tc.expectedErr {
				require.NotNil(t, err)
			} else {
				assert.EqualValues(t, tc.expectedCfg, cfg)
			}
		})
	}
}
