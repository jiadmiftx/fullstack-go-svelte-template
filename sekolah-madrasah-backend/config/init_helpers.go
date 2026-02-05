package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

func loadEnvFile() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Infof("Warning: .env file found but failed to load: %v", err)
			log.Infof("Using environment variables or default values")
		} else {
			log.Infof("✅ .env file loaded successfully")
		}
	} else {
		log.Infof("No .env file found, using environment variables or default values")
	}
}

func newAppConfig() *AppConfig {
	return &AppConfig{
		Version:     "12 December 2025 - 02:17:48 WIB",
		ServiceName: "UMKM Service",
		Rest: Rest{
			Port:            getEnvAsInt("REST_PORT", 8080),
			Origin:          getEnv("REST_CORS_ORIGIN", "*"),
			RestDebugMode:   getEnvAsBool("REST_DEBUG_MODE", true),
			JWTSecret:       getEnv("REST_SECRET", ""),
			XAuthCron:       getEnv("X_AUTH_CRON", ""),
			SwaggerUser:     getEnv("SWAGGER_USER", ""),
			SwaggerPassword: getEnv("SWAGGER_PASSWORD", ""),
		},
		MainDB: DBConfig{
			Id:       MainDB,
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnvAsInt("DB_PORT", 5432),
			SslMode:  getEnv("DB_SSLMODE", ""),
			Username: getEnv("DB_USERNAME", ""),
			Name:     getEnv("DB_NAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			DbEvent:  getEnv("DB_EVENT", ""),
		},
		APM: APMConfig{
			Enabled:          getEnvAsBool("ELASTIC_APM_ENABLED", true),
			ServiceName:      getEnv("ELASTIC_APM_SERVICE_NAME", ""),
			ServerURL:        getEnv("ELASTIC_APM_SERVER_URL", ""),
			SecretToken:      getEnv("ELASTIC_APM_SECRET_TOKEN", ""),
			VerifyServerCert: getEnvAsBool("ELASTIC_APM_VERIFY_SERVER_CERT", true),
		},
		// Broker: MsgBrokerConfig{
		// 	Host:     getEnv("MQTT_HOST", ""),
		// 	Port:     getEnv("MQTT_PORT", "1883"),
		// 	User:     getEnv("MQTT_USER", ""),
		// 	Password: getEnv("MQTT_PASSWORD", ""),
		// },
		// Bucket: BucketConfig{
		// 	BaseUrl:         getEnv("AWS_S3_BASE_URL", ""),
		// 	BucketName:      getEnv("AWS_S3_BUCKET_NAME", ""),
		// 	AccessKeyId:     getEnv("AWS_S3_ACCESS_KEY_ID", ""),
		// 	SecretAccessKey: getEnv("AWS_S3_SECRET_ACCESS_KEY", ""),
		// 	PrefixPath:      getEnv("AWS_S3_PREFIX_PATH", ""),
		// 	Region:          getEnv("AWS_S3_REGION", ""),
		// },
		Elasticsearch: ElasticsearchConfig{
			Enabled:     getEnvAsBool("ELASTIC_ENABLED", false),
			Addresses:   []string{getEnv("ELASTIC_HOST", "http://localhost:9200")},
			Username:    getEnv("ELASTIC_USERNAME", ""),
			Password:    getEnv("ELASTIC_PASSWORD", ""),
			APIKey:      getEnv("ELASTIC_API_KEY", ""),
			CloudID:     getEnv("ELASTIC_CLOUD_ID", ""),
			IndexPrefix: getEnv("ELASTIC_INDEX", ""),
			CertPath:    getEnv("ELASTIC_CERT_PATH", ""),
		},
		Security: Security{
			AesKey: getEnv("AES_KEY", "default-aes-key-32-chars-long!!"),
		},
	}
}

func validateOperationalConfig(app *AppConfig) {
	// validateBrokerConfig(app)
	validateMainDBConfig(app)
	if app.Elasticsearch.Enabled {
		validateElasticsearchConfig(app)
	} else {
		log.Info("Elasticsearch disabled, skipping validation")
	}
	validateAPMConfig(app)
}

// func validateBrokerConfig(app *AppConfig) {
// 	if app.Broker.Host == "" {
// 		log.Fatalf(" You need to set Broker Host First !")
// 	}
// 	if app.Broker.User == "" {
// 		log.Fatalf("You need to set Broker User ")
// 	}
// }

func validateMainDBConfig(app *AppConfig) {
	if app.MainDB.Host == "" {
		log.Fatal("you need to set DB_HOST")
	}
	if app.MainDB.Name == "" {
		log.Fatal("you need to set DB_NAME")
	}
	if app.MainDB.Username == "" {
		log.Fatalf("You need to set DB_USERNAME ")
	}
	// DB_PASSWORD can be empty for local development with peer authentication
}

func validateElasticsearchConfig(app *AppConfig) {
	if app.Elasticsearch.Addresses == nil || len(app.Elasticsearch.Addresses) == 0 {
		log.Fatal("you need to set ELASTIC_HOST (or ELASTICSEARCH_ADDRESSES)")
	}

	// Check if any authentication method is provided
	hasUsername := app.Elasticsearch.Username != ""
	hasPassword := app.Elasticsearch.Password != ""
	hasAPIKey := app.Elasticsearch.APIKey != ""
	hasCloudID := app.Elasticsearch.CloudID != ""

	if !hasUsername && !hasAPIKey && !hasCloudID {
		log.Fatal("you need to set ELASTIC_USERNAME and ELASTIC_PASSWORD, or ELASTIC_API_KEY, or ELASTIC_CLOUD_ID for Elasticsearch authentication")
	}

	// If username is provided, password should also be provided (unless using API key)
	if hasUsername && !hasPassword && !hasAPIKey {
		log.Fatal("you need to set ELASTIC_PASSWORD when using ELASTIC_USERNAME")
	}

	// Validate certificate path if provided
	if app.Elasticsearch.CertPath != "" {
		if _, err := os.Stat(app.Elasticsearch.CertPath); os.IsNotExist(err) {
			log.Fatalf("ELASTIC_CERT_PATH is set but certificate file not found: %s", app.Elasticsearch.CertPath)
		}
	}

	// Log which authentication method is being used
	if hasCloudID {
		log.Infof("✅ Using Elastic Cloud authentication (Cloud ID: %s)", app.Elasticsearch.CloudID)
	} else if hasAPIKey {
		log.Infof("✅ Using Elasticsearch API key authentication")
	} else if hasUsername && hasPassword {
		log.Infof("✅ Using Elasticsearch username/password authentication (user: %s)", app.Elasticsearch.Username)
	}

	// Log certificate usage
	if app.Elasticsearch.CertPath != "" {
		log.Infof("✅ Using custom Elasticsearch certificate: %s", app.Elasticsearch.CertPath)
	}

	// Validate index prefix
	if app.Elasticsearch.IndexPrefix == "" {
		log.Fatal("you need to set ELASTIC_INDEX")
	}
}

func validateAPMConfig(app *AppConfig) {
	if app.APM.Enabled {
		if app.APM.ServerURL == "" {
			log.Fatal("APM_ENABLED is true but ELASTIC_APM_SERVER_URL is not set")
		}
		if app.APM.SecretToken == "" {
			log.Fatal("APM_ENABLED is true but ELASTIC_APM_SECRET_TOKEN is not set")
		}
		if app.APM.ServiceName == "" {
			log.Fatal("APM_ENABLED is true but ELASTIC_APM_SERVICE_NAME is not set")
		}

		log.Infof("✅ Elastic APM configuration validated:")
		log.Infof("   - Server URL: %s", app.APM.ServerURL)
		log.Infof("   - Service Name: %s", app.APM.ServiceName)

		// Log certificate verification setting
		if app.APM.VerifyServerCert {
			log.Info("✅ APM server certificate verification enabled")
		} else {
			log.Warn("⚠️  APM server certificate verification disabled")
		}

		log.Info("🔄 APM connection will be tested during application startup...")
	} else {
		log.Info("APM monitoring disabled")
	}
}

func validateRestAndSecurity(app *AppConfig) {
	if app.Rest.JWTSecret == "" {
		log.Fatal("you neet to set REST_SECRET to start rest server")
	}
	if app.Security.AesKey == "" {
		log.Fatal("you need to set AES_KEY env")
	}
}
