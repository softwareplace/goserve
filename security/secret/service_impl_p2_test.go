package secret

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/internal/utils"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to generate a temporary RSA key file
func createTempRSAKeyFile(t *testing.T, fileName string) string {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal RSA key: %v", err)
	}

	return createTempKeyFile(t, privBytes, fileName)
}

// Helper function to generate a temporary ECDSA key file
func createTempECDSAKeyFile(t *testing.T, fileName string) string {
	t.Helper()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal ECDSA key: %v", err)
	}

	return createTempKeyFile(t, privBytes, fileName)
}

// Helper function to create a temporary key file
func createTempKeyFile(t *testing.T, keyBytes []byte, fileName string) string {
	t.Helper()

	tmpDir := utils.ProjectBasePath() + "/.log/tmp"
	err := os.MkdirAll(tmpDir, 0700)

	require.NoError(t, err)

	filePath := filepath.Join(tmpDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}

	if err := pem.Encode(file, block); err != nil {
		t.Fatalf("Failed to write PEM block: %v", err)
	}

	return filePath
}

func TestSecretImplValidationP2(t *testing.T) {
	tests := []struct {
		name            string
		setup           func(t *testing.T) string
		wantErr         bool
		expectedKeyType interface{}
	}{
		{
			name: "successful RSA key load",
			setup: func(t *testing.T) string {
				return createTempRSAKeyFile(t, "test_key.pem")
			},
			wantErr:         false,
			expectedKeyType: &rsa.PrivateKey{},
		},
		{
			name: "successful ECDSA key load",
			setup: func(t *testing.T) string {
				return createTempECDSAKeyFile(t, "test_key.pem")
			},
			wantErr:         false,
			expectedKeyType: &ecdsa.PrivateKey{},
		},
		{
			name: "nonexistent file",
			setup: func(t *testing.T) string {
				return "/nonexistent/path/to/key.pem"
			},
			wantErr: true,
		},
		{
			name: "invalid PEM data",
			setup: func(t *testing.T) string {
				tmpDir := utils.ProjectBasePath() + "/.log/tmp"
				filePath := filepath.Join(tmpDir, "invalid.txt")

				if err := os.WriteFile(filePath, []byte("not a valid key"), 0600); err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			wantErr: true,
		},
		{
			name: "unsupported key type",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "unsupported.pem")

				// Create a PEM block with unsupported content
				block := &pem.Block{
					Type:  "PRIVATE KEY",
					Bytes: []byte("invalid key data"),
				}

				file, err := os.Create(filePath)
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()

				if err := pem.Encode(file, block); err != nil {
					t.Fatal(err)
				}

				return filePath
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test file
			keyPath := tt.setup(t)

			// Create test handler
			handler := &apiSecretHandlerImpl[*goservectx.DefaultContext]{
				secretKey: keyPath,
			}

			if tt.wantErr {
				goserveerror.Handler(func() {
					handler.InitAPISecretKey()
					require.Error(t, nil)
				}, func(err error) {
					require.Error(t, err)
				})
			} else {
				handler.InitAPISecretKey()
			}
		})
	}
}
