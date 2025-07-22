# DeviceModel JSON 序列化/反序列化示例

本文档展示如何使用 DeviceModel 的 JSON 序列化和反序列化功能。

## 基本用法

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	edgeexpr "github.com/thinkontrol/edge-expr"
)

func main() {
	// 创建一个示例 DeviceModel
	publishCycle := 5 * time.Second
	cacheDuration := 2 * time.Minute
	diffThreshold := 0.5
	scale := 2.0
	offset := 1.0

	deviceModel := &edgeexpr.DeviceModel{
		Connections: map[string]string{
			"plc1": "modbus",
			"plc2": "ethernet",
		},
		Variables: map[string]*edgeexpr.Variable{
			"temperature": {
				Key:           "temperature",
				Connection:    "plc1",
				Address:       "DB1.DBD0",
				DataTypeStr:   "Float32",
				PublishCycle:  &publishCycle,
				CacheDuration: &cacheDuration,
				DiffThreshold: &diffThreshold,
				Scale:         &scale,
				Offset:        &offset,
				Writable:      false,
			},
			"status": {
				Key:         "status",
				Connection:  "plc1",
				Address:     "DB1.DBX1.0",
				DataTypeStr: "Bool",
				Writable:    true,
			},
			"calculated_value": {
				Key:         "calculated_value",
				Script:      "10 + 5",
				DataTypeStr: "Float64",
			},
		},
	}

	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(deviceModel, "", "  ")
	if err != nil {
		log.Fatalf("序列化失败: %v", err)
	}

	fmt.Println("序列化结果:")
	fmt.Println(string(jsonData))

	// 反序列化
	var unmarshaledModel edgeexpr.DeviceModel
	err = json.Unmarshal(jsonData, &unmarshaledModel)
	if err != nil {
		log.Fatalf("反序列化失败: %v", err)
	}

	fmt.Println("反序列化成功!")
	fmt.Printf("Hash: %s\n", unmarshaledModel.Hash())
}
```

## 测试功能

运行测试以验证 JSON 序列化和反序列化功能：

```bash
go test -v -run TestDeviceModel
```

测试涵盖了以下场景：

1. **基本序列化/反序列化**：测试 DeviceModel 的完整往返序列化
2. **错误处理**：测试无效的变量键、键不匹配、无效脚本等错误情况
3. **复杂场景**：测试空模型、nil 映射处理、所有可选字段、多种变量类型等
4. **数据完整性**：验证 Hash 一致性和序列化顺序保持

## JSON 格式示例

```json
{
  "connections": {
    "plc1": "modbus",
    "plc2": "ethernet"
  },
  "variables": {
    "temperature": {
      "key": "temperature",
      "connection": "plc1",
      "address": "DB1.DBD0",
      "data_type": "Float32",
      "diff_threshold": 0.5,
      "pct_threshold": 10,
      "scale": 2,
      "offset": 1,
      "publish_cycle": "5s",
      "cache_duration": "2m0s"
    },
    "status": {
      "key": "status", 
      "connection": "plc1",
      "address": "DB1.DBX1.0",
      "data_type": "Bool",
      "writable": true
    },
    "calculated_value": {
      "key": "calculated_value",
      "script": "10 + 5",
      "data_type": "Float64"
    }
  }
}
```

## 特性

- ✅ 完整的 JSON 序列化/反序列化支持
- ✅ 自动脚本编译（用于计算变量）
- ✅ 变量键验证
- ✅ 时间持续时间的字符串格式化
- ✅ nil 映射的自动初始化
- ✅ 错误处理和验证
- ✅ Hash 一致性保证
- ✅ 全面的测试覆盖
