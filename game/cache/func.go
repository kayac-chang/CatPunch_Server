package cache

func NewCache(setting Setting) *GameCache {
	cache := &GameCache{
		Setting: setting,
	}
	return cache
}
