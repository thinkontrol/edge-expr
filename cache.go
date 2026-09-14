package edgeexpr

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"

	"github.com/samber/lo"
)

type Point struct {
	Value     any
	Timestamp *time.Time
}

type Cache struct {
	Points         []Point
	ExpireDuration time.Duration
	mu             sync.RWMutex // 读写锁保护Points切片
}

func NewCache(expireDuration *time.Duration) *Cache {
	cache := &Cache{
		Points:         make([]Point, 0),
		ExpireDuration: time.Minute,
	}
	if expireDuration != nil {
		cache.ExpireDuration = *expireDuration
	}
	return cache
}

// func (c *Cache) PushValue() *PushValue {
// 	c.mu.RLock()
// 	defer c.mu.RUnlock()

// 	if len(c.Points) == 0 {
// 		return nil
// 	}
// 	return &PushValue{
// 		// Key:       key,
// 		Value:     c.Points[len(c.Points)-1].Value,
// 		Timestamp: c.Points[len(c.Points)-1].Timestamp,
// 	}
// }

func (c *Cache) Value() any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return nil
	}
	return c.Points[len(c.Points)-1].Value
}

// func (c *Cache) Latest() any {
// 	c.mu.RLock()
// 	defer c.mu.RUnlock()

// 	if len(c.Points) == 0 {
// 		return nil
// 	}
// 	return c.Points[len(c.Points)-1].Value
// }

// Timestamp returns the timestamp of the latest value
func (c *Cache) Timestamp() *time.Time {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return nil
	}
	return c.Points[len(c.Points)-1].Timestamp
}

// Point returns the latest point (value and timestamp)
// func (c *Cache) Point() *Point {
// 	if c == nil {
// 		return nil
// 	}

// 	c.mu.RLock()
// 	defer c.mu.RUnlock()

// 	if len(c.Points) == 0 {
// 		return nil
// 	}
// 	// 返回最新点的副本
// 	latest := c.Points[len(c.Points)-1]
// 	return &latest
// }

// Len returns the number of points in the cache
func (c *Cache) Len() int {
	if c == nil {
		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.Points)
}

// MA calculates Moving Average within the specified time window
func (c *Cache) MA(window string) (float64, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}
	points := c.getPointsInWindow(window)
	if len(points) == 0 {
		return 0, fmt.Errorf("no data yet")
	}

	// 使用类型断言检查是否为 float64
	return lo.MeanByErr(points, func(p Point) (float64, error) {
		return ConvertToFloat64(p.Value)
	})
}

// StdDev calculates Standard Deviation within the specified time window
func (c *Cache) StdDev(window string) (float64, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}
	points := c.getPointsInWindow(window)
	if len(points) == 0 {
		return 0, fmt.Errorf("no data yet")
	}

	if len(points) == 1 {
		return 0, fmt.Errorf("at least two data points are required to calculate standard deviation")
	}

	mean, err := lo.MeanByErr(points, func(p Point) (float64, error) {
		return ConvertToFloat64(p.Value)
	})
	if err != nil {
		return 0, err
	}

	squaredDiffSum, err := lo.SumByErr(points, func(p Point) (float64, error) {
		f, err := ConvertToFloat64(p.Value)
		if err != nil {
			return 0, err
		}
		return math.Pow(f-mean, 2), nil
	})

	if err != nil {
		return 0, err
	}

	// 计算方差
	variance := squaredDiffSum / float64(len(points))
	// 计算标准差（方差的平方根）
	standardDeviation := math.Sqrt(variance)
	return standardDeviation, nil
}

// PctChange calculates Percentage Change between the latest two points
func (c *Cache) PctChange() (float64, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 2 {
		return 0, fmt.Errorf("not enough data points")
	}

	// 获取最新的两个点
	currentVal, err1 := ConvertToFloat64(c.Points[len(c.Points)-1].Value)
	previousVal, err2 := ConvertToFloat64(c.Points[len(c.Points)-2].Value)

	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("value(%v)[%T] or (%v)[%T] is not a float64 type", c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-2].Value, c.Points[len(c.Points)-2].Value)
	}

	// 如果前一个值为0，无法计算百分比变化
	if previousVal == 0 {
		if currentVal == 0 {
			return 0, nil // 0到0没有变化
		}
		return 0, errors.New("cannot calculate percentage change from zero")
	}

	// 计算百分比变化：((current - previous) / previous) * 100
	percentageChange := ((currentVal - previousVal) / previousVal) * 100
	return percentageChange, nil
}

func (c *Cache) PctChangeWith(val float64) (float64, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 1 {
		return 0, fmt.Errorf("not enough data points")
	}

	// 获取最新的两个点
	currentVal, err1 := ConvertToFloat64(c.Points[len(c.Points)-1].Value)
	if err1 != nil {
		return 0, fmt.Errorf("value(%v)[%T] is not a float64 type", c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-1].Value)
	}

	// 如果前一个值为0，无法计算百分比变化
	if val == 0 {
		if currentVal == 0 {
			return 0, nil // 0到0没有变化
		}
		return 0, errors.New("cannot calculate percentage change from zero")
	}

	// 计算百分比变化：((current - previous) / previous) * 100
	percentageChange := ((currentVal - val) / val) * 100
	return percentageChange, nil
}

// Diff calculates the difference between the latest two points (current - previous)
func (c *Cache) Diff() (float64, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 2 {
		return 0, fmt.Errorf("not enough data points")
	}

	// 获取最新的两个点
	currentVal, err1 := ConvertToFloat64(c.Points[len(c.Points)-1].Value)
	previousVal, err2 := ConvertToFloat64(c.Points[len(c.Points)-2].Value)

	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("value(%v)[%T] or (%v)[%T] is not a float64 type", c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-2].Value, c.Points[len(c.Points)-2].Value)
	}

	// 计算差值：current - previous
	difference := currentVal - previousVal
	return difference, nil
}

func (c *Cache) DiffWith(val float64) (float64, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 21 {
		return 0, fmt.Errorf("not enough data points")
	}

	// 获取最新的两个点
	currentVal, err1 := ConvertToFloat64(c.Points[len(c.Points)-1].Value)
	if err1 != nil {
		return 0, fmt.Errorf("value(%v)[%T] is not a float64 type", c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-1].Value)
	}

	// 计算差值：current - previous
	difference := currentVal - val
	return difference, nil
}

// PctChangeExceeds checks if the percentage change between the latest two points exceeds the specified threshold
func (c *Cache) PctChangeExceeds(threshold float64) (bool, error) {
	if c == nil {
		return false, fmt.Errorf("cache is nil")
	}
	pctChange, err := c.PctChange()
	if err != nil {
		return false, err
	}

	// 使用绝对值比较，因为超过阈值可能是正向或负向的
	return math.Abs(pctChange) > threshold, nil
}

// DiffExceeds checks if the absolute difference between the latest two points exceeds the specified threshold
func (c *Cache) DiffExceeds(threshold float64) (bool, error) {
	if c == nil {
		return false, fmt.Errorf("cache is nil")
	}
	diff, err := c.Diff()
	if err != nil {
		return false, err
	}

	// 使用绝对值比较，因为超过阈值可能是正向或负向的
	return math.Abs(diff) > threshold, nil
}

// Changed checks if the latest two values are different
func (c *Cache) Changed() bool {
	if c == nil {
		return false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 2 {
		return false
	}

	// 比较最新的两个点的值是否不同
	return !reflect.DeepEqual(c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-2].Value)
}

// PctChangeSince calculates Percentage Change between the latest value and the value from the specified time window ago
func (c *Cache) PctChangeSince(window string) (float64, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return 0, fmt.Errorf("no data points available")
	}

	// 获取最新值
	currentVal, err := ConvertToFloat64(c.Points[len(c.Points)-1].Value)
	if err != nil {
		return 0, fmt.Errorf("value(%v)[%T] is not a float64 type", c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-1].Value)
	}

	// 解析时间窗口
	duration, err := time.ParseDuration(window)
	if err != nil {
		return 0, errors.New("invalid time window format")
	}

	// 计算目标时间点
	now := time.Now()
	targetTime := now.Add(-duration)

	// 找到时间窗口前最接近的点
	var baseVal float64
	var found bool

	for i := len(c.Points) - 1; i >= 0; i-- {
		if c.Points[i].Timestamp != nil && c.Points[i].Timestamp.Before(targetTime) {
			if val, err := ConvertToFloat64(c.Points[i].Value); err == nil {
				baseVal = val
				found = true
				break
			} else {
				return 0, fmt.Errorf("value(%v)[%T] is not a float64 type", c.Points[i].Value, c.Points[i].Value)
			}
		}
	}

	if !found {
		return 0, errors.New("no data point found before the specified time window")
	}

	// 如果基准值为0，无法计算百分比变化
	if baseVal == 0 {
		if currentVal == 0 {
			return 0, nil // 0到0没有变化
		}
		return 0, errors.New("cannot calculate percentage change from zero")
	}

	// 计算百分比变化：((current - base) / base) * 100
	percentageChange := ((currentVal - baseVal) / baseVal) * 100
	return percentageChange, nil
}

// DiffSince calculates the difference between the latest value and the value from the specified time window ago
func (c *Cache) DiffSince(window string) (float64, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return 0, fmt.Errorf("no data points available")
	}

	// 获取最新值
	currentVal, err := ConvertToFloat64(c.Points[len(c.Points)-1].Value)
	if err != nil {
		return 0, fmt.Errorf("value(%v)[%T] is not a float64 type", c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-1].Value)
	}

	// 解析时间窗口
	duration, err := time.ParseDuration(window)
	if err != nil {
		return 0, errors.New("invalid time window format")
	}

	// 计算目标时间点
	now := time.Now()
	targetTime := now.Add(-duration)

	// 找到时间窗口前最接近的点
	var baseVal float64
	var found bool

	for i := len(c.Points) - 1; i >= 0; i-- {
		if c.Points[i].Timestamp != nil && c.Points[i].Timestamp.Before(targetTime) {
			if val, err := ConvertToFloat64(c.Points[i].Value); err == nil {
				baseVal = val
				found = true
				break
			} else {
				return 0, fmt.Errorf("value(%v)[%T] is not a float64 type", c.Points[i].Value, c.Points[i].Value)
			}
		}
	}

	if !found {
		return 0, errors.New("no data point found before the specified time window")
	}

	// 计算差值：current - base
	difference := currentVal - baseVal
	return difference, nil
}

func (c *Cache) Count(window string) int {
	points := c.getPointsInWindow(window)
	if len(points) <= 1 {
		return len(points)
	}

	// 计算数据变化次数，相邻重复的不计数
	changeCount := 1 // 第一个点始终计数

	for i := 1; i < len(points); i++ {
		// 比较当前点与前一个点的值是否不同
		if !reflect.DeepEqual(points[i].Value, points[i-1].Value) {
			changeCount++
		}
	}

	return changeCount
}

// getPointsInWindow gets points within the specified time window
// This method will acquire its own read lock
func (c *Cache) getPointsInWindow(window string) []Point {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return nil
	}

	// 解析时间窗口字符串
	duration, err := time.ParseDuration(window)
	if err != nil {
		// 如果解析失败，返回所有点的副本
		result := make([]Point, len(c.Points))
		copy(result, c.Points)
		return result
	}

	now := time.Now()
	cutoffTime := now.Add(-duration)

	var result []Point
	for _, point := range c.Points {
		if point.Timestamp != nil && point.Timestamp.After(cutoffTime) {
			result = append(result, point)
		}
	}

	return result
}

// Only for bool type
// newest point is true and second newest point is false
func (c *Cache) Rising() (bool, error) {
	if c == nil {
		return false, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 2 {
		return false, nil
	}

	if val, ok := c.Points[len(c.Points)-1].Value.(bool); ok && val {
		if val, ok := c.Points[len(c.Points)-2].Value.(bool); ok && !val {
			return true, nil
		} else {
			return false, fmt.Errorf("value(%v)[%T] is not a bool type", c.Points[len(c.Points)-2].Value, c.Points[len(c.Points)-2].Value)
		}
	} else {
		return false, fmt.Errorf("value(%v)[%T] is not a bool type", c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-1].Value)
	}
}

func (c *Cache) Falling() (bool, error) {
	if c == nil {
		return false, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 2 {
		return false, nil
	}

	if val, ok := c.Points[len(c.Points)-1].Value.(bool); ok && !val {
		if val, ok := c.Points[len(c.Points)-2].Value.(bool); ok && val {
			return true, nil
		} else {
			return false, fmt.Errorf("value(%v)[%T] is not a bool type", c.Points[len(c.Points)-2].Value, c.Points[len(c.Points)-2].Value)
		}
	} else {
		return false, fmt.Errorf("value(%v)[%T] is not a bool type", c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-1].Value)
	}
}

// RC calculates Rising Count (false to true transitions) within the specified time window
func (c *Cache) RC(window string) (int, error) {
	points := c.getPointsInWindow(window)
	if len(points) < 2 {
		return 0, nil
	}

	// 检查是否为 bool 类型
	if _, ok := points[0].Value.(bool); !ok {
		return 0, fmt.Errorf("value(%v)[%T] is not a bool type", points[0].Value, points[0].Value)
	}

	risingCount := 0
	for i := 1; i < len(points); i++ {
		prevVal, ok1 := points[i-1].Value.(bool)
		currVal, ok2 := points[i].Value.(bool)

		if !ok1 || !ok2 {
			return 0, fmt.Errorf("value(%v)[%T] is not a bool type", points[i].Value, points[i].Value)
		}

		// 从 false 到 true 的变化
		if !prevVal && currVal {
			risingCount++
		}
	}

	return risingCount, nil
}

// FC calculates Falling Count (true to false transitions) within the specified time window
func (c *Cache) FC(window string) (int, error) {
	points := c.getPointsInWindow(window)
	if len(points) < 2 {
		return 0, nil
	}

	// 检查是否为 bool 类型
	if _, ok := points[0].Value.(bool); !ok {
		return 0, fmt.Errorf("value(%v)[%T] is not a bool type", points[0].Value, points[0].Value)
	}

	fallingCount := 0
	for i := 1; i < len(points); i++ {
		prevVal, ok1 := points[i-1].Value.(bool)
		currVal, ok2 := points[i].Value.(bool)

		if !ok1 || !ok2 {
			return 0, fmt.Errorf("value(%v)[%T] is not a bool type", points[i].Value, points[i].Value)
		}

		// 从 true 到 false 的变化
		if prevVal && !currVal {
			fallingCount++
		}
	}

	return fallingCount, nil
}

// bitValue converts supported bit values to a little-endian uint64.
func bitValue(value any) (uint64, error) {
	switch val := value.(type) {
	case []byte:
		var result uint64
		for i, b := range val {
			if i >= 8 {
				break
			}
			result |= uint64(b) << (i * 8)
		}
		return result, nil
	case uint32:
		return uint64(val), nil
	case uint64:
		return val, nil
	default:
		return 0, fmt.Errorf("value(%v)[%T] is not a supported bit type", value, value)
	}
}

func bitAt(value any, index int) (bool, error) {
	if index < 0 {
		return false, fmt.Errorf("index(%d) is out of range", index)
	}

	bitWidth := 64
	switch val := value.(type) {
	case []byte:
		bitWidth = len(val) * 8
		if index >= bitWidth {
			return false, fmt.Errorf("index(%d) is out of range", index)
		}
		return (val[index/8] & (uint8(1) << (index % 8))) != 0, nil
	case uint32:
		bitWidth = 32
	case uint64:
	default:
		return false, fmt.Errorf("value(%v)[%T] is not a supported bit type", value, value)
	}
	if index >= bitWidth {
		return false, fmt.Errorf("index(%d) is out of range", index)
	}

	converted, err := bitValue(value)
	if err != nil {
		return false, err
	}
	return (converted & (uint64(1) << index)) != 0, nil
}

func bitNot(value any) (uint, error) {
	switch val := value.(type) {
	case uint32:
		return uint(^val), nil
	case uint64:
		return uint(^val), nil
	case []byte:
		converted, err := bitValue(val)
		if err != nil {
			return 0, err
		}
		return uint(^converted), nil
	default:
		return 0, fmt.Errorf("value(%v)[%T] is not a supported bit type", value, value)
	}
}

// []byte values act like a whole bit array; integer values use their native width.
// index is the bit position, starting from 0
// returns true if the bit at the specified index is set
func (c *Cache) Bit(index int) (bool, error) {
	if c == nil {
		return false, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return false, fmt.Errorf("no data points available")
	}

	value := c.Points[len(c.Points)-1].Value
	return bitAt(value, index)
}

// ByteBit returns the i-th bit of the n-th byte in the latest bit value.
// ByteBit(n, i) gets bit i (0-7) from byte n (0-based indexing).
func (c *Cache) ByteBit(n, i int) (bool, error) {
	if c == nil {
		return false, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return false, fmt.Errorf("no data points available")
	}

	if n < 0 {
		return false, fmt.Errorf("byte index(%d) is out of range", n)
	}
	if i < 0 || i > 7 {
		return false, fmt.Errorf("bit index(%d) is out of range (must be 0-7)", i)
	}

	return bitAt(c.Points[len(c.Points)-1].Value, n*8+i)
}

// BitAnd performs a bitwise AND operation between the latest []byte value and the mask
// The []byte value is interpreted as a Little-Endian integer
func (c *Cache) BitAnd(mask uint64) (uint, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return 0, fmt.Errorf("no data points available")
	}

	value, err := bitValue(c.Points[len(c.Points)-1].Value)
	if err != nil {
		return 0, err
	}
	return uint(value & mask), nil
}

// BitOr performs a bitwise OR operation between the latest []byte value and the mask
// The []byte value is interpreted as a Little-Endian integer
func (c *Cache) BitOr(mask uint64) (uint, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return 0, fmt.Errorf("no data points available")
	}

	value, err := bitValue(c.Points[len(c.Points)-1].Value)
	if err != nil {
		return 0, err
	}
	return uint(value | mask), nil
}

// BitXor performs a bitwise XOR operation between the latest []byte value and the mask
// The []byte value is interpreted as a Little-Endian integer
func (c *Cache) BitXor(mask uint64) (uint, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return 0, fmt.Errorf("no data points available")
	}

	value, err := bitValue(c.Points[len(c.Points)-1].Value)
	if err != nil {
		return 0, err
	}
	return uint(value ^ mask), nil
}

// BitClear performs a bitwise AND NOT operation (bit clear) between the latest []byte value and the mask
// The []byte value is interpreted as a Little-Endian integer
func (c *Cache) BitClear(mask uint64) (uint, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return 0, fmt.Errorf("no data points available")
	}

	value, err := bitValue(c.Points[len(c.Points)-1].Value)
	if err != nil {
		return 0, err
	}
	return uint(value &^ mask), nil
}

// BitNot performs a bitwise NOT operation on the latest []byte value
// The []byte value is interpreted as a Little-Endian integer
func (c *Cache) BitNot() (uint, error) {
	if c == nil {
		return 0, fmt.Errorf("cache is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return 0, fmt.Errorf("no data points available")
	}

	return bitNot(c.Points[len(c.Points)-1].Value)
}

func (c *Cache) AddPoint(value any, timestamp *time.Time) {
	if c == nil {
		return
	}

	if timestamp == nil {
		now := time.Now()
		timestamp = &now
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查是否已经存在相同timestamp的point
	for i, point := range c.Points {
		if point.Timestamp != nil && timestamp != nil && point.Timestamp.Equal(*timestamp) {
			// 如果存在相同的时间戳，更新值并返回
			c.Points[i].Value = value
			c.cleanExpiredPointsUnsafe()
			return
		}
	}

	c.Points = append(c.Points, Point{Value: value, Timestamp: timestamp})
	c.cleanExpiredPointsUnsafe()
}

func (c *Cache) cleanExpiredPointsUnsafe() {
	if c.ExpireDuration <= 0 || len(c.Points) <= 1 {
		return
	}

	now := time.Now()
	validPoints := make([]Point, 0, len(c.Points))

	for _, point := range c.Points {
		if point.Timestamp != nil && now.Sub(*point.Timestamp) <= c.ExpireDuration {
			validPoints = append(validPoints, point)
		}
	}

	if len(validPoints) == 0 && len(c.Points) > 0 {
		// 保留最新的一个点，即使它已经过期
		validPoints = append(validPoints, c.Points[len(c.Points)-1])
	}

	c.Points = validPoints
}

// TODO:  增加功能： 某个时间段的变化值
// TODO: 取bit
