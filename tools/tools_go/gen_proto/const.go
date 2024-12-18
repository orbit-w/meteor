package gen_proto

/*
   @Author: orbit-w
   @File: const
   @2024 11月 周六 20:35
*/

const (
	fieldTypeInt32  = "int32"
	fieldTypeInt64  = "int64"
	fieldTypeUint32 = "uint32"
	fieldTypeUint64 = "uint64"
	fieldTypeFloat  = "float"
	fieldTypeDouble = "double"
	fieldTypeString = "string"
	fieldTypeBool   = "bool"
	fieldTypeLong   = "long"
	fieldTypeInt    = "int"
	fieldTypeBytes  = "bytes"
)

const (
	RowChineseName = iota
	RowDesc
	RowName
	RowType
	RowPermission
	RowKey
	RawDataStart
)
