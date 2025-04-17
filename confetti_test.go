package confetti_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eriktate/confetti"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	String     string `conf:"TEST_NAME"`
	Bool       bool   `conf:"TEST_BOOL"`
	Int        int    `conf:"TEST_INT"`
	Uint       uint   `conf:"TEST_UINT"`
	ByteSlice  []byte `conf:"TEST_BYTE_SLICE"`
	DefaultKey string
}

func TestApplyEnv(t *testing.T) {
	os.Setenv("TEST_NAME", "test")
	os.Setenv("TEST_BOOL", "true")
	os.Setenv("TEST_INT", "-42")
	os.Setenv("TEST_UINT", "42")
	os.Setenv("TEST_BYTE_SLICE", "bytes")
	os.Setenv("DefaultKey", "default")

	cfg := testConfig{}
	err := confetti.ApplyEnv(&cfg)
	require.NoError(t, err)

	require.Equal(t, "test", cfg.String)
	require.Equal(t, true, cfg.Bool)
	require.Equal(t, -42, cfg.Int)
	require.Equal(t, uint(42), cfg.Uint)
	require.Equal(t, []byte("bytes"), cfg.ByteSlice)
	require.Equal(t, "default", cfg.DefaultKey)
}

func TestApplyFiles(t *testing.T) {
	content1 := `TEST_NAME=test
TEST_BOOL=true
TEST_INT=-42
TEST_UINT=10`
	content2 := `TEST_UINT=42
TEST_BYTE_SLICE=bytes
DefaultKey=default`

	path1 := filepath.Join(t.TempDir(), ".env1")
	path2 := filepath.Join(t.TempDir(), ".env2")

	file, err := os.Create(path1)
	require.NoError(t, err)
	defer file.Close()
	_, err = file.WriteString(content1)
	require.NoError(t, err)

	file, err = os.Create(path2)
	require.NoError(t, err)
	defer file.Close()
	_, err = file.WriteString(content2)
	require.NoError(t, err)

	cfg := testConfig{}
	err = confetti.ApplyFiles(&cfg, path1, path2)
	require.NoError(t, err)

	require.Equal(t, "test", cfg.String)
	require.Equal(t, true, cfg.Bool)
	require.Equal(t, -42, cfg.Int)
	require.Equal(t, uint(42), cfg.Uint)
	require.Equal(t, []byte("bytes"), cfg.ByteSlice)
	require.Equal(t, "default", cfg.DefaultKey)
}

func TestFromEnv(t *testing.T) {
	os.Setenv("TEST_NAME", "test")
	os.Setenv("TEST_BOOL", "true")
	os.Setenv("TEST_INT", "-42")
	os.Setenv("TEST_UINT", "42")
	os.Setenv("TEST_BYTE_SLICE", "bytes")
	os.Setenv("DefaultKey", "default")

	cfg, err := confetti.FromEnv[testConfig]()
	require.NoError(t, err)

	require.Equal(t, "test", cfg.String)
	require.Equal(t, true, cfg.Bool)
	require.Equal(t, -42, cfg.Int)
	require.Equal(t, uint(42), cfg.Uint)
	require.Equal(t, []byte("bytes"), cfg.ByteSlice)
	require.Equal(t, "default", cfg.DefaultKey)
}

func TestFromFiles(t *testing.T) {
	content1 := `TEST_NAME=test
TEST_BOOL=true
TEST_INT=-42
TEST_UINT=10`
	content2 := `TEST_UINT=42
TEST_BYTE_SLICE=bytes
DefaultKey=default`

	path1 := filepath.Join(t.TempDir(), ".env1")
	path2 := filepath.Join(t.TempDir(), ".env2")

	file, err := os.Create(path1)
	require.NoError(t, err)
	defer file.Close()
	_, err = file.WriteString(content1)
	require.NoError(t, err)

	file, err = os.Create(path2)
	require.NoError(t, err)
	defer file.Close()
	_, err = file.WriteString(content2)
	require.NoError(t, err)

	cfg, err := confetti.FromFiles[testConfig](path1, path2)
	require.NoError(t, err)

	require.Equal(t, "test", cfg.String)
	require.Equal(t, true, cfg.Bool)
	require.Equal(t, -42, cfg.Int)
	require.Equal(t, uint(42), cfg.Uint)
	require.Equal(t, []byte("bytes"), cfg.ByteSlice)
	require.Equal(t, "default", cfg.DefaultKey)
}
