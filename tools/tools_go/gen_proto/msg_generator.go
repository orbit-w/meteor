package gen_proto

import (
	"fmt"
	"log"
	"strings"
)

/*
   @Author: orbit-w
   @File: parser
   @2024 12月 周日 07:58
*/

type MessageGenerator struct {
	dict map[string]bool
}

func NewMessageGenerator(dict map[string]bool) *MessageGenerator {
	return &MessageGenerator{
		dict: dict,
	}
}

func (p *MessageGenerator) GenProtocol3Message(msgName, structure string) string {
	// 移除结构中的大括号
	structure = strings.Trim(structure, "{}")

	// 创建proto3 message定义
	protoDef := fmt.Sprintf("message %s {\n", msgName)
	// 根据结构添加成员变量
	fields := strings.Split(structure, ",")
	for i, field := range fields {
		fieldParts := strings.Split(field, " ")
		if len(fieldParts) != 2 {
			log.Fatalf("结构字段格式错误: %s", field)
		}
		fieldType := fieldParts[0]
		fieldName := fieldParts[1]

		fieldType = p.parseFieldStr(fieldType)
		protoDef += fmt.Sprintf("  %s %s = %d;\n", fieldType, fieldName, i+1)
	}
	protoDef += "}\n"
	return protoDef
}

func (p *MessageGenerator) parseFieldStr(str string) (fieldType string) {
	switch {
	case strings.HasPrefix(str, "[]"):
		//处理数组类型
		elementType := str[2:]
		fieldType = "repeated " + p.checkAndRepl(elementType)
	default:
		fieldType = "optional " + p.checkAndRepl(str)
	}

	return
}

func (p *MessageGenerator) checkAndRepl(str string) (res string) {
	switch str {
	case fieldTypeInt:
		res = fieldTypeInt32
	case fieldTypeLong:
		res = fieldTypeDouble
	case fieldTypeInt32, fieldTypeInt64, fieldTypeUint32, fieldTypeUint64, fieldTypeFloat, fieldTypeDouble, fieldTypeString, fieldTypeBool, fieldTypeBytes:
		res = str
	default:
		if p.dict != nil && p.dict[str] {
			res = str
			break
		}
		panic(fmt.Sprintf("未知的字段类型: %s", str))
	}
	return
}
