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
	envPath    string // using env as a file to check
	currentDir string
	fileName   string
	dir        string
}

// Make sure that Account is set to five
// before each test
func (test *TestPathSuite) SetupTest() {
	test.fileName = ".test.env"
	test.dir = "configs"
	test.envPath = filepath.Join("configs", ".test.env")

	currentDir, err := CurrentDir()
	test.Require().NoError(err)
	test.currentDir = currentDir

	err = MakeDir(filepath.Join(currentDir, test.dir))
	test.Require().NoError(err)

	file, err := os.OpenFile(filepath.Join(currentDir, test.envPath), os.O_CREATE|os.O_RDWR, 0755)
	test.Require().NoError(err)

	_, err = file.WriteString("KEY=123")
	test.Require().NoError(err, "failed to write the data into: "+test.envPath)
	err = file.Close()
	test.Require().NoError(err, "close the file: "+test.envPath)
}

func (test *TestPathSuite) TearDownTest() {
	absPath := filepath.Join(test.currentDir, test.envPath)
	err := os.Remove(absPath)
	test.Require().NoError(err, "delete the dump file: "+absPath)

}

// All methods that begin with "Test" are run as tests within a
// suite.
func (test *TestPathSuite) TestRun() {
	currentDir, err := CurrentDir()
	test.Require().NoError(err)
	test.Require().True(filepath.IsAbs(currentDir))

	// AbsDir
	expected := filepath.Join(currentDir, test.envPath)
	absPath := AbsDir(currentDir, test.envPath)
	test.Require().Equal(expected, absPath)

	// FileName
	test.Require().Equal(test.fileName, FileName(absPath))

	// DirAndFileName
	actualDir, actualFileName := DirAndFileName(absPath)
	test.Require().Equal(filepath.Join(currentDir, test.dir), actualDir)
	test.Require().Equal(test.fileName, actualFileName)

	// FileExist
	fmt.Printf("abs path: %s\n", absPath)
	exist, err := FileExist(absPath)
	test.Require().NoError(err)
	test.Require().True(exist)

	_, err = FileExist(actualDir)
	test.Require().Error(err) // checking directory should throw an error

	// DirExist
	exist, err = DirExist(actualDir)
	test.Require().NoError(err)
	test.Require().True(exist)

	_, err = DirExist(absPath) // checking against file name should throw error
	test.Require().Error(err)

	// MakeDir
	newPath := AbsDir(currentDir, "lvl_1/lvl_2/lvl_3")
	err = MakeDir(newPath)
	test.Require().NoError(err)

	err = MakeDir(newPath) // creating the existing directory skips it.
	test.Require().NoError(err)

	err = MakeDir(absPath) // it includes the file names
	test.Require().Error(err)

	err = os.RemoveAll(newPath)
	test.Require().NoError(err, "deleting dump files failed")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPath(t *testing.T) {
	suite.Run(t, new(TestPathSuite))
}
