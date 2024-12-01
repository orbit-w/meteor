package gen_proto

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
)

/*
   @Author: orbit-w
   @File: load
   @2024 12月 周日 22:55
*/

func Load(filename string) {
	// 打开XLSX文件
	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// 获取所有工作表的名称
	sheets := f.GetSheetMap()
	for _, name := range sheets {
		// 读取指定工作表的数据
		rows, err := f.GetRows(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 假设第一行是列名，第二行是类型信息
		if len(rows) < RawDataStart {
			fmt.Println("数据行不足")
			return
		}

		columns := initColumns(rows)

		// 从第三行开始遍历数据行
		for i := RawDataStart; i < len(rows); i++ {
			row := rows[i]

			// 创建map来存储行数据
			rowData := fill(row, columns)

			// 打印map，这里可以根据需要处理map中的数据
			fmt.Println(rowData)
		}
	}
}

func fill(row []string, columns Columns) map[string]any {
	// 创建map来存储行数据
	rowData := make(map[string]any)

	// 遍历每一列，将列名和值存入map
	for j, colCell := range row {
		col := columns[j]

		// 根据类型信息转换单元格的值
		switch col.Type {
		case "int32":
			if value, err := strconv.ParseInt(colCell, 10, 32); err == nil {
				rowData[col.Name] = int32(value)
			}
		case "int64":
			if value, err := strconv.ParseInt(colCell, 10, 64); err == nil {
				rowData[col.Name] = value
			}
		case "float64":
			if value, err := strconv.ParseFloat(colCell, 64); err == nil {
				rowData[col.Name] = value
			}
		case "bool":
			if value, err := strconv.ParseBool(colCell); err == nil {
				rowData[col.Name] = value
			}
		case "string":
			rowData[col.Name] = colCell
		// 根据需要添加更多类型转换
		default:
			rowData[col.Name] = colCell // 默认作为字符串处理
		}
	}
	return rowData
}

func initColumns(rows [][]string) Columns {
	columns := make(Columns, 0)
	head := rows[0]
	for i := range head {
		chName := rows[RowChineseName][i]
		desc := rows[RowDesc][i]
		t := rows[RowType][i]
		permission := rows[RowPermission][i]
		key := rows[RowKey][i] == "key"
		columns = append(columns, newColumn(head[i], chName, desc, t, permission, key))
	}
	return columns
}
