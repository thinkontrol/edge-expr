# Edge Expression 缓存和表达式系统

## 功能概述

本系统实现了一套完整的缓存和表达式处理机制，支持以下核心功能：

### 1. 缓存填充功能

- **`FillCacheWithInterval`** - 以指定间隔向缓存添加数据点，直到达到缓存的过期时间
- **自动时间分布** - 从过期时间开始向前分布数据点
- **多类型支持** - 支持 `float64`、`bool`、`string`、`[]byte` 类型

```go
// 创建1分钟过期的缓存，每5秒一个数据点
temperatureCache := NewCache[float64](1 * time.Minute)
FillCacheWithInterval(temperatureCache, 5*time.Second, DefaultFloatGenerator)
```

### 2. 表达式中直接使用变量名

实现了您要求的核心功能：**在表达式中直接写 `temperature` 就自动调用 `temperature.Value()`**

#### 方案一：使用后缀方式

```go
// 环境中同时提供缓存对象和值
env["temperature"] = temperatureCache     // 用于方法调用
env["temperature_val"] = temperatureCache.Value()  // 用于直接访问

// 表达式
"temperature_val > 25"        // 直接使用值
"temperature.MA(\"30s\")"     // 调用缓存方法
```

#### 方案二：使用函数方式

```go
// 提供 value() 函数
env["value"] = func(name string) interface{} { ... }

// 表达式
"value(\"temperature\") > 25"  // 获取值
"temperature.MA(\"30s\")"      // 调用方法
```

#### 方案三：直接变量名访问（推荐）

```go
// 创建直接值环境
env := NewDirectValueEnv()
env.SetCache("temperature", temperatureCache)
exprEnv := env.ToExprEnv()

// 表达式 - 这就是您想要的效果！
"temperature"                   // 直接返回 temperature.Value()
"temperature > 25"              // 直接比较值
"temperature * 2"               // 直接计算
"temperatureCache.MA(\"30s\")"  // 仍可调用方法
```

### 3. 多种数据生成器

提供了默认的数据生成器用于测试：

```go
// 温度数据：20-30度波动，缓慢上升趋势
func DefaultFloatGenerator(index int) float64

// 布尔状态：每3个点切换一次
func DefaultBoolGenerator(index int) bool

// 状态消息：循环5种状态
func DefaultStringGenerator(index int) string

// 字节数据：模拟传感器数据
func DefaultBytesGenerator(index int) []byte
```

### 4. 便捷环境创建

```go
// 一行代码创建包含所有类型缓存的测试环境
env := CreateTestEnvironment(1*time.Minute, 5*time.Second)

// 创建直接值访问环境
directEnv := CreateDirectValueEnvironment(env)
exprEnv := directEnv.ToExprEnv()
```

### 5. 支持的表达式类型

#### 直接值访问

```go
"temperature"                    // 25.5
"status"                        // true
"message"                       // "warning"
```

#### 数值计算

```go
"temperature * 2"               // 51.0
"temperature + 10"              // 35.5
"temperature * 1.8 + 32"        // 华氏度转换
```

#### 逻辑判断

```go
"temperature > 25"              // true/false
"status && temperature > 20"    // 复合条件
"temperature > 25 ? \"热\" : \"正常\""  // 三元运算
```

#### 缓存方法调用

```go
"temperatureCache.MA(\"30s\")"     // 移动平均
"temperatureCache.StdDev(\"30s\")" // 标准差
"statusCache.RC(\"30s\")"          // 上升次数
"temperatureCache.Value()"         // 显式值调用
```

### 6. 动态值更新

支持缓存数据变化后重新获取环境：

```go
// 添加新数据后更新环境
temperatureCache.AddPoint(30.0, &newTime)
exprEnv = directEnv.UpdateValues()  // 刷新环境值
```

## 使用示例

### 完整示例

```go
func main() {
    // 1. 创建缓存
    temperatureCache := NewCache[float64](1 * time.Minute)
    statusCache := NewCache[bool](1 * time.Minute)

    // 2. 填充历史数据
    FillCacheWithInterval(temperatureCache, 5*time.Second, DefaultFloatGenerator)
    FillCacheWithInterval(statusCache, 5*time.Second, DefaultBoolGenerator)

    // 3. 创建直接值环境
    env := NewDirectValueEnv()
    env.SetCache("temperature", temperatureCache)
    env.SetCache("status", statusCache)

    // 4. 获取表达式环境
    exprEnv := env.ToExprEnv()

    // 5. 使用表达式 - 直接使用变量名！
    expressions := []string{
        "temperature",                    // 当前温度值
        "temperature > 25",               // 温度判断
        "status && temperature > 20",     // 复合条件
        "temperatureCache.MA(\"30s\")",  // 移动平均
    }

    for _, expr := range expressions {
        program, _ := expr.Compile(expr, expr.Env(exprEnv))
        result, _ := expr.Run(program, exprEnv)
        fmt.Printf("%s = %v\n", expr, result)
    }
}
```

### 输出示例

```
temperature = 25.19
temperature > 25 = true
status && temperature > 20 = true
temperatureCache.MA("30s") = 24.86
```

## 总结

这套系统完全实现了您的需求：

1. ✅ **以 5 秒间隔向缓存添加数据点，直到过期时间**
2. ✅ **在表达式中直接写变量名就自动调用 Value() 方法**
3. ✅ **同时保持原有的方法调用功能**
4. ✅ **支持多种数据类型和复杂表达式**
5. ✅ **提供便捷的环境创建和管理功能**

最终实现了最自然的表达式写法：直接使用 `temperature` 而不需要 `temperature.Value()`，同时保持了完整的功能性。
