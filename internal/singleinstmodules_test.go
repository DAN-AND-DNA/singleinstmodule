package internal

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type A struct {
	Count int
}

func (a *A) ModuleConstruct()       { a.Count = 0 }
func (a *A) ModuleDestruct()        { a.Count = 0 }
func (a *A) ModuleBeforeRun()       { a.Count += 100 }
func (a *A) ModuleLock() ModuleCore { return nil }
func (a *A) ModuleUnlock()          {}
func (a *A) Run1()                  { a.Count += 1 }
func (a *A) Run2()                  { a.Count += 2 }
func (a *A) Run3()                  { a.Count += 3 }
func (a *A) AfterRun()              { a.Count -= 100 }
func (a *A) ModuleShutdown()        { a.Count -= 1 }

type B struct {
}

func (b *B) ModuleConstruct()       {}
func (b *B) ModuleDestruct()        {}
func (b *B) ModuleLock() ModuleCore { return nil }
func (b *B) ModuleUnlock()          {}

type C struct {
	B
}

type SingleInstModuleTestSuite struct {
	suite.Suite
	target *SingleInstModules
	a      *A
	b      *B
	c      *C
}

func (suite *SingleInstModuleTestSuite) SetupSuite() {

}

func (suite *SingleInstModuleTestSuite) TearDownSuite() {

}

func (suite *SingleInstModuleTestSuite) SetupTest() {
	suite.target = new(SingleInstModules)
	suite.target.Construct()
	suite.a = new(A)
	suite.b = new(B)
	suite.c = new(C)
}

func (suite *SingleInstModuleTestSuite) TearDownTest() {
	suite.a = nil
	suite.b = nil
	suite.c = nil
	//suite.target.Destruct()
	suite.target = nil
}

func (suite *SingleInstModuleTestSuite) Test_1_RegisterModule() {
	suite.Equal(true, suite.target.Register(suite.a))
	suite.Equal(false, suite.target.Register(suite.a))
	suite.Equal(true, suite.target.Register(suite.b))
	suite.Equal(false, suite.target.Register(suite.b))
	suite.Equal(true, suite.target.Register(suite.c))
	suite.Equal(false, suite.target.Register(suite.c))
}

func (suite *SingleInstModuleTestSuite) Test_2_GetModuleName() {
	suite.Equal("A", suite.target.GetModuleName(suite.a))
	suite.Equal("B", suite.target.GetModuleName(suite.b))
	suite.Equal("C", suite.target.GetModuleName(suite.c))
}

func (suite *SingleInstModuleTestSuite) Test_3_Run() {
	suite.target.Register(suite.a)
	suite.target.Run(true)
	suite.Equal(39, suite.a.Count)
	suite.target.Shutdown()
	suite.Equal(0, suite.a.Count)
}

func TestSingleInstModule(t *testing.T) {
	suite.Run(t, new(SingleInstModuleTestSuite))
}
