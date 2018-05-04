package main

import (
	"os"

	validator "gopkg.in/go-playground/validator.v8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pulpfree/gales-fuelsale-report/config"
	"github.com/pulpfree/gales-fuelsale-report/handlers"
	"github.com/pulpfree/gales-fuelsale-report/model/mongo"
	"github.com/pulpfree/gales-fuelsale-report/validators"

	log "github.com/sirupsen/logrus"
)

const (
	defaultsFP = "./config/defaults.yaml"
)

func main() {

	stg := os.Getenv("Stage")
	if stg == "" {
		os.Setenv("Stage", "test")
	}
	c := &config.Config{DefaultsFilePath: defaultsFP}
	err := c.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	db, err := mongo.NewDB(c.GetMongoConnectURL())
	if err != nil {
		log.Fatalf("Error connecting to mongo: %s", err)
	}
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(handlers.MiddleDB(db))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validatedate", validators.ValidateDate)
	}

	router.GET("/", handlers.HeartBeat)
	router.GET("/fuel-sale/:type", handlers.FuelSale)
	router.Run(":3011")
}

const (
	dateTimeFormat = "2006-01-02 15:04 MST"
	dateFormat     = "2006-01-02"
	timeFormat     = "15:04"
)
