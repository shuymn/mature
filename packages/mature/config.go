package mature

type Config struct {
	Region         string
	AccessTokenKey string
	DeviceIDKey    string
}

func NewConfig(region, accessTokenKey, deviceIDKey string) *Config {
	return &Config{
		Region:         region,
		AccessTokenKey: accessTokenKey,
		DeviceIDKey:    deviceIDKey,
	}
}
