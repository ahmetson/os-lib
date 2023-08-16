package path

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing orchestra
type TestPathSuite struct {
	suite.Suite
	envPath  string // using env as a file to check
	fileName string
	dir      string
}

// Make sure that Account is set to five
// before each test
func (suite *TestPathSuite) SetupTest() {
	suite.fileName = ".test.env"
	suite.dir = "configs"
	suite.envPath = filepath.Join("configs", ".test.env")

	currentDir, err := CurrentDir()
	suite.Require().NoError(err)

	err = MakeDir(filepath.Join(currentDir, suite.dir))
	suite.Require().NoError(err)

	file, err := os.OpenFile(filepath.Join(currentDir, suite.envPath), os.O_CREATE|os.O_RDWR, 0755)
	suite.Require().NoError(err)

	_, err = file.WriteString("KEY=123")
	suite.Require().NoError(err, "failed to write the data into: "+suite.envPath)
	err = file.Close()
	suite.Require().NoError(err, "delete the dump file: "+suite.envPath)
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *TestPathSuite) TestRun() {
	currentDir, err := CurrentDir()
	suite.Require().NoError(err)
	suite.Require().True(filepath.IsAbs(currentDir))

	// AbsDir
	expected := filepath.Join(currentDir, suite.envPath)
	absPath := AbsDir(currentDir, suite.envPath)
	suite.Require().Equal(expected, absPath)

	// FileName
	suite.Require().Equal(suite.fileName, FileName(absPath))

	// DirAndFileName
	actualDir, actualFileName := DirAndFileName(absPath)
	suite.Require().Equal(filepath.Join(currentDir, suite.dir), actualDir)
	suite.Require().Equal(suite.fileName, actualFileName)

	// FileExist
	fmt.Printf("abs path: %s\n", absPath)
	exist, err := FileExist(absPath)
	suite.Require().NoError(err)
	suite.Require().True(exist)

	_, err = FileExist(actualDir)
	suite.Require().Error(err) // checking directory should throw an error

	// DirExist
	exist, err = DirExist(actualDir)
	suite.Require().NoError(err)
	suite.Require().True(exist)

	_, err = DirExist(absPath) // checking against file name should throw error
	suite.Require().Error(err)

	// MakeDir
	newPath := AbsDir(currentDir, "lvl_1/lvl_2/lvl_3")
	err = MakeDir(newPath)
	suite.Require().NoError(err)

	err = MakeDir(newPath) // creating the existing directory skips it.
	suite.Require().NoError(err)

	err = MakeDir(absPath) // it includes the file names
	suite.Require().Error(err)

	err = os.Remove(suite.envPath)
	suite.Require().NoError(err, "delete the dump file: "+suite.envPath)

	err = os.RemoveAll(newPath)
	suite.Require().NoError(err, "deleting dump files failed")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPath(t *testing.T) {
	suite.Run(t, new(TestPathSuite))
}
