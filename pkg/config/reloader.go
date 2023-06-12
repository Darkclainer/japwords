package config

import (
	"fmt"
	"sync"
)

type Manager struct {
	config     *UserConfig
	configPath string

	mu        sync.Mutex
	lastParts map[Reloader]Part
	reloaders []Reloader
}

func New(configPath string) (*Manager, error) {
	uc, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return &Manager{
		config:     uc,
		configPath: configPath,
		lastParts:  map[Reloader]Part{},
	}, nil
}

type Part interface {
	Equal(any) bool
}

type Consumer interface {
	// Config extracts parts of config that will be used.
	// Also its recommended that you implement validation in this function.
	Config(*UserConfig) (Part, error)
}

// ConsumerFunc is helper implementation for Consumer interface
type ConsumerFunc func(*UserConfig) (Part, error)

func (f ConsumerFunc) Config(uc *UserConfig) (Part, error) {
	return f(uc)
}

type Reloader interface {
	Consumer
	Reload(Part) error
}

type UpdateConfigFunc func(updateFn func(*UserConfig) error) error

// Register registers config consumer. It's not concurrent safe.
// If consumer implement Reloader interface then it will be
// registered as consumer that can reload it's self on config change.
// Also Reloader gets it's initial reload with acquired part.
func (m *Manager) Register(consumer Consumer) (Part, UpdateConfigFunc, error) {
	part, err := consumer.Config(m.config.Clone())
	if err != nil {
		return part, nil, err
	}
	reloader, ok := consumer.(Reloader)
	if !ok {
		return part, nil, nil
	}
	if err := m.addReloader(reloader, part); err != nil {
		return part, nil, err
	}
	return part, m.UpdateConfig, nil
}

func (m *Manager) addReloader(reloader Reloader, part Part) error {
	_, ok := m.lastParts[reloader]
	if ok {
		return fmt.Errorf("reloader %T already registered", reloader)
	}
	m.reloaders = append(m.reloaders, reloader)
	if err := reloader.Reload(part); err != nil {
		return fmt.Errorf("initload reload with reloader %T failed: %w", reloader, err)
	}
	m.lastParts[reloader] = part
	return nil
}

// UpdateConfig guarded by simple mutex.
// This function will rewrite config specified while construction Manager.
// Saving can be failed despite reloading all affected parts,
// this indicated by SaveFailedError.
func (m *Manager) UpdateConfig(updateFn func(*UserConfig) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	newConfig := m.config.Clone()
	err := updateFn(newConfig)
	if err != nil {
		return err
	}

	type ReloaderPart struct {
		Reloader Reloader
		OldPart  Part
		NewPart  Part
	}

	var toReload []ReloaderPart
	for _, reloader := range m.reloaders {
		newPart, err := reloader.Config(newConfig)
		if err != nil {
			return err
		}
		oldPart := m.lastParts[reloader]
		// we compare newPart with oldPart in case oldPart not exists (if on register we get error)
		if !newPart.Equal(oldPart) {
			toReload = append(toReload, ReloaderPart{
				Reloader: reloader,
				OldPart:  oldPart,
				NewPart:  newPart,
			})
		}
	}
	// If nothing was updated we can handle it as error, to
	// escape some unwanted behavior, but I will leave it to user
	for i, reloaderPart := range toReload {
		// TODO: load new config in case of error reload old config
		err := reloaderPart.Reloader.Reload(reloaderPart.NewPart)
		if err != nil {
			// we start reloading back starting failed reloader
			for i := i; i >= 0; i-- {
				// Ignore error here. Maybe log it?
				_ = toReload[i].Reloader.Reload(toReload[i].OldPart)
			}
			return err
		}
	}

	m.config = newConfig
	return SaveConfig(m.configPath, newConfig)
}
