package notify

type Config struct {
	Region        string
	SlackTokenKey string
	ChannelIDKey  string
}

func NewConfig(region, slackTokenKey, channelIDKey string) *Config {
	return &Config{
		Region:        region,
		SlackTokenKey: slackTokenKey,
		ChannelIDKey:  channelIDKey,
	}
}
