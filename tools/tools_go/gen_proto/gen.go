package gen_proto

/*
   @Author: orbit-w
   @File: gen
   @2024 11月 周六 15:49
*/

import (
	"fmt"
	"github.com/tealeg/xlsx" // 假设使用tealeg/xlsx库来处理Excel文件
	"log"
	"os"
)

const (
	rowStructName = "名称"
	rowStructDef  = "结构"

	structPackage = "Structs"
)

// GenerateProtoFromBaseExcel 读取Base.xlsx文件并生成proto3 message定义
// GenerateProtoFromBaseExcel reads the Base.xlsx file and generates proto3 message definitions
// 参数 (Parameters):
// - filename: Excel文件的路径 (The path to the Excel file)
// - dst: 生成的.proto文件的名称 (The name of the generated .proto file)
// - outputDir: 输出目录 (The output directory)
func GenerateProtoFromBaseExcel(filename, dst, outputDir string) {
	file, err := xlsx.OpenFile(filename)
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}

	// 假设我们只处理第一个sheet
	sheet := file.Sheets[0]

	// 存储proto3 message定义
	var protoDefs []string

	var dict = map[string]bool{}

	// 收集所有的结构名称
	for _, row := range sheet.Rows {
		// 跳过标题行
		if row.Cells[0].String() == rowStructName {
			continue
		}
		name := row.Cells[0].String()
		dict[name] = true
	}

	gen := NewMessageGenerator(dict)

	// 遍历sheet中的每一行
	for _, row := range sheet.Rows {
		// 跳过标题行
		if row.Cells[0].String() == rowStructName {
			continue
		}

		// 获取名称和结构
		name := row.Cells[0].String()
		structure := row.Cells[1].String()

		protoDef := gen.GenProtocol3Message(name, structure)

		protoDefs = append(protoDefs, protoDef)
	}

	// 将protoDefs写入到proto文件中
	writeProtoToFile(dst, structPackage, protoDefs, outputDir)
}

// writeProtoToFile 将proto3 message定义写入到文件
// writeProtoToFile writes the proto3 message definitions to a file
// 参数 (Parameters):
// - filename: 生成的.proto文件的名称 (The name of the generated .proto file)
// - packageName: proto文件的包名 (The package name for the proto file)
// - protoDefs: proto3 message定义的切片 (A slice of proto3 message definitions)
// - outputDir: 输出目录 (The output directory)
func writeProtoToFile(filename, packageName string, protoDefs []string, outputDir string) {
	content := fmt.Sprintf("syntax = \"proto3\";\n\npackage %s;\n\n", packageName)
	for _, def := range protoDefs {
		content += def + "\n"
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("无法创建目录: %v", err)
	}

	filePath := fmt.Sprintf("%s/%s", outputDir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Fatalf("无法写入文件: %v", err)
	}
}
