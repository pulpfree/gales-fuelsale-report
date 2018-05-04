package xlsx

import (
	"os"
	"testing"

	"github.com/pulpfree/gales-fuelsale-report/config"
	"github.com/pulpfree/gales-fuelsale-report/model"
	"github.com/pulpfree/gales-fuelsale-report/model/mongo"
	"github.com/stretchr/testify/suite"
)

// UnitSuite struct
type UnitSuite struct {
	suite.Suite
	c  *config.Config
	db model.DBHandler
}

const (
	defaultsFP = "../config/defaults.yaml"
	startDate  = "2018-03-01"
	endDate    = "2018-03-31"
)

// SetupTest method
func (suite *UnitSuite) SetupTest() {
	os.Setenv("Stage", "test")
	suite.c = &config.Config{DefaultsFilePath: defaultsFP}
	err := suite.c.Load()
	suite.NoError(err)

	suite.db, err = mongo.NewDB(suite.c.GetMongoConnectURL())
	suite.NoError(err)
	suite.IsType(new(mongo.DB), suite.db)
}

// TestGetFuelSales method
func (suite *UnitSuite) TestGetFuelSales() {
	sales, err := suite.db.GetFuelSales(startDate, endDate)
	suite.NoError(err)
	suite.True(len(sales.Sales) > 0)
}

// TestOutput method
func (suite *UnitSuite) TestOutput() {

	sales, _ := suite.db.GetFuelSales(startDate, endDate)
	xlsx, err := NewFile(sales)
	suite.NoError(err)
	err = xlsx.SaveAs("../tmp/FuelSalesReport.xlsx")
	suite.NoError(err)
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}
