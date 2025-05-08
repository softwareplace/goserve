package generator

import (
	"github.com/softwareplace/goserve/cmd/goserve-generator/config"
	"github.com/softwareplace/goserve/cmd/goserve-generator/version"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetTargetVersion(t *testing.T) {
	t.Run("should return default goserve version when config.GoServerVersion was not set", func(t *testing.T) {
		defaultGoServeLatest := version.GoServeLatest
		defer func() {
			version.GoServeLatest = defaultGoServeLatest
		}()

		version.GoServeLatest = func() string {
			return "1.0.0-test-only"
		}

		require.Equal(t, "1.0.0-test-only", getGoServeVersion())
	})

	t.Run("should return provide goserve version when config.GoServerVersion was set", func(t *testing.T) {
		config.GoServerVersion = "1.0.0-test-only"
		require.Equal(t, "1.0.0-test-only", getGoServeVersion())
	})
}
