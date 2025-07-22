package edgeexpr

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestDeviceModel_JSONSerialization(t *testing.T) {
	// 创建测试用的变量
	publishCycle := 5 * time.Second
	cacheDuration := 2 * time.Minute
	diffThreshold := 0.5
	pctThreshold := 10.0
	scale := 2.0
	offset := 1.0

	var1 := &Variable{
		Key:           "temperature",
		Connection:    "plc1",
		Address:       "DB1.DBD0",
		DataTypeStr:   "Float32",
		PublishCycle:  &publishCycle,
		CacheDuration: &cacheDuration,
		DiffThreshold: &diffThreshold,
		PctThreshold:  &pctThreshold,
		Scale:         &scale,
		Offset:        &offset,
		Writable:      false,
	}

	var2 := &Variable{
		Key:         "status",
		Connection:  "plc1",
		Address:     "DB1.DBX1.0",
		DataTypeStr: "Bool",
		Writable:    true,
	}

	var3 := &Variable{
		Key:         "calculated_value",
		Script:      "10 + 5",
		DataTypeStr: "Float64",
	}

	// 创建DeviceModel
	deviceModel := &DeviceModel{
		Connections: map[string]string{
			"plc1": "modbus",
			"plc2": "ethernet",
		},
		Variables: map[string]*Variable{
			"temperature":      var1,
			"status":           var2,
			"calculated_value": var3,
		},
	}

	// 测试序列化
	t.Run("Marshal", func(t *testing.T) {
		jsonData, err := json.MarshalIndent(deviceModel, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal DeviceModel: %v", err)
		}

		t.Logf("Serialized JSON:\n%s", string(jsonData))

		// 验证JSON包含期望的字段
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Fatalf("Failed to unmarshal serialized JSON: %v", err)
		}

		// 检查connections字段
		if connections, ok := result["connections"].(map[string]interface{}); ok {
			if connections["plc1"] != "modbus" {
				t.Errorf("Expected plc1 connection to be 'modbus', got %v", connections["plc1"])
			}
			if connections["plc2"] != "ethernet" {
				t.Errorf("Expected plc2 connection to be 'ethernet', got %v", connections["plc2"])
			}
		} else {
			t.Error("connections field not found or invalid type")
		}

		// 检查variables字段
		if variables, ok := result["variables"].(map[string]interface{}); ok {
			if len(variables) != 3 {
				t.Errorf("Expected 3 variables, got %d", len(variables))
			}

			// 检查temperature变量
			if tempVar, ok := variables["temperature"].(map[string]interface{}); ok {
				if tempVar["key"] != "temperature" {
					t.Errorf("Expected temperature key, got %v", tempVar["key"])
				}
				if tempVar["connection"] != "plc1" {
					t.Errorf("Expected plc1 connection, got %v", tempVar["connection"])
				}
				if tempVar["data_type"] != "Float32" {
					t.Errorf("Expected Float32 data type, got %v", tempVar["data_type"])
				}
			} else {
				t.Error("temperature variable not found")
			}
		} else {
			t.Error("variables field not found or invalid type")
		}
	})

	// 测试反序列化
	t.Run("Unmarshal", func(t *testing.T) {
		// 首先序列化
		jsonData, err := json.Marshal(deviceModel)
		if err != nil {
			t.Fatalf("Failed to marshal DeviceModel: %v", err)
		}

		// 然后反序列化
		var unmarshaledModel DeviceModel
		err = json.Unmarshal(jsonData, &unmarshaledModel)
		if err != nil {
			t.Fatalf("Failed to unmarshal DeviceModel: %v", err)
		}

		// 验证反序列化结果
		if len(unmarshaledModel.Connections) != 2 {
			t.Errorf("Expected 2 connections, got %d", len(unmarshaledModel.Connections))
		}

		if unmarshaledModel.Connections["plc1"] != "modbus" {
			t.Errorf("Expected plc1 connection to be 'modbus', got %v", unmarshaledModel.Connections["plc1"])
		}

		if len(unmarshaledModel.Variables) != 3 {
			t.Errorf("Expected 3 variables, got %d", len(unmarshaledModel.Variables))
		}

		// 验证temperature变量
		tempVar := unmarshaledModel.Variables["temperature"]
		if tempVar == nil {
			t.Fatal("temperature variable not found")
		}

		if tempVar.Key != "temperature" {
			t.Errorf("Expected temperature key, got %v", tempVar.Key)
		}

		if tempVar.Connection != "plc1" {
			t.Errorf("Expected plc1 connection, got %v", tempVar.Connection)
		}

		if tempVar.DataTypeStr != "Float32" {
			t.Errorf("Expected Float32 data type, got %v", tempVar.DataTypeStr)
		}

		if tempVar.DiffThreshold == nil || *tempVar.DiffThreshold != 0.5 {
			t.Errorf("Expected diff threshold 0.5, got %v", tempVar.DiffThreshold)
		}

		if tempVar.PublishCycle == nil || *tempVar.PublishCycle != 5*time.Second {
			t.Errorf("Expected publish cycle 5s, got %v", tempVar.PublishCycle)
		}

		// 验证calculated_value变量（脚本变量）
		calcVar := unmarshaledModel.Variables["calculated_value"]
		if calcVar == nil {
			t.Fatal("calculated_value variable not found")
		}

		if calcVar.Script != "10 + 5" {
			t.Errorf("Expected script '10 + 5', got %v", calcVar.Script)
		}

		// 验证脚本已编译
		if calcVar.Program == nil {
			t.Error("Expected script to be compiled, but Program is nil")
		}
	})

	// 测试往返序列化（roundtrip）
	t.Run("Roundtrip", func(t *testing.T) {
		// 序列化
		jsonData, err := json.Marshal(deviceModel)
		if err != nil {
			t.Fatalf("Failed to marshal DeviceModel: %v", err)
		}

		// 反序列化
		var unmarshaledModel DeviceModel
		err = json.Unmarshal(jsonData, &unmarshaledModel)
		if err != nil {
			t.Fatalf("Failed to unmarshal DeviceModel: %v", err)
		}

		// 再次序列化
		jsonData2, err := json.Marshal(&unmarshaledModel)
		if err != nil {
			t.Fatalf("Failed to marshal unmarshaled DeviceModel: %v", err)
		}

		// 比较两次序列化的结果应该相同（除了可能的字段顺序）
		var original, roundtrip map[string]interface{}
		if err := json.Unmarshal(jsonData, &original); err != nil {
			t.Fatalf("Failed to unmarshal original JSON: %v", err)
		}
		if err := json.Unmarshal(jsonData2, &roundtrip); err != nil {
			t.Fatalf("Failed to unmarshal roundtrip JSON: %v", err)
		}

		// 验证关键字段
		if len(original["connections"].(map[string]interface{})) != len(roundtrip["connections"].(map[string]interface{})) {
			t.Error("Connection counts don't match after roundtrip")
		}

		if len(original["variables"].(map[string]interface{})) != len(roundtrip["variables"].(map[string]interface{})) {
			t.Error("Variable counts don't match after roundtrip")
		}
	})
}

func TestDeviceModel_UnmarshalJSON_ErrorHandling(t *testing.T) {
	t.Run("InvalidVariableKey", func(t *testing.T) {
		jsonStr := `{
			"connections": {"plc1": "modbus"},
			"variables": {
				"invalid-key": {
					"key": "invalid-key",
					"connection": "plc1",
					"address": "DB1.DBD0",
					"data_type": "Float32"
				}
			}
		}`

		var deviceModel DeviceModel
		err := json.Unmarshal([]byte(jsonStr), &deviceModel)
		if err == nil {
			t.Error("Expected error for invalid variable key, but got none")
		}

		if !contains(err.Error(), "Invalid variable key") {
			t.Errorf("Expected error message to contain 'Invalid variable key', got: %v", err)
		}
	})

	t.Run("KeyMismatch", func(t *testing.T) {
		jsonStr := `{
			"connections": {"plc1": "modbus"},
			"variables": {
				"temperature": {
					"key": "different_key",
					"connection": "plc1",
					"address": "DB1.DBD0",
					"data_type": "Float32"
				}
			}
		}`

		var deviceModel DeviceModel
		err := json.Unmarshal([]byte(jsonStr), &deviceModel)
		if err == nil {
			t.Error("Expected error for key mismatch, but got none")
		}

		if !contains(err.Error(), "key mismatch") {
			t.Errorf("Expected error message to contain 'key mismatch', got: %v", err)
		}
	})

	t.Run("InvalidScript", func(t *testing.T) {
		jsonStr := `{
			"connections": {},
			"variables": {
				"calculated": {
					"key": "calculated",
					"script": "invalid syntax +++",
					"data_type": "Float64"
				}
			}
		}`

		var deviceModel DeviceModel
		err := json.Unmarshal([]byte(jsonStr), &deviceModel)
		if err == nil {
			t.Error("Expected error for invalid script, but got none")
		}

		if !contains(err.Error(), "calculated") {
			t.Errorf("Expected error message to contain variable name 'calculated', got: %v", err)
		}
	})

	t.Run("ValidJSON", func(t *testing.T) {
		jsonStr := `{
			"connections": {"plc1": "modbus"},
			"variables": {
				"temperature": {
					"key": "temperature",
					"connection": "plc1",
					"address": "DB1.DBD0",
					"data_type": "Float32",
					"publish_cycle": "5s",
					"cache_duration": "1m",
					"diff_threshold": 0.1,
					"pct_threshold": 5.0,
					"scale": 1.5,
					"offset": 2.0,
					"writable": true
				}
			}
		}`

		var deviceModel DeviceModel
		err := json.Unmarshal([]byte(jsonStr), &deviceModel)
		if err != nil {
			t.Errorf("Unexpected error for valid JSON: %v", err)
		}

		// 验证变量正确解析
		tempVar := deviceModel.Variables["temperature"]
		if tempVar == nil {
			t.Fatal("temperature variable not found")
		}

		if tempVar.Key != "temperature" {
			t.Errorf("Expected key 'temperature', got %v", tempVar.Key)
		}

		if tempVar.PublishCycle == nil || *tempVar.PublishCycle != 5*time.Second {
			t.Errorf("Expected publish cycle 5s, got %v", tempVar.PublishCycle)
		}

		if tempVar.CacheDuration == nil || *tempVar.CacheDuration != time.Minute {
			t.Errorf("Expected cache duration 1m, got %v", tempVar.CacheDuration)
		}
	})
}

func TestDeviceModel_Hash(t *testing.T) {
	// 创建两个相同的DeviceModel
	createDeviceModel := func() *DeviceModel {
		return &DeviceModel{
			Connections: map[string]string{
				"plc1": "modbus",
				"plc2": "ethernet",
			},
			Variables: map[string]*Variable{
				"temp": {
					Key:         "temp",
					Connection:  "plc1",
					Address:     "DB1.DBD0",
					DataTypeStr: "Float32",
				},
			},
		}
	}

	model1 := createDeviceModel()
	model2 := createDeviceModel()

	hash1 := model1.Hash()
	hash2 := model2.Hash()

	if hash1 != hash2 {
		t.Errorf("Expected same hash for identical models, got %s != %s", hash1, hash2)
	}

	// 修改一个模型，确保hash不同
	model2.Connections["plc3"] = "tcp"
	hash3 := model2.Hash()

	if hash1 == hash3 {
		t.Errorf("Expected different hash for different models, but got same hash: %s", hash1)
	}
}

func TestDeviceModel_ComplexSerialization(t *testing.T) {
	// 测试更复杂的JSON序列化/反序列化场景
	t.Run("EmptyDeviceModel", func(t *testing.T) {
		deviceModel := &DeviceModel{}

		// 序列化
		jsonData, err := json.Marshal(deviceModel)
		if err != nil {
			t.Fatalf("Failed to marshal empty DeviceModel: %v", err)
		}

		// 反序列化
		var unmarshaledModel DeviceModel
		err = json.Unmarshal(jsonData, &unmarshaledModel)
		if err != nil {
			t.Fatalf("Failed to unmarshal empty DeviceModel: %v", err)
		}

		// 验证初始化后的maps不为nil
		if unmarshaledModel.Connections == nil {
			t.Error("Expected connections map to be initialized")
		}
		if unmarshaledModel.Variables == nil {
			t.Error("Expected variables map to be initialized")
		}
	})

	t.Run("NilMapsHandling", func(t *testing.T) {
		jsonStr := `{"connections": null, "variables": null}`

		var deviceModel DeviceModel
		err := json.Unmarshal([]byte(jsonStr), &deviceModel)
		if err != nil {
			t.Fatalf("Failed to unmarshal DeviceModel with null maps: %v", err)
		}

		// 验证nil maps被正确初始化
		if deviceModel.Connections == nil {
			t.Error("Expected connections map to be initialized after unmarshaling null")
		}
		if deviceModel.Variables == nil {
			t.Error("Expected variables map to be initialized after unmarshaling null")
		}
	})

	t.Run("VariableWithAllOptionalFields", func(t *testing.T) {
		cacheDuration := 30 * time.Second
		publishCycle := 2 * time.Second
		diffThreshold := 1.5
		pctThreshold := 25.0
		scale := 0.5
		offset := -10.0

		jsonStr := fmt.Sprintf(`{
			"connections": {"modbus1": "tcp"},
			"variables": {
				"pressure": {
					"key": "pressure",
					"connection": "modbus1",
					"address": "40001",
					"data_type": "Float32",
					"diff_threshold": %f,
					"pct_threshold": %f,
					"scale": %f,
					"offset": %f,
					"writable": true,
					"publish_cycle": "%s",
					"cache_duration": "%s"
				}
			}
		}`, diffThreshold, pctThreshold, scale, offset, publishCycle.String(), cacheDuration.String())

		var deviceModel DeviceModel
		err := json.Unmarshal([]byte(jsonStr), &deviceModel)
		if err != nil {
			t.Fatalf("Failed to unmarshal complex DeviceModel: %v", err)
		}

		// 验证变量属性
		pressureVar := deviceModel.Variables["pressure"]
		if pressureVar == nil {
			t.Fatal("pressure variable not found")
		}

		if pressureVar.DiffThreshold == nil || *pressureVar.DiffThreshold != diffThreshold {
			t.Errorf("Expected diff threshold %f, got %v", diffThreshold, pressureVar.DiffThreshold)
		}

		if pressureVar.PctThreshold == nil || *pressureVar.PctThreshold != pctThreshold {
			t.Errorf("Expected pct threshold %f, got %v", pctThreshold, pressureVar.PctThreshold)
		}

		if pressureVar.Scale == nil || *pressureVar.Scale != scale {
			t.Errorf("Expected scale %f, got %v", scale, pressureVar.Scale)
		}

		if pressureVar.Offset == nil || *pressureVar.Offset != offset {
			t.Errorf("Expected offset %f, got %v", offset, pressureVar.Offset)
		}

		if !pressureVar.Writable {
			t.Error("Expected variable to be writable")
		}

		if pressureVar.PublishCycle == nil || *pressureVar.PublishCycle != publishCycle {
			t.Errorf("Expected publish cycle %v, got %v", publishCycle, pressureVar.PublishCycle)
		}

		if pressureVar.CacheDuration == nil || *pressureVar.CacheDuration != cacheDuration {
			t.Errorf("Expected cache duration %v, got %v", cacheDuration, pressureVar.CacheDuration)
		}
	})

	t.Run("MultipleVariableTypes", func(t *testing.T) {
		jsonStr := `{
			"connections": {
				"plc1": "modbus",
				"plc2": "ethernet"
			},
			"variables": {
				"temperature": {
					"key": "temperature",
					"connection": "plc1",
					"address": "DB1.DBD0",
					"data_type": "Float32"
				},
				"running": {
					"key": "running",
					"connection": "plc1", 
					"address": "DB1.DBX0.0",
					"data_type": "Bool"
				},
				"counter": {
					"key": "counter",
					"connection": "plc2",
					"address": "DB2.DBW0",
					"data_type": "UInt16"
				},
				"message": {
					"key": "message",
					"connection": "plc2",
					"address": "DB2.DBB10",
					"data_type": "String"
				},
				"calculated": {
					"key": "calculated",
					"script": "15 * 3",
					"data_type": "Int32"
				}
			}
		}`

		var deviceModel DeviceModel
		err := json.Unmarshal([]byte(jsonStr), &deviceModel)
		if err != nil {
			t.Fatalf("Failed to unmarshal multi-type DeviceModel: %v", err)
		}

		// 验证各种数据类型
		expectedTypes := map[string]string{
			"temperature": "Float32",
			"running":     "Bool",
			"counter":     "UInt16",
			"message":     "String",
			"calculated":  "Int32",
		}

		for varName, expectedType := range expectedTypes {
			variable := deviceModel.Variables[varName]
			if variable == nil {
				t.Errorf("Variable %s not found", varName)
				continue
			}
			if variable.DataTypeStr != expectedType {
				t.Errorf("Expected %s data type %s, got %s", varName, expectedType, variable.DataTypeStr)
			}
		}

		// 验证计算变量的脚本已编译
		calcVar := deviceModel.Variables["calculated"]
		if calcVar.Program == nil {
			t.Error("Expected calculated variable script to be compiled")
		}
	})

	t.Run("SerializationPreservesOrder", func(t *testing.T) {
		// 创建一个包含多个变量的DeviceModel
		deviceModel := &DeviceModel{
			Connections: map[string]string{
				"conn1": "type1",
				"conn2": "type2",
				"conn3": "type3",
			},
			Variables: map[string]*Variable{
				"var1": {Key: "var1", DataTypeStr: "Float32"},
				"var2": {Key: "var2", DataTypeStr: "Bool"},
				"var3": {Key: "var3", DataTypeStr: "Int32"},
			},
		}

		// 多次序列化，结果应该一致
		jsonData1, err := json.Marshal(deviceModel)
		if err != nil {
			t.Fatalf("Failed to marshal DeviceModel: %v", err)
		}

		jsonData2, err := json.Marshal(deviceModel)
		if err != nil {
			t.Fatalf("Failed to marshal DeviceModel second time: %v", err)
		}

		if string(jsonData1) != string(jsonData2) {
			t.Error("Serialization results should be consistent")
		}
	})

	t.Run("HashConsistency", func(t *testing.T) {
		cacheDuration := 1 * time.Minute
		deviceModel := &DeviceModel{
			Connections: map[string]string{
				"plc1": "modbus",
			},
			Variables: map[string]*Variable{
				"temp": {
					Key:           "temp",
					Connection:    "plc1",
					Address:       "DB1.DBD0",
					DataTypeStr:   "Float32",
					CacheDuration: &cacheDuration, // 明确设置cache duration
				},
			},
		}

		// 序列化后反序列化
		jsonData, err := json.Marshal(deviceModel)
		if err != nil {
			t.Fatalf("Failed to marshal DeviceModel: %v", err)
		}

		var unmarshaledModel DeviceModel
		err = json.Unmarshal(jsonData, &unmarshaledModel)
		if err != nil {
			t.Fatalf("Failed to unmarshal DeviceModel: %v", err)
		}

		// hash应该一致
		originalHash := deviceModel.Hash()
		unmarshaledHash := unmarshaledModel.Hash()

		if originalHash != unmarshaledHash {
			t.Errorf("Hash should be consistent after round-trip serialization. Original: %s, Unmarshaled: %s", originalHash, unmarshaledHash)
		}
	})
}

// 辅助函数：检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
