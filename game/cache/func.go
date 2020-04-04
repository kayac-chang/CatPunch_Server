package cache

func NewCacheSetting() Setting {
	return Setting{}
}

func NewGameCache() *GameCache {
	return &GameCache{
		cacheMap: make(map[string]interface{}),
	}
}
