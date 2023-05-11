package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Test_DefaultConfigPath is a dumb test, but we will check that we contain our config in right dir at least
func Test_DefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	assert.Equal(t, "japwords", filepath.Base(filepath.Dir(path)))
}

func Test_EnsureConfigFile(t *testing.T) {
	t.Run("FileExists", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "testfile")
		err := os.WriteFile(path, []byte("hello world"), 0o600)
		assert.NoError(t, err)
		err = EnsureConfigFile(path)
		assert.NoError(t, err)
		data, err := os.ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, []byte("hello world"), data)
	})
	t.Run("FileExistsButItsDir", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "testdir")
		err := os.Mkdir(path, 0o700)
		assert.NoError(t, err)
		err = EnsureConfigFile(path)
		assert.Error(t, err)
	})
	t.Run("FileNotExists", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "testfile")
		err := EnsureConfigFile(path)
		assert.NoError(t, err)
		config, err := LoadConfig(path)
		assert.NoError(t, err)
		assert.Equal(t, DefaultUserConfig(), config)
	})
	t.Run("FileAndDirNotExists", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "dir", "testfile")
		err := EnsureConfigFile(path)
		assert.NoError(t, err)
		config, err := LoadConfig(path)
		assert.NoError(t, err)
		assert.Equal(t, DefaultUserConfig(), config)
	})
}

func Test_SaveConfig(t *testing.T) {
	requireConfig := func(t *testing.T, path string, expected *UserConfig) {
		buffer, err := os.ReadFile(path)
		require.NoError(t, err)
		var actual UserConfig
		err = yaml.Unmarshal(buffer, &actual)
		require.NoError(t, err)
		require.Equal(t, expected, &actual)
	}
	t.Run("Write", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "myfile")
		conf := DefaultUserConfig()
		conf.Dictionary.UserAgent = "myagent"
		err := SaveConfig(path, conf)
		require.NoError(t, err)
		requireConfig(t, path, conf)
	})
	t.Run("Overwrite", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "myfile")
		err := os.WriteFile(path, []byte("hello"), 0o600)
		require.NoError(t, err)
		conf := DefaultUserConfig()
		conf.Dictionary.UserAgent = "myagent"
		err = SaveConfig(path, conf)
		require.NoError(t, err)
		requireConfig(t, path, conf)
	})
}

func Test_LoadConfig(t *testing.T) {
	// TODO: there is no test for env, it's doable, but too combersome
	t.Run("OK", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "myconfig")
		config := &UserConfig{
			Addr: "someaddr",
			Anki: Anki{},
			Dictionary: Dictionary{
				Workers:   4,
				UserAgent: "hello",
				Headers: map[string]string{
					"here": "there",
				},
				Jisho: Jisho{
					URL: "jisho",
				},
				Wadoku: Wadoku{
					URL: "wadoku",
				},
			},
		}
		err := SaveConfig(path, config)
		require.NoError(t, err)
		actual, err := LoadConfig(path)
		require.NoError(t, err)
		require.Equal(t, config, actual)
	})
	t.Run("NoFile", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "notexists")
		_, err := LoadConfig(path)
		require.ErrorContains(t, err, "could not load file config")
	})
	t.Run("IncorrectFile", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "notexists")
		err := os.WriteFile(path, []byte("somenonsens: sdf:"), 0o600)
		require.NoError(t, err)
		_, err = LoadConfig(path)
		require.ErrorContains(t, err, "could not load file config")
	})
	t.Run("IncorrectField", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "notexists")
		err := os.WriteFile(path, []byte("aa: 123"), 0o600)
		require.NoError(t, err)
		_, err = LoadConfig(path)
		require.ErrorContains(t, err, "reading config failed")
	})
}
