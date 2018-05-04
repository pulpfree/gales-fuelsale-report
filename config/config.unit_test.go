package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	tAWSRegion = "ca-central-1"
)

// UnitSuite struct
type UnitSuite struct {
	suite.Suite
	cfg *Config
}

// SetupTest method
func (suite *UnitSuite) SetupTest() {
	suite.cfg = &Config{}
}

// TestSetDefaults method
func (suite *UnitSuite) TestSetDefaults() {
	err := suite.cfg.setDefaults()
	suite.NoError(err)
	suite.Equal(tAWSRegion, defs.AWSRegion)
}

// TestSetEnvVars method
func (suite *UnitSuite) TestSetEnvVars() {

	var err error
	err = suite.cfg.setEnvVars()
	suite.NoError(err)

	// Change a var
	os.Setenv("Stage", "noexist")
	err = suite.cfg.setEnvVars()
	suite.Error(err)

	// Reset to valid stage
	os.Setenv("Stage", "test")
	suite.cfg.setEnvVars()

	os.Setenv("DBName", "testdb")
	err = suite.cfg.setEnvVars()
	suite.Equal("testdb", defs.DBName)
}

// TestValidateStage method
func (suite *UnitSuite) TestValidateStage() {
	err := suite.cfg.validateStage()
	suite.NoError(err)
}

// TestSetFinal function
func (suite *UnitSuite) TestSetFinal() {

	var se StageEnvironment
	err := suite.cfg.setFinal()

	suite.NoError(err)
	suite.Equal(suite.cfg.AWSRegion, defs.AWSRegion, "Expected Config.AWSRegion (%s) to equal defs.AWSRegion (%s)", suite.cfg.AWSRegion, defs.AWSRegion)
	suite.Equal(suite.cfg.S3Bucket, defs.S3Bucket, "Expected Config.S3Bucket (%s) to equal defs.S3Bucket (%s)", suite.cfg.S3Bucket, defs.S3Bucket)
	suite.IsType(se, suite.cfg.Stage)
}

// TestSetDBConnectURL function
func (suite *UnitSuite) TestSetDBConnectURL() {

	os.Setenv("DBUser", "test")
	os.Setenv("DBPassword", "test")
	suite.cfg.setEnvVars()
	suite.cfg.setDBConnectURL()

	suite.True(strings.HasPrefix(suite.cfg.DBConnectURL, "mongo"))
	suite.True(strings.Index(suite.cfg.DBConnectURL, "test:test@") > 0)
}

// TestConfigUnit function
func TestConfigUnit(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}
