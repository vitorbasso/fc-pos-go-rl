package ratel

import (
	"log"
	"os"
	"ratel/pkg/ratel/key"
	"ratel/pkg/ratel/rconfig"
	"ratel/pkg/ratel/request"
	"ratel/pkg/ratel/request/local"
	"ratel/pkg/ratel/request/rredis"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Option func(*Ratel)

func WithEnvRules() Option {
	tokenString, _ := os.LookupEnv("RATEL_TOKEN_CONFIG")
	ipString, _ := os.LookupEnv("RATEL_IP_CONFIG")
	tokenConfigs := getTokenConfig(tokenString)
	ipConfigs := getIPConfig(ipString)
	return func(r *Ratel) {
		r.configStore = r.configStore.AddConfigs(tokenConfigs)
		r.configStore = r.configStore.AddConfigs(ipConfigs)
	}
}

func getIPConfig(configString string) map[key.RatelKey]rconfig.ConfigItem {
	if configString == "" {
		return make(map[key.RatelKey]rconfig.ConfigItem)
	}
	configStringPerIP := strings.Split(configString, " ")
	configItems := make(map[key.RatelKey]rconfig.ConfigItem, len(configStringPerIP))
	for _, ipConfigString := range configStringPerIP {
		ipConfig := strings.Split(ipConfigString, ",")
		k := key.IPKey(ipConfig[0])
		item, err := getConfig(ipConfig[1], ipConfig[2], ipConfig[3])
		if err != nil {
			log.Printf("error parsing ip config: %v. config string: %s", err, ipConfig)
			continue
		}
		configItems[k] = item
	}
	return configItems
}

func getTokenConfig(configString string) map[key.RatelKey]rconfig.ConfigItem {
	if configString == "" {
		return make(map[key.RatelKey]rconfig.ConfigItem)
	}
	configStringPerToken := strings.Split(configString, " ")
	configItems := make(map[key.RatelKey]rconfig.ConfigItem, len(configStringPerToken))
	for _, tokenConfigString := range configStringPerToken {
		tokenConfig := strings.Split(tokenConfigString, ",")
		k := key.TokenKey(tokenConfig[0])
		item, err := getConfig(tokenConfig[1], tokenConfig[2], tokenConfig[3])
		if err != nil {
			log.Printf("error parsing token config: %v. config string: %s", err, tokenConfig)
			continue
		}
		configItems[k] = item
	}
	return configItems
}
func getConfig(capacity, timeWindowSecond, blockDurationSecond string) (rconfig.ConfigItem, error) {
	capacityInt, err := strconv.Atoi(capacity)
	if err != nil {
		return rconfig.ConfigItem{}, err
	}
	timeWindowSecondInt, err := strconv.Atoi(timeWindowSecond)
	if err != nil {
		return rconfig.ConfigItem{}, err
	}
	blockDurationSecondInt, err := strconv.Atoi(blockDurationSecond)
	if err != nil {
		return rconfig.ConfigItem{}, err
	}
	return rconfig.ConfigItem{
		Capacity:            capacityInt,
		TimeWindowSecond:    int64(timeWindowSecondInt),
		BlockDurationSecond: int64(blockDurationSecondInt),
	}, nil
}

func WithDefaultRule(item rconfig.ConfigItem) Option {
	return func(r *Ratel) {
		r.configStore.SetDefaultConfig(item)
	}
}

func WithRuleItem(key key.RatelKey, item rconfig.ConfigItem) Option {
	return func(r *Ratel) {
		r.configStore.AddConfig(key, item)
	}
}

func WithRequestStore(store request.RequestStore) Option {
	return func(r *Ratel) {
		r.requestStore = store
	}
}

func WithLocalRequestStore() Option {
	return func(r *Ratel) {
		r.requestStore = local.NewLocalStore()
	}
}

func WithRedisRequestStore(options *redis.Options) Option {
	client := redis.NewClient(options)
	return func(r *Ratel) {
		r.requestStore = rredis.NewRedisStore(client)
	}
}
