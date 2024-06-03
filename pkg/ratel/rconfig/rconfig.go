package rconfig

import (
	"ratel/pkg/ratel/key"
)

type ConfigItem struct {
	Capacity            int
	TimeWindowSecond    int64
	BlockDurationSecond int64
}

type ConfigStore struct {
	configs       map[key.RatelKey]ConfigItem
	defaultConfig ConfigItem
}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{
		configs: make(map[key.RatelKey]ConfigItem),
		defaultConfig: ConfigItem{
			Capacity:            10,
			TimeWindowSecond:    1,
			BlockDurationSecond: 10,
		},
	}
}

func (c *ConfigStore) AddConfig(key key.RatelKey, config ConfigItem) *ConfigStore {
	c.configs[key] = config
	return c
}

func (c *ConfigStore) GetConfig(key key.RatelKey) (ConfigItem, bool) {
	config, ok := c.configs[key]
	return config, ok
}

func (c *ConfigStore) AddConfigs(configs map[key.RatelKey]ConfigItem) *ConfigStore {
	for k, v := range configs {
		c.configs[k] = v
	}
	return c
}

func (c *ConfigStore) SetDefaultConfig(config ConfigItem) *ConfigStore {
	c.defaultConfig = config
	return c
}

func (c *ConfigStore) GetDefaultConfig() ConfigItem {
	return c.defaultConfig
}
