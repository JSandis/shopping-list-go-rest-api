package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoDBURI() string {
	loadGoDotEnv()
	return os.Getenv("MONGO_URI")
}

func EnvMongoDBName() string {
	loadGoDotEnv()
	return os.Getenv("MONGO_DATABASE_NAME")
}

func EnvMongoDBCollectionName() string {
	loadGoDotEnv()
	return os.Getenv("MONGO_COLLECTION_NAME")
}

func Port() string {
	loadGoDotEnv()
	return os.Getenv("PORT")
}

func loadGoDotEnv() {
	error := godotenv.Load()
	if error != nil {
		log.Fatal("Error loading .env file")
	}
}
