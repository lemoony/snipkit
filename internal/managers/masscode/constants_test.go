package masscode

import "path/filepath"

const (
	testDataUserHomeV1 = "testdata/userhome-v1"
	testDataUserHomeV2 = "testdata/userhome-v2"
)

var (
	testDataMassCodeV2Path = filepath.Join(testDataUserHomeV2, defaultMassCodeHomePath)
	testDataLibraryV2Path  = filepath.Join(testDataUserHomeV2, defaultMassCodeHomePath, v2DatabaseFile)
)
