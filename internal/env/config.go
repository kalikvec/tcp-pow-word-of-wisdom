package env

type Config struct {
	ServerHost            string `envconfig:"SERVER_HOST"`
	ServerPort            int    `envconfig:"SERVER_PORT"`
	CacheHost             string `envconfig:"CACHE_HOST"`
	CachePort             int    `envconfig:"CACHE_PORT"`
	HashcashZerosCount    int
	ChallengeTimeout      int64
	HashcashMaxIterations int
}
