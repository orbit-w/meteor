package gen_proto

import "testing"

/*
   @Author: orbit-w
   @File: gen_test
   @2024 11月 周六 15:49
*/

func Test_generateProtoFromExcel(t *testing.T) {
	filename := "./src/Base-定义表.xlsx"
	GenerateProtoFromBaseExcel(filename, "BaseStructs.proto", "./proto")
}
