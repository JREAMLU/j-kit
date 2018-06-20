package redis

// Load load redis
func Load(consulAddr string, isWatching bool, names ...string) error {
	return LoadConfig(consulAddr, isWatching, names...)
}

// LoadConfig load config
func LoadConfig(consulAddr string, isWatching bool, names ...string) error {
	return nil
}
