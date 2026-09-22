# Cache expr 脚本方法手册

本文档用于指导 expr 脚本编写者通过变量缓存查询当前值、历史统计和状态变化。

## 基本写法

设备模型中的变量会以变量名放入 expr 环境。假设模型中有以下变量：

- `temperature`：数值型温度
- `running`：布尔型运行状态
- `status`：`uint32`、`uint64` 或字节数组状态值

可在脚本中直接调用它们的 Cache 方法：

```expr
temperature.Value() > 80
temperature.MA("5m") > 60
running.Rising()
status.Bit(3)
```

方法返回值可继续参与 expr 的比较、算术和逻辑运算：

```expr
temperature.MA("5m") > 60 && temperature.Changed()
running.Rising() || running.Falling()
status.BitAnd(15) == 5
```

> 以下示例中的 `temperature`、`running` 和 `status` 仅为示例变量名。实际脚本应替换为设备模型中定义的变量名。

## 时间窗口

带 `window` 参数的方法用于查询最近一段时间内的数据。窗口必须写成字符串：

```expr
temperature.MA("500ms")
temperature.MA("10s")
temperature.MA("5m")
temperature.MA("2h30m")
```

常用单位：

| 单位 | 含义 | 示例 |
| --- | --- | --- |
| `ms` | 毫秒 | `"500ms"` |
| `s` | 秒 | `"30s"` |
| `m` | 分钟 | `"5m"` |
| `h` | 小时 | `"2h"` |

窗口以脚本执行时刻为终点。例如 `"5m"` 表示执行时刻之前的 5 分钟。

> `MA`、`StdDev`、`Count`、`RC`、`FC`、`Contains`、`EveryBy`、`SomeBy` 和 `NoneBy` 收到无效窗口时，会使用缓存中的全部数据，而不是报告窗口格式错误。`PctChangeSince` 和 `DiffSince` 收到无效窗口时会报错。因此脚本中应始终使用有效窗口。

## 方法速查

| 适用类型 | 方法 | 返回值 | 说明 |
| --- | --- | --- | --- |
| 所有类型 | `Value()` | 当前值类型 | 最新值 |
| 所有类型 | `Timestamp()` | 时间 | 最新值的时间戳 |
| 所有类型 | `Len()` | 整数 | 缓存中的数据点数量 |
| 所有类型 | `Count(window)` | 整数 | 窗口内连续值段数量 |
| 所有类型 | `Contains(value, window)` | 布尔 | 窗口内是否出现过指定值 |
| 所有类型 | `Changed()` | 布尔 | 最新值是否相对前一个值发生变化 |
| 所有类型 | `EveryBy(condition, window)` | 布尔 | 窗口内是否每个值都满足条件 |
| 所有类型 | `SomeBy(condition, window)` | 布尔 | 窗口内是否至少一个值满足条件 |
| 所有类型 | `NoneBy(condition, window)` | 布尔 | 窗口内是否没有值满足条件 |
| 数值型 | `MA(window)` | 数值 | 窗口内平均值 |
| 数值型 | `StdDev(window)` | 数值 | 窗口内总体标准差 |
| 数值型 | `PctChange()` | 数值 | 最新两个值的百分比变化 |
| 数值型 | `PctChangeWith(value)` | 数值 | 最新值相对指定值的百分比变化 |
| 数值型 | `PctChangeSince(window)` | 数值 | 最新值相对指定时间前的百分比变化 |
| 数值型 | `PctChangeExceeds(threshold)` | 布尔 | 百分比变化绝对值是否超过阈值 |
| 数值型 | `Diff()` | 数值 | 最新值减前一个值 |
| 数值型 | `DiffWith(value)` | 数值 | 最新值减指定值 |
| 数值型 | `DiffSince(window)` | 数值 | 最新值减指定时间前的值 |
| 数值型 | `DiffExceeds(threshold)` | 布尔 | 差值绝对值是否超过阈值 |
| 布尔型 | `Rising()` | 布尔 | 最新两个值是否为 `false -> true` |
| 布尔型 | `Falling()` | 布尔 | 最新两个值是否为 `true -> false` |
| 布尔型 | `RC(window)` | 整数 | 窗口内上升沿次数 |
| 布尔型 | `FC(window)` | 整数 | 窗口内下降沿次数 |
| 位值类型 | `Bit(index)` | 布尔 | 指定位是否为 1 |
| 位值类型 | `ByteBit(byte, bit)` | 布尔 | 指定字节的指定位是否为 1 |
| 位值类型 | `BitAnd(mask)` | 整数 | 与掩码按位与 |
| 位值类型 | `BitOr(mask)` | 整数 | 与掩码按位或 |
| 位值类型 | `BitXor(mask)` | 整数 | 与掩码按位异或 |
| 位值类型 | `BitClear(mask)` | 整数 | 清除掩码指定的位 |
| 位值类型 | `BitNot()` | 整数 | 按位取反 |

## 当前值和缓存状态

### `Value()`

读取变量的最新值。

```expr
temperature.Value() > 80
running.Value() == true
```

返回值类型与变量类型一致。没有缓存数据时返回 `nil`。

### `Timestamp()`

读取最新值的时间戳。

```expr
temperature.Timestamp()
```

没有缓存数据时返回 `nil`。

### `Len()`

返回缓存中的原始数据点数量。

```expr
temperature.Len() >= 10
```

适合在统计计算前确认数据量是否充足：

```expr
temperature.Len() >= 2 && temperature.StdDev("5m") > 3
```

### `Count(window)`

返回窗口内连续值段的数量，不是原始数据点数量。相邻重复值只计数一次。

例如历史值为 `[1, 1, 2, 2, 1]` 时：

```expr
variable.Count("5m") == 3
```

可用于判断状态在窗口内是否发生过多次变化：

```expr
mode.Count("10m") > 5
```

### `Contains(value, window)`

判断窗口内是否出现过指定值，适用于数值、布尔、字符串等类型。

```expr
temperature.Contains(100, "5m")
running.Contains(false, "1h")
state.Contains("alarm", "10m")
```

窗口内没有数据时返回 `false`。

### `Changed()`

判断最新值是否与前一个值不同。

```expr
temperature.Changed()
state.Changed() && state.Value() == "alarm"
```

少于两个数据点时返回 `false`。

## 数值统计

本节方法用于数值型变量。

### `MA(window)`

计算窗口内数据的算术平均值。

```expr
temperature.MA("5m") > 80
pressure.Value() > pressure.MA("10m") * 1.2
```

窗口内没有数据或值不是数值类型时，脚本执行会报错。

### `StdDev(window)`

计算窗口内数据的总体标准差，用于判断数值波动程度。

```expr
temperature.StdDev("5m") > 3
```

至少需要两个数据点。该方法计算总体标准差，分母为数据点数量 `N`。

## 数值变化

### `PctChange()`

计算最新值相对前一个值的百分比变化：

```text
((最新值 - 前一个值) / 前一个值) * 100
```

```expr
temperature.PctChange() > 10
abs(temperature.PctChange()) > 10
```

至少需要两个数据点。如果前一个值为 `0` 且最新值不为 `0`，脚本执行会报错。

### `PctChangeWith(value)`

计算最新值相对指定基准值的百分比变化：

```text
((最新值 - 基准值) / 基准值) * 100
```

```expr
temperature.PctChangeWith(80) > 20
```

至少需要一个数据点。基准值为 `0` 且最新值不为 `0` 时会报错。

### `PctChangeSince(window)`

计算最新值相对指定时间之前最近一个数据点的百分比变化。

```expr
temperature.PctChangeSince("5m") > 10
```

例如当前时间为 `10:10`，窗口为 `"5m"`，方法会查找 `10:05` 之前最接近 `10:05` 的数据点作为基准。

如果指定时间之前没有数据点，脚本执行会报错。

### `PctChangeExceeds(threshold)`

判断最新两个值的百分比变化绝对值是否严格大于阈值。

```expr
temperature.PctChangeExceeds(10)
```

以下两种变化都会返回 `true`：

```text
+12% > 10%
-12% 的绝对值 > 10%
```

### `Diff()`

计算最新值减去前一个值。

```expr
temperature.Diff() > 5
temperature.Diff() < -5
abs(temperature.Diff()) > 5
```

至少需要两个数据点。

### `DiffWith(value)`

计算最新值减去指定值，至少需要一个数据点。

```expr
temperature.DiffWith(80) > 5
```

当最新值为 `90` 时，`temperature.DiffWith(80)` 返回 `10`。

### `DiffSince(window)`

计算最新值减去指定时间之前最近一个数据点的值。

```expr
temperature.DiffSince("5m") > 10
```

数据点选择规则与 `PctChangeSince(window)` 相同。如果指定时间之前没有数据点，脚本执行会报错。

### `DiffExceeds(threshold)`

判断最新两个值的差值绝对值是否严格大于阈值。

```expr
temperature.DiffExceeds(5)
```

无论数值上升超过 `5` 还是下降超过 `5`，都会返回 `true`。

## 布尔值和边沿

本节方法用于布尔型变量。

### `Rising()`

判断最新两个值是否从 `false` 变为 `true`。

```expr
running.Rising()
```

常用于设备启动事件：

```expr
motor_running.Rising() && temperature.Value() < 80
```

少于两个数据点时返回 `false`。

### `Falling()`

判断最新两个值是否从 `true` 变为 `false`。

```expr
running.Falling()
```

常用于设备停止事件：

```expr
motor_running.Falling() && alarm.Value() == true
```

少于两个数据点时返回 `false`。

### `RC(window)`

统计窗口内 `false -> true` 的次数，即上升沿次数。

```expr
running.RC("10m") > 5
```

可用于检测设备频繁启动。

### `FC(window)`

统计窗口内 `true -> false` 的次数，即下降沿次数。

```expr
running.FC("10m") > 5
```

可用于检测设备频繁停止。

## 位操作

本节方法用于 `uint32`、`uint64` 或字节数组变量。字节数组按小端序解释，即第 0 个字节是最低有效字节。

位和字节索引都从 `0` 开始。

### `Bit(index)`

判断最新值的指定位是否为 `1`。

```expr
status.Bit(0)
status.Bit(7) && running.Value()
```

`uint32` 的有效索引为 `0-31`，`uint64` 的有效索引为 `0-63`。字节数组的有效索引取决于数组长度。

### `ByteBit(byteIndex, bitIndex)`

判断第 `byteIndex` 个字节的第 `bitIndex` 位是否为 `1`。`bitIndex` 的范围为 `0-7`。

```expr
status.ByteBit(0, 3)
status.ByteBit(1, 7)
```

例如 `ByteBit(1, 2)` 等价于读取全局位索引 `10`。

### `BitAnd(mask)`

将最新值与掩码进行按位与，常用于判断一组状态位。

```expr
status.BitAnd(15) == 5
status.BitAnd(8) != 0
```

### `BitOr(mask)`

将最新值与掩码进行按位或。

```expr
status.BitOr(1) == 5
```

### `BitXor(mask)`

将最新值与掩码进行按位异或。

```expr
status.BitXor(3) == 6
```

### `BitClear(mask)`

清除掩码中为 `1` 的位，返回清除后的整数值。

```expr
status.BitClear(8) == 5
```

### `BitNot()`

对最新值按位取反。

```expr
status.BitNot()
```

`uint32` 按 32 位取反，`uint64` 按 64 位取反。字节数组转换为最多 64 位的小端序整数后按 64 位取反。

## 条件查询

条件查询用于判断一个时间窗口内的历史值是否满足给定条件。

条件参数本身是一段 expr 表达式字符串，其中 `x` 代表正在检查的历史值：

```expr
temperature.SomeBy("x > 80", "5m")
```

注意条件必须放在字符串中。正确与错误写法对比：

```expr
// 正确
temperature.SomeBy("x > 80", "5m")

// 错误：条件没有写成字符串
temperature.SomeBy(x > 80, "5m")

// 错误：不使用 JavaScript 箭头函数语法
temperature.SomeBy("(x) => x > 80", "5m")
```

条件必须对每个历史值都返回布尔值。窗口内的数据类型也应保持一致。

### `EveryBy(condition, window)`

窗口内每个值都满足条件时返回 `true`。

```expr
temperature.EveryBy("x >= 0 && x <= 100", "5m")
running.EveryBy("x == true", "30s")
```

窗口内没有数据时返回 `false`。

> 当前实现还要求窗口内的数据点数量少于缓存总数据点数量。如果窗口覆盖缓存中的全部数据，方法返回 `false`。这用于确认指定窗口之前仍有缓存数据。

### `SomeBy(condition, window)`

窗口内至少一个值满足条件时返回 `true`。

```expr
temperature.SomeBy("x > 80", "5m")
state.SomeBy("x == 'alarm'", "10m")
```

窗口内没有数据时返回 `false`。

### `NoneBy(condition, window)`

窗口内没有任何值满足条件时返回 `true`。

```expr
temperature.NoneBy("x > 100", "5m")
running.NoneBy("x == false", "30s")
```

窗口内没有数据时返回 `false`。

> 与 `EveryBy` 相同，如果窗口覆盖缓存中的全部数据，当前实现返回 `false`。

## 常见脚本示例

### 温度持续过高

最近 5 分钟内所有温度都大于 80：

```expr
temperature.EveryBy("x > 80", "5m")
```

### 温度异常波动

最近 5 分钟标准差超过 3，并且最新值相对前值变化超过 10%：

```expr
temperature.StdDev("5m") > 3 && temperature.PctChangeExceeds(10)
```

### 温度出现过越界值

最近 10 分钟至少出现一次低于 0 或高于 100：

```expr
temperature.SomeBy("x < 0 || x > 100", "10m")
```

### 设备刚刚启动

```expr
running.Rising()
```

### 设备频繁启停

最近 10 分钟上升沿或下降沿超过 5 次：

```expr
running.RC("10m") > 5 || running.FC("10m") > 5
```

### 检查状态字中的报警位

第 3 位为 `1`：

```expr
status.Bit(3)
```

低 4 位的值等于 `5`：

```expr
status.BitAnd(15) == 5
```

### 数据量保护

只有缓存数据足够时才计算标准差：

```expr
temperature.Len() >= 2 && temperature.StdDev("5m") > 3
```

## 使用注意事项

- 调用数值方法前，应确保变量是数值型。
- 调用 `Rising`、`Falling`、`RC`、`FC` 前，应确保变量是布尔型。
- 调用位操作前，应确保变量是 `uint32`、`uint64` 或字节数组类型。
- `PctChange`、`Diff`、`Changed`、`Rising` 和 `Falling` 都依赖至少两个最新数据点。
- `DiffWith` 和 `PctChangeWith` 只需要一个最新数据点。
- `StdDev` 至少需要两个窗口内数据点。
- `PctChangeSince` 和 `DiffSince` 需要在指定时间点之前存在历史数据。
- `PctChangeExceeds` 和 `DiffExceeds` 比较的是绝对值，因此上升和下降都可能触发。
- 阈值判断使用严格大于 `>`；变化量刚好等于阈值时返回 `false`。
- 条件查询的条件参数必须是返回布尔值的 expr 表达式字符串。
