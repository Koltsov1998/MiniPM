package MiniPM

type Config struct {
	PostgresqlUser     string `yaml:"postgresql_user" env:"POSTGRESQL_USER" default:"postgres"`
	PostgresqlPassword string `yaml:"postgresql_password" env:"POSTGRESQL_PASSWORD" default:"postgres"`
	PostgresqlHost     string `yaml:"postgresql_host" env:"POSTGRESQL_HOST" default:"localhost"`
	PostgresqlPort     int    `yaml:"postgresql_port" env:"POSTGRESQL_PORT" default:"5432"`
	Port               string `yaml:"port" env:"PORT" default:"8080"`
	DefaultSchedule    string `yaml:"default_schedule" env:"DEFAULT_SCHEDULE"`
	Reddy              struct {
		BotDomain string `yaml:"bot_domain" env:"REDDY_BOT_DOMAIN"`
		BotToken  string `yaml:"bot_token" env:"REDDY_BOT_TOKEN"`
	} `yaml:"reddy"`
	YouTrack struct {
		Url   string `yaml:"url" env:"YOUTRACK_URL"`
		Token string `yaml:"token" env:"YOUTRACK_TOKEN"`
	}
}
