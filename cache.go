package edgeexpr

import (
	"errors"
	"math"
	"sync"
	"time"
)

type Point[T float64 | bool | string | []byte] struct {
	Value     T
	Timestamp *time.Time
}

type Cache[T float64 | bool | string | []byte] struct {
	Points         []Point[T]
	ExpireDuration time.Duration
	mu             sync.RWMutex // 读写锁保护Points切片
}

func NewCache[T float64 | bool | string | []byte](expireDuration time.Duration) *Cache[T] {
	return &Cache[T]{
		Points:         make([]Point[T], 0),
		ExpireDuration: expireDuration,
	}
}

func (c *Cache[T]) Value() T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var zeroValue T
	if len(c.Points) == 0 {
		return zeroValue
	}
	return c.Points[len(c.Points)-1].Value
}

func (c *Cache[T]) Latest() T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var zeroValue T
	if len(c.Points) == 0 {
		return zeroValue
	}
	return c.Points[len(c.Points)-1].Value
}

// Timestamp returns the timestamp of the latest value
func (c *Cache[T]) Timestamp() *time.Time {
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
func (c *Cache[T]) Point() *Point[T] {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return nil
	}
	// 返回最新点的副本
	latest := c.Points[len(c.Points)-1]
	return &latest
}

// Len returns the number of points in the cache
func (c *Cache[T]) Len() int {
	if c == nil {
		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.Points)
}

// MA calculates Moving Average within the specified time window
func (c *Cache[T]) MA(window string) (float64, error) {
	points := c.getPointsInWindow(window)
	if len(points) == 0 {
		return 0, nil
	}

	// 使用类型断言检查是否为 float64
	var sum float64
	for _, point := range points {
		if val, ok := any(point.Value).(float64); ok {
			sum += val
		} else {
			return 0, errors.New("value is not a float64 type")
		}
	}
	mean := sum / float64(len(points))
	return mean, nil
}

// StdDev calculates Standard Deviation within the specified time window
func (c *Cache[T]) StdDev(window string) (float64, error) {
	points := c.getPointsInWindow(window)
	if len(points) == 0 {
		return 0, nil
	}

	if len(points) == 1 {
		return 0, nil // 单个点的标准差为0
	}

	// 检查所有值是否为 float64 类型并计算平均值
	var sum float64
	var values []float64

	for _, point := range points {
		if val, ok := any(point.Value).(float64); ok {
			sum += val
			values = append(values, val)
		} else {
			return 0, errors.New("value is not a float64 type")
		}
	}

	mean := sum / float64(len(values))

	// 计算方差
	var variance float64
	for _, val := range values {
		diff := val - mean
		variance += diff * diff
	}
	variance = variance / float64(len(values))

	// 计算标准差（方差的平方根）
	standardDeviation := math.Sqrt(variance)
	return standardDeviation, nil
}

// PctChange calculates Percentage Change between the latest two points
func (c *Cache[T]) PctChange() (float64, error) {
	if c == nil {
		return 0, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 2 {
		return 0, nil
	}

	// 获取最新的两个点
	currentVal, ok1 := any(c.Points[len(c.Points)-1].Value).(float64)
	previousVal, ok2 := any(c.Points[len(c.Points)-2].Value).(float64)

	if !ok1 || !ok2 {
		return 0, errors.New("value is not a float64 type")
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

// Diff calculates the difference between the latest two points (current - previous)
func (c *Cache[T]) Diff() (float64, error) {
	if c == nil {
		return 0, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 2 {
		return 0, nil
	}

	// 获取最新的两个点
	currentVal, ok1 := any(c.Points[len(c.Points)-1].Value).(float64)
	previousVal, ok2 := any(c.Points[len(c.Points)-2].Value).(float64)

	if !ok1 || !ok2 {
		return 0, errors.New("value is not a float64 type")
	}

	// 计算差值：current - previous
	difference := currentVal - previousVal
	return difference, nil
}

// PctChangeExceeds checks if the percentage change between the latest two points exceeds the specified threshold
func (c *Cache[T]) PctChangeExceeds(threshold float64) (bool, error) {
	pctChange, err := c.PctChange()
	if err != nil {
		return false, err
	}

	// 使用绝对值比较，因为超过阈值可能是正向或负向的
	return math.Abs(pctChange) > threshold, nil
}

// DiffExceeds checks if the absolute difference between the latest two points exceeds the specified threshold
func (c *Cache[T]) DiffExceeds(threshold float64) (bool, error) {
	diff, err := c.Diff()
	if err != nil {
		return false, err
	}

	// 使用绝对值比较，因为超过阈值可能是正向或负向的
	return math.Abs(diff) > threshold, nil
}

// Changed checks if the latest two values are different
func (c *Cache[T]) Changed() bool {
	if c == nil {
		return false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 2 {
		return false
	}

	// 比较最新的两个点的值是否不同
	return !isValueEqual(c.Points[len(c.Points)-1].Value, c.Points[len(c.Points)-2].Value)
}

// PctChangeSince calculates Percentage Change between the latest value and the value from the specified time window ago
func (c *Cache[T]) PctChangeSince(window string) (float64, error) {
	if c == nil {
		return 0, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return 0, nil
	}

	// 获取最新值
	currentVal, ok := any(c.Points[len(c.Points)-1].Value).(float64)
	if !ok {
		return 0, errors.New("value is not a float64 type")
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
			if val, ok := any(c.Points[i].Value).(float64); ok {
				baseVal = val
				found = true
				break
			} else {
				return 0, errors.New("value is not a float64 type")
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
func (c *Cache[T]) DiffSince(window string) (float64, error) {
	if c == nil {
		return 0, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return 0, nil
	}

	// 获取最新值
	currentVal, ok := any(c.Points[len(c.Points)-1].Value).(float64)
	if !ok {
		return 0, errors.New("value is not a float64 type")
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
			if val, ok := any(c.Points[i].Value).(float64); ok {
				baseVal = val
				found = true
				break
			} else {
				return 0, errors.New("value is not a float64 type")
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

func (c *Cache[T]) Count(window string) int {
	points := c.getPointsInWindow(window)
	if len(points) <= 1 {
		return len(points)
	}

	// 计算数据变化次数，相邻重复的不计数
	changeCount := 1 // 第一个点始终计数

	for i := 1; i < len(points); i++ {
		// 比较当前点与前一个点的值是否不同
		if !isValueEqual(points[i].Value, points[i-1].Value) {
			changeCount++
		}
	}

	return changeCount
}

// getPointsInWindow gets points within the specified time window
// This method will acquire its own read lock
func (c *Cache[T]) getPointsInWindow(window string) []Point[T] {
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
		result := make([]Point[T], len(c.Points))
		copy(result, c.Points)
		return result
	}

	now := time.Now()
	cutoffTime := now.Add(-duration)

	var result []Point[T]
	for _, point := range c.Points {
		if point.Timestamp != nil && point.Timestamp.After(cutoffTime) {
			result = append(result, point)
		}
	}

	return result
}

// 辅助函数：比较两个值是否相等，处理不同类型
func isValueEqual[T float64 | bool | string | []byte](a, b T) bool {
	// 使用 any 类型转换来处理不同类型的比较
	aVal := any(a)
	bVal := any(b)

	switch aVal := aVal.(type) {
	case []byte:
		if bBytes, ok := bVal.([]byte); ok {
			if len(aVal) != len(bBytes) {
				return false
			}
			for i := range aVal {
				if aVal[i] != bBytes[i] {
					return false
				}
			}
			return true
		}
		return false
	default:
		return aVal == bVal
	}
}

// Only for bool type
// newest point is true and second newest point is false
func (c *Cache[T]) Rising() (bool, error) {
	if c == nil {
		return false, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 2 {
		return false, nil
	}

	if val, ok := any(c.Points[len(c.Points)-1].Value).(bool); ok && val {
		if val, ok := any(c.Points[len(c.Points)-2].Value).(bool); ok && !val {
			return true, nil
		} else {
			return false, errors.New("value is not a bool type")
		}
	} else {
		return false, errors.New("value is not a bool type")
	}
}

func (c *Cache[T]) Falling() (bool, error) {
	if c == nil {
		return false, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) < 2 {
		return false, nil
	}

	if val, ok := any(c.Points[len(c.Points)-1].Value).(bool); ok && !val {
		if val, ok := any(c.Points[len(c.Points)-2].Value).(bool); ok && val {
			return true, nil
		} else {
			return false, errors.New("value is not a bool type")
		}
	} else {
		return false, errors.New("value is not a bool type")
	}
}

// RC calculates Rising Count (false to true transitions) within the specified time window
func (c *Cache[T]) RC(window string) (int, error) {
	points := c.getPointsInWindow(window)
	if len(points) < 2 {
		return 0, nil
	}

	// 检查是否为 bool 类型
	if _, ok := any(points[0].Value).(bool); !ok {
		return 0, errors.New("value is not a bool type")
	}

	risingCount := 0
	for i := 1; i < len(points); i++ {
		prevVal, ok1 := any(points[i-1].Value).(bool)
		currVal, ok2 := any(points[i].Value).(bool)

		if !ok1 || !ok2 {
			return 0, errors.New("value is not a bool type")
		}

		// 从 false 到 true 的变化
		if !prevVal && currVal {
			risingCount++
		}
	}

	return risingCount, nil
}

// FC calculates Falling Count (true to false transitions) within the specified time window
func (c *Cache[T]) FC(window string) (int, error) {
	points := c.getPointsInWindow(window)
	if len(points) < 2 {
		return 0, nil
	}

	// 检查是否为 bool 类型
	if _, ok := any(points[0].Value).(bool); !ok {
		return 0, errors.New("value is not a bool type")
	}

	fallingCount := 0
	for i := 1; i < len(points); i++ {
		prevVal, ok1 := any(points[i-1].Value).(bool)
		currVal, ok2 := any(points[i].Value).(bool)

		if !ok1 || !ok2 {
			return 0, errors.New("value is not a bool type")
		}

		// 从 true 到 false 的变化
		if prevVal && !currVal {
			fallingCount++
		}
	}

	return fallingCount, nil
}

// Only for []byte type
// []byte value act like a whole bit array
// index is the bit position, starting from 0
// returns true if the bit at the specified index is set
func (c *Cache[T]) Bit(index int) (bool, error) {
	if c == nil {
		return false, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return false, nil
	}

	if val, ok := any(c.Points[len(c.Points)-1].Value).([]byte); ok {
		if index >= 0 && index < len(val)*8 {
			byteIndex := index / 8
			bitIndex := index % 8
			if byteIndex < len(val) {
				return (val[byteIndex] & (1 << bitIndex)) != 0, nil
			} else {
				return false, errors.New("index out of range")
			}
		} else {
			return false, errors.New("index out of range")
		}
	} else {
		return false, errors.New("value is not a []byte type")
	}
}

// ByteBit returns the i-th bit of the n-th byte in the latest []byte value
// ByteBit(n, i) gets bit i (0-7) from byte n (0-based indexing)
func (c *Cache[T]) ByteBit(n, i int) (bool, error) {
	if c == nil {
		return false, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Points) == 0 {
		return false, nil
	}

	if val, ok := any(c.Points[len(c.Points)-1].Value).([]byte); ok {
		// 检查字节索引是否有效
		if n < 0 || n >= len(val) {
			return false, errors.New("byte index out of range")
		}

		// 检查位索引是否有效 (0-7)
		if i < 0 || i > 7 {
			return false, errors.New("bit index out of range (must be 0-7)")
		}

		// 获取第n个字节的第i位
		return (val[n] & (1 << i)) != 0, nil
	} else {
		return false, errors.New("value is not a []byte type")
	}
}

func (c *Cache[T]) AddPoint(value T, timestamp *time.Time) {
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

	c.Points = append(c.Points, Point[T]{Value: value, Timestamp: timestamp})
	c.cleanExpiredPointsUnsafe()
}

func (c *Cache[T]) cleanExpiredPointsUnsafe() {
	if c.ExpireDuration <= 0 {
		return
	}

	now := time.Now()
	validPoints := make([]Point[T], 0, len(c.Points))

	for _, point := range c.Points {
		if point.Timestamp != nil && now.Sub(*point.Timestamp) <= c.ExpireDuration {
			validPoints = append(validPoints, point)
		}
	}

	c.Points = validPoints
}
