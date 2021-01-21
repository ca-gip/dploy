package ansible

import "fmt"

const testPath = "./../../test"

var (
	projectMultiLevelPath  = fmt.Sprint(testPath, "/projectMultiLevel")
	projectSimpleLevelPath = fmt.Sprint(testPath, "/projectSimpleLevel")
)
