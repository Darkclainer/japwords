package config

import (
	"errors"
	"maps"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Manager_New(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "config.yaml")
		config := DefaultUserConfig()
		config.Addr = "myaddr"
		err := SaveConfig(configPath, config)
		require.NoError(t, err)
		mgr, err := New(configPath)
		require.NoError(t, err)
		assert.Equal(t, config, mgr.config)
		assert.Equal(t, configPath, mgr.configPath)
	})
	t.Run("NotExists", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "notconfig.yaml")
		_, err := New(configPath)
		require.Error(t, err)
	})
}

type TestGenericConsumer struct {
	consume func(*UserConfig) (Part, error)
}

func (c *TestGenericConsumer) Config(uc *UserConfig) (Part, error) {
	return c.consume(uc)
}

type TestGenericReloader struct {
	TestGenericConsumer
	reload func(TestGenericPart) error
}

func (c *TestGenericReloader) Reload(p Part) error {
	genericPart, ok := p.(TestGenericPart)
	if !ok {
		return errors.New("reload with wrong part type")
	}
	return c.reload(genericPart)
}

type TestGenericPart map[string]any

func (p TestGenericPart) Equal(other any) bool {
	otherPart, ok := other.(TestGenericPart)
	if !ok {
		return false
	}
	return maps.Equal(p, otherPart)
}

func Test_Manager_Register(t *testing.T) {
	t.Run("ConsumerOK", func(t *testing.T) {
		mgr := newManager(t, &UserConfig{Addr: "myaddr"})
		consumer := &TestGenericConsumer{
			consume: func(uc *UserConfig) (Part, error) {
				data := map[string]any{
					"addr": uc.Addr,
				}
				return TestGenericPart(data), nil
			},
		}
		part, updateFn, err := mgr.Register(consumer)
		require.NoError(t, err)
		assert.Nil(t, updateFn)
		genericPart, ok := part.(TestGenericPart)
		require.True(t, ok)
		assert.Equal(t, "myaddr", genericPart["addr"])
		assert.Len(t, mgr.reloaders, 0)
	})
	t.Run("ConsumerFail", func(t *testing.T) {
		mgr := newManager(t, &UserConfig{Addr: "myaddr"})
		consumer := &TestGenericConsumer{
			consume: func(uc *UserConfig) (Part, error) {
				data := map[string]any{
					"addr": uc.Addr,
				}
				return TestGenericPart(data), errors.New("myerr")
			},
		}
		part, updateFn, err := mgr.Register(consumer)
		require.Error(t, err)
		assert.Nil(t, updateFn)
		genericPart, ok := part.(TestGenericPart)
		require.True(t, ok)
		assert.Equal(t, "myaddr", genericPart["addr"])
	})
	t.Run("ReloaderOK", func(t *testing.T) {
		mgr := newManager(t, &UserConfig{Addr: "myaddr"})
		var initialReloadPart Part
		reloader := &TestGenericReloader{
			TestGenericConsumer: TestGenericConsumer{
				consume: func(uc *UserConfig) (Part, error) {
					data := map[string]any{
						"addr": uc.Addr,
					}
					return TestGenericPart(data), nil
				},
			},
			reload: func(p TestGenericPart) error {
				initialReloadPart = p
				return nil
			},
		}
		part, updateFn, err := mgr.Register(reloader)
		require.NoError(t, err)
		assert.NotNil(t, updateFn)
		assert.Equal(t, part, initialReloadPart)
		genericPart, ok := part.(TestGenericPart)
		require.True(t, ok)
		assert.Equal(t, "myaddr", genericPart["addr"])
		assert.Len(t, mgr.reloaders, 1)
	})
	t.Run("ReloaderDuplicate", func(t *testing.T) {
		mgr := newManager(t, &UserConfig{Addr: "myaddr"})
		reloader := &TestGenericReloader{
			TestGenericConsumer: TestGenericConsumer{
				consume: func(uc *UserConfig) (Part, error) {
					data := map[string]any{
						"addr": uc.Addr,
					}
					return TestGenericPart(data), nil
				},
			},
			reload: func(TestGenericPart) error {
				return nil
			},
		}
		part, updateFn, err := mgr.Register(reloader)
		require.NoError(t, err)
		assert.NotNil(t, updateFn)
		genericPart, ok := part.(TestGenericPart)
		require.True(t, ok)
		assert.Equal(t, "myaddr", genericPart["addr"])
		assert.Len(t, mgr.reloaders, 1)
		part, updateFn, err = mgr.Register(reloader)
		assert.Error(t, err)
		assert.Nil(t, updateFn)
		genericPart, ok = part.(TestGenericPart)
		require.True(t, ok)
		assert.Equal(t, "myaddr", genericPart["addr"])
		assert.Len(t, mgr.reloaders, 1)
	})
}

type TrackedReloaderFactory struct {
	mu     sync.Mutex
	count  int
	called []int
}

func (tr *TrackedReloaderFactory) New(consume func(*UserConfig) (Part, error), reloaderr error) *TestGenericReloader {
	tr.count++
	myindex := tr.count
	return &TestGenericReloader{
		TestGenericConsumer: TestGenericConsumer{
			consume: consume,
		},
		reload: func(TestGenericPart) error {
			tr.mu.Lock()
			defer tr.mu.Unlock()
			tr.called = append(tr.called, myindex)
			return reloaderr
		},
	}
}

func Test_Manager_UpdateConfig(t *testing.T) {
	t.Run("NoChange", func(t *testing.T) {
		mgr := newManager(t, DefaultUserConfig())
		first := fileModTime(t, mgr.configPath)
		// we need some sleep to see that change int modtime
		time.Sleep(5 * time.Millisecond)
		err := mgr.UpdateConfig(func(_ *UserConfig) error { return nil })
		assert.NoError(t, err)
		second := fileModTime(t, mgr.configPath)
		assert.True(t, second.After(first))
		newConfig, err := LoadConfig(mgr.configPath)
		assert.NoError(t, err)
		assert.Equal(t, mgr.config, newConfig)
	})
	t.Run("ChangeWithoutReloaders", func(t *testing.T) {
		mgr := newManager(t, DefaultUserConfig())
		err := mgr.UpdateConfig(func(uc *UserConfig) error {
			uc.Addr = "myaddr"
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, "myaddr", mgr.config.Addr)
		newConfig, err := LoadConfig(mgr.configPath)
		assert.NoError(t, err)
		assert.Equal(t, mgr.config, newConfig)
	})
	t.Run("OneReloader", func(t *testing.T) {
		mgr := newManager(t, DefaultUserConfig())
		factory := &TrackedReloaderFactory{}
		reloader := factory.New(func(uc *UserConfig) (Part, error) {
			data := map[string]any{
				"addr": uc.Addr,
			}
			return TestGenericPart(data), nil
		}, nil)
		_, _, err := mgr.Register(reloader)
		require.NoError(t, err)
		err = mgr.UpdateConfig(func(uc *UserConfig) error {
			uc.Addr = "myaddr"
			return nil
		})
		require.NoError(t, err)
		// two 1 because we should account inital reload on register
		assert.Equal(t, []int{1, 1}, factory.called)
	})
	t.Run("TwoReloaders", func(t *testing.T) {
		mgr := newManager(t, DefaultUserConfig())
		factory := &TrackedReloaderFactory{}
		firstReloader := factory.New(func(uc *UserConfig) (Part, error) {
			data := map[string]any{
				"addr": uc.Addr,
			}
			return TestGenericPart(data), nil
		}, nil)
		secondReloader := factory.New(func(uc *UserConfig) (Part, error) {
			data := map[string]any{
				"addr": uc.Addr,
			}
			return TestGenericPart(data), nil
		}, nil)
		_, _, err := mgr.Register(firstReloader)
		require.NoError(t, err)
		_, _, err = mgr.Register(secondReloader)
		require.NoError(t, err)
		err = mgr.UpdateConfig(func(uc *UserConfig) error {
			uc.Addr = "myaddr"
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2, 1, 2}, factory.called)
	})
	t.Run("TwoReloadersFirstOff", func(t *testing.T) {
		mgr := newManager(t, DefaultUserConfig())
		factory := &TrackedReloaderFactory{}
		firstReloader := factory.New(func(uc *UserConfig) (Part, error) {
			data := map[string]any{
				"ua": uc.Dictionary.UserAgent,
			}
			return TestGenericPart(data), nil
		}, nil)
		secondReloader := factory.New(func(uc *UserConfig) (Part, error) {
			data := map[string]any{
				"addr": uc.Addr,
			}
			return TestGenericPart(data), nil
		}, nil)
		_, _, err := mgr.Register(firstReloader)
		require.NoError(t, err)
		_, _, err = mgr.Register(secondReloader)
		require.NoError(t, err)
		err = mgr.UpdateConfig(func(uc *UserConfig) error {
			uc.Addr = "myaddr"
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2, 2}, factory.called)
	})
	t.Run("TwoReloadersSecondErr", func(t *testing.T) {
		mgr := newManager(t, DefaultUserConfig())
		factory := &TrackedReloaderFactory{}
		firstReloader := factory.New(func(uc *UserConfig) (Part, error) {
			data := map[string]any{
				"addr": uc.Addr,
			}
			return TestGenericPart(data), nil
		}, nil)
		secondReloader := factory.New(func(uc *UserConfig) (Part, error) {
			data := map[string]any{
				"addr": uc.Addr,
			}
			return TestGenericPart(data), nil
		}, errors.New("mytesterror"))
		_, _, err := mgr.Register(firstReloader)
		require.NoError(t, err)
		_, _, err = mgr.Register(secondReloader)
		require.ErrorContains(t, err, "mytesterror")
		err = mgr.UpdateConfig(func(uc *UserConfig) error {
			uc.Addr = "myaddr"
			return nil
		})
		require.ErrorContains(t, err, "mytesterror")
		// First two appears because of register. Third and fourth called on update config, and last
		// because we reload back after error in second reloader
		assert.Equal(t, []int{1, 2, 1, 2, 1}, factory.called)
	})
	t.Run("TwoReloadersAsync", func(t *testing.T) {
		mgr := newManager(t, DefaultUserConfig())
		factory := &TrackedReloaderFactory{}
		firstReloader := factory.New(func(uc *UserConfig) (Part, error) {
			data := map[string]any{
				"addr": uc.Addr,
			}
			return TestGenericPart(data), nil
		}, nil)
		secondReloader := factory.New(func(uc *UserConfig) (Part, error) {
			data := map[string]any{
				"addr": uc.Addr,
			}
			return TestGenericPart(data), nil
		}, nil)
		_, _, err := mgr.Register(firstReloader)
		require.NoError(t, err)
		_, _, err = mgr.Register(secondReloader)
		require.NoError(t, err)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			updateErr := mgr.UpdateConfig(func(uc *UserConfig) error {
				uc.Addr = "myfirstaddr"
				return nil
			})
			require.NoError(t, updateErr)
		}()
		go func() {
			defer wg.Done()
			updateErr := mgr.UpdateConfig(func(uc *UserConfig) error {
				uc.Addr = "mysecondaddr"
				return nil
			})
			require.NoError(t, updateErr)
		}()
		wg.Wait()
		assert.Equal(t, []int{1, 2, 1, 2, 1, 2}, factory.called)
	})
	// this tests chechks that manager config cache is actually updated.
	// There was a bug, when you update config, and then revert it, but reverting
	// did nothing, because manager always kept original config.
	t.Run("ReloaderRevertConfig", func(t *testing.T) {
		uc := DefaultUserConfig()
		uc.Addr = "initialaddr"
		mgr := newManager(t, uc)
		var calledWith []string
		reloader := &TestGenericReloader{
			TestGenericConsumer: TestGenericConsumer{
				consume: func(uc *UserConfig) (Part, error) {
					data := map[string]any{
						"addr": uc.Addr,
					}
					return TestGenericPart(data), nil
				},
			},
			reload: func(p TestGenericPart) error {
				calledWith = append(calledWith, p["addr"].(string))
				return nil
			},
		}

		_, _, err := mgr.Register(reloader)
		require.NoError(t, err)
		err = mgr.UpdateConfig(func(uc *UserConfig) error {
			uc.Addr = "newaddr"
			return nil
		})
		require.NoError(t, err)
		err = mgr.UpdateConfig(func(uc *UserConfig) error {
			uc.Addr = "initialaddr"
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, []string{"initialaddr", "newaddr", "initialaddr"}, calledWith)
	})
	// Just like above, but we check that after error, we revert cache
	t.Run("ReloaderRevertConfigError", func(t *testing.T) {
		uc := DefaultUserConfig()
		uc.Addr = "initialAddr"
		uc.Anki.Addr = "initalAnkiAddr"
		mgr := newManager(t, uc)
		var calledWith []string
		firstReloader := &TestGenericReloader{
			TestGenericConsumer: TestGenericConsumer{
				consume: func(uc *UserConfig) (Part, error) {
					data := map[string]any{
						"addr": uc.Addr,
					}
					return TestGenericPart(data), nil
				},
			},
			reload: func(p TestGenericPart) error {
				calledWith = append(calledWith, p["addr"].(string))
				return nil
			},
		}
		secondReloader := &TestGenericReloader{
			TestGenericConsumer: TestGenericConsumer{
				consume: func(uc *UserConfig) (Part, error) {
					data := map[string]any{
						"ankiaddr": uc.Anki.Addr,
					}
					return TestGenericPart(data), nil
				},
			},
			reload: func(p TestGenericPart) error {
				if p["ankiaddr"] == "newAnkiAddr" {
					return errors.New("hello")
				}
				return nil
			},
		}

		_, _, err := mgr.Register(firstReloader)
		require.NoError(t, err)
		_, _, err = mgr.Register(secondReloader)
		require.NoError(t, err)
		err = mgr.UpdateConfig(func(uc *UserConfig) error {
			uc.Addr = "newAddr"
			uc.Anki.Addr = "newAnkiAddr"
			return nil
		})
		require.Error(t, err)
		// we check first reloader, it should get `newAddr` and then revert to `initialAddr`
		assert.Equal(t, []string{"initialAddr", "newAddr", "initialAddr"}, calledWith)
		err = mgr.UpdateConfig(func(uc *UserConfig) error {
			uc.Addr = "initialAddr"
			return nil
		})
		require.NoError(t, err)
		// second time we update only first reloader and it should not be called, because it was left in `initialAddr` state
		assert.Equal(t, []string{"initialAddr", "newAddr", "initialAddr"}, calledWith)
	})
}

func fileModTime(t testing.TB, path string) time.Time {
	stat, err := os.Stat(path)
	require.NoError(t, err)
	return stat.ModTime()
}

func newManager(t testing.TB, uc *UserConfig) *Manager {
	t.Helper()
	file, err := os.CreateTemp("", "test")
	require.NoError(t, err)
	filepath := file.Name()
	require.NoError(t, file.Close())
	t.Cleanup(func() { os.Remove(filepath) })
	err = SaveConfig(filepath, uc)
	require.NoError(t, err)
	mgr, err := New(filepath)
	require.NoError(t, err)
	return mgr
}
