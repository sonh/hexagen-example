package main

import "testing"

func TestUsrGen(t *testing.T) {

	typeName := "User"
	outputPath := "./testdata"
	outpkg := "testdata"

	genMock(typeName, outputPath, outpkg, []string{"../internal/usermgmt/modules/user/core/entity"}...)
}
