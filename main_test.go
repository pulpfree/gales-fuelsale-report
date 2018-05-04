package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	validator "gopkg.in/go-playground/validator.v8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pulpfree/gales-fuelsale-report/config"
	"github.com/pulpfree/gales-fuelsale-report/handlers"
	"github.com/pulpfree/gales-fuelsale-report/model/mongo"
	"github.com/pulpfree/gales-fuelsale-report/validators"
	"github.com/stretchr/testify/suite"
)

type response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// UnitSuite struct
type UnitSuite struct {
	suite.Suite
	r   *gin.Engine
	res response
}

// SetupTest method
func (suite *UnitSuite) SetupTest() {
	os.Setenv("Stage", "test")
	c := &config.Config{DefaultsFilePath: defaultsFP}
	err := c.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	db, err := mongo.NewDB(c.GetMongoConnectURL())
	if err != nil {
		log.Fatalf("Error connecting to mongo: %s", err)
	}

	gin.SetMode(gin.ReleaseMode)
	suite.r = gin.Default()
	suite.r.Use(handlers.MiddleDB(db))
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validatedate", validators.ValidateDate)
	}
	suite.r.GET("/", handlers.HeartBeat)
	suite.r.GET("/fuel-sale/:type", handlers.FuelSale)
}

// TestHeartBeatRoute method
func (suite *UnitSuite) TestHeartBeatRoute() {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	suite.r.ServeHTTP(w, req)
	suite.Equal(200, w.Code)

	json.Unmarshal(w.Body.Bytes(), &suite.res)
	suite.True(strings.HasPrefix(suite.res.Message, "Pong"))
}

// TestFuelReportRoute method
func (suite *UnitSuite) TestFuelReportRoute() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/fuel-sale/xlsx?startDate=2018-04-01&endDate=2018-04-30", nil)
	suite.r.ServeHTTP(w, req)
	suite.Equal(200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/fuel-sale/json?startDate=2018-04-01&endDate=2018-04-30", nil)
	suite.r.ServeHTTP(w, req)
	suite.Equal(200, w.Code)

	// Invalid required parameter
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/fuel-sale/xlsx?statDate=2018-04-01&endDate=2018-04-30", nil)
	suite.r.ServeHTTP(w, req)
	suite.Equal(400, w.Code)
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}
