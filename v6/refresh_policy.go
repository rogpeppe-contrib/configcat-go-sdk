package configcat

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sync"
)

const (
	CacheBase = "go_"+ ConfigJsonName +"_%s"
)

type refreshPolicy interface {
	getConfigurationAsync() *asyncResult
	getLastCachedConfig() string
	refreshAsync() *async
	close()
}

type configRefresher struct {
	configFetcher configProvider
	cache         ConfigCache
	logger        Logger
	inMemoryValue string
	cacheKey      string
	sync.RWMutex
}

type RefreshMode interface {
	getModeIdentifier() string
	accept(visitor pollingModeVisitor) refreshPolicy
}

func newConfigRefresher(configFetcher configProvider, cache ConfigCache, logger Logger, sdkKey string) configRefresher {
	sha := sha1.New()
	sha.Write([]byte(sdkKey))
	hash := hex.EncodeToString(sha.Sum(nil))
	cacheKey := fmt.Sprintf(CacheBase, hash)
	return configRefresher{configFetcher: configFetcher, cache: cache, logger: logger, cacheKey: cacheKey}
}

func (refresher *configRefresher) refreshAsync() *async {
	return refresher.configFetcher.getConfigurationAsync().accept(func(result interface{}) {
		response := result.(fetchResponse)
		if result.(fetchResponse).isFetched() {
			refresher.set(response.body)
		}
	})
}

func (refresher *configRefresher) getLastCachedConfig() string {
	return refresher.inMemoryValue
}

// get reads the configuration.
func (refresher *configRefresher) get() string {
	refresher.RLock()
	defer refresher.RUnlock()
	value, err := refresher.cache.Get(refresher.cacheKey)
	if err != nil {
		refresher.logger.Errorf("Reading from the cache failed, %s", err)
		return refresher.inMemoryValue
	}

	return value
}

// set writes the configuration.
func (refresher *configRefresher) set(value string) {
	refresher.Lock()
	defer refresher.Unlock()
	refresher.inMemoryValue = value
	err := refresher.cache.Set(refresher.cacheKey, value)
	if err != nil {
		refresher.logger.Errorf("Saving into the cache failed, %s", err)
	}
}
