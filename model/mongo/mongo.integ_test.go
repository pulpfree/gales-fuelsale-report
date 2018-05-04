package mongo

import (
	"os"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/pulpfree/gales-fuelsale-report/config"
	"github.com/stretchr/testify/suite"
)

// IntegSuite struct
type IntegSuite struct {
	suite.Suite
	c  *config.Config
	db *DB
}

const (
	defaultsFP = "../../config/defaults.yaml"
	startDate  = "2018-04-01"
	endDate    = "2018-04-30"
)

// SetupTest method
func (suite *IntegSuite) SetupTest() {
	// setup config
	os.Setenv("Stage", "test")
	suite.c = &config.Config{DefaultsFilePath: defaultsFP}
	err := suite.c.Load()
	suite.NoError(err)

	// setup db
	s, err := mgo.Dial(suite.c.GetMongoConnectURL())
	suite.NoError(err)

	suite.db = &DB{session: s}
}

// TestNewDB method
func (suite *IntegSuite) TestNewDB() {
	_, err := NewDB(suite.c.GetMongoConnectURL())
	suite.NoError(err)
}

// TestFetchFuelSales method
func (suite *IntegSuite) TestFetchFuelSales() {

	qp, _ := suite.db.setQueryParams(startDate, endDate)
	fs, err := suite.db.fetchFuelSales(qp)
	suite.NoError(err)
	suite.NotNil(fs)
	suite.True(len(fs.Sales) > 0)
	suite.IsType(time.Time{}, fs.DateStart)
	suite.IsType(time.Time{}, fs.DateEnd)
}

// TestSetQueryParams method
func (suite *IntegSuite) TestSetQueryParams() {

	qp, err := suite.db.setQueryParams(startDate, endDate)
	suite.NoError(err)
	suite.NotNil(qp)
	suite.IsType(time.Time{}, qp.DateEnd)
	suite.IsType(time.Time{}, qp.DateStart)
}

// TestGetStationNames method
func (suite *IntegSuite) TestGetStationNames() {

	nms, err := suite.db.getStationNames()
	suite.NoError(err)
	suite.NotNil(nms)
	suite.True(len(nms) > 0)
}

// TestGetFuelSales method
func (suite *IntegSuite) TestGetFuelSales() {
	sales, err := suite.db.GetFuelSales(startDate, endDate)
	suite.NoError(err)
	suite.True(len(sales.Sales) > 0)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
