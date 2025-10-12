package MiniPM

type Config struct {
	Port            string `yaml:"port" env:"PORT" default:"8080"`
	DefaultSchedule string `yaml:"default_schedule" env:"DEFAULT_SCHEDULE"`
}
