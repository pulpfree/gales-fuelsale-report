package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

// IntegSuite struct
type IntegSuite struct {
	suite.Suite
	c *Config
}

// SetupTest method
func (suite *IntegSuite) SetupTest() {

	suite.c = &Config{}

	os.Setenv("Stage", "test")
	suite.c.setDefaults()
	suite.c.setEnvVars()
}

// TestSetSSMParams function
// this test assumes that the S3Bucket is set
func (suite *IntegSuite) TestSetSSMParams() {

	DBNameBefore := defs.DBName
	err := suite.c.setSSMParams()
	suite.NoError(err)

	DBNameAfter := defs.DBName
	suite.True(strings.Compare(DBNameBefore, DBNameAfter) != 0)
}

// TestLoad function
func (suite *IntegSuite) TestLoad() {

	suite.Empty(suite.c.AWSRegion)

	// cfg, err := Load()
	err := suite.c.Load()
	suite.NoError(err)
	suite.NotEmpty(suite.c.AWSRegion)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
