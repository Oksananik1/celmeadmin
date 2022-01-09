package config

import "os"

type Config struct {
	Port          string
	MongoURI      string
	DBName        string
	Base          string
	FilePath      string
	TelegramToken string
}

func configLookup(key string) (string, bool) {
	return os.LookupEnv("CELME_" + key)
}

func (conf *Config) defval() {
	conf.Port = ":8080"
	conf.Base = "/celmeapi"
	conf.DBName = "celme"
	conf.FilePath = "/Users/Admin/Downloads/celme/"
}

func (conf *Config) Env() {
	conf.defval()
	if port, ok := configLookup("PORT"); ok {
		conf.Port = port
	}
	if host, ok := configLookup("MONGOURI"); ok {
		conf.MongoURI = host
	}
	if name, ok := configLookup("DBNAME"); ok {
		conf.DBName = name
	}

	if storage, ok := configLookup("STORAGE_PATH"); ok {
		conf.FilePath = storage
	}

	if telegramToken, ok := configLookup("TELEGRAM_TOKEN"); ok {
		conf.TelegramToken = telegramToken
	}

}
func (conf *Config) listen() {

}
