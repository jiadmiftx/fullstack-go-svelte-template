package main

import (
	"log"

	"sekolah-madrasah/config"
	"sekolah-madrasah/routes"

	"go.elastic.co/apm/v2"
)

// @title sekolah-madrasah API
// @version 1.0
// @description API for sekolah-madrasah using Clean Architecture
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description API for sekolah-madrasah using Clean Architecture
func main() {
	cfg := config.Init()
	log.Println("✅ Configuration loaded")

	if err := initAPM(cfg); err != nil {
		log.Fatalf("❌ Failed to initialize Elastic APM: %v", err)
	}

	engine := routes.InitRoutes()

	routes.StartServer(engine)
}

func initAPM(cfg *config.AppConfig) error {
	if !cfg.APM.Enabled {
		log.Println("📊 APM is disabled - skipping initialization")
		return nil
	}

	log.Println("🔄 Initializing Elastic APM...")

	// Initialize APM tracer with default options
	// The tracer will automatically read configuration from environment variables
	tracer, err := apm.NewTracer(cfg.APM.ServiceName, "1.0.0")
	if err != nil {
		return err
	}

	// Set the default tracer
	apm.SetDefaultTracer(tracer)

	log.Println("📊 APM initialized successfully - connection will be tested in service validation phase")

	return nil
}
