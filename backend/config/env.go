package config

import "github.com/ganeshdipdumbare/goenv"

type EnvVar struct {
	LogLevel               string `json:"log_level"`
	MongoUri               string `json:"mongo_uri"`
	MongoDb                string `json:"mongo_db"`
	ConversationCollection string `json:"conversation_collection"`
	ReviewCollection       string `json:"review_collection"`
	Port                   string `json:"port"`
}

var (
	envVars = &EnvVar{
		LogLevel:               "info",
		Port:                   "5174",
		MongoDb:                "reviewbot",
		ConversationCollection: "conversation",
		ReviewCollection:       "review",
	}
)

func init() {
	goenv.SyncEnvVar(&envVars)
}

func Get() *EnvVar {
	return envVars
}
