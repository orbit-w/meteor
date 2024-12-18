### Plan
1. Create a `README.md` file for the `gen_proto` package.
2. Include sections for the project description, installation, usage, and examples.
3. Provide detailed instructions and code snippets where necessary.

## Usage

### GenerateProtoFromBaseExcel

The `GenerateProtoFromBaseExcel` function reads an Excel file and generates `proto3` message definitions.

#### Parameters

- `filename`: The path to the Excel file.
- `dst`: The name of the generated `.proto` file.
- `outputDir`: The output directory.

#### Example

```go
package main

import (
    "github.com/yourusername/gen_proto"
)

func main() {
    filename := "./src/Base-定义表.xlsx"
    gen_proto.GenerateProtoFromBaseExcel(filename, "BaseStructs.proto", "./proto")
}
```

### writeProtoToFile

The `writeProtoToFile` function writes the `proto3` message definitions to a file.

#### Parameters

- `filename`: The name of the generated `.proto` file.
- `packageName`: The package name for the proto file.
- `protoDefs`: A slice of `proto3` message definitions.
- `outputDir`: The output directory.

## Example

```go
package main

import (
    "github.com/yourusername/gen_proto"
)

func main() {
    protoDefs := []string{
        `message Example {
            int32 id = 1;
            string name = 2;
        }`,
    }
    gen_proto.writeProtoToFile("Example.proto", "example", protoDefs, "./proto")
}
```
```