package edgeexpr

import "time"

type Point[T float64 | bool | string | []byte] struct {
	Value     T
	Timestamp *time.Time
}

type Cache[T float64 | bool | string | []byte] struct {
	Points         []Point[T]
	ExpireDuration time.Duration
}

func NewCache[T float64 | bool | string | []byte](expireDuration time.Duration) *Cache[T] {
	return &Cache[T]{
		Points:         make([]Point[T], 0),
		ExpireDuration: expireDuration,
	}
}

func (c *Cache[T]) Value() T {
	var zeroValue T
	if len(c.Points) == 0 {
		return zeroValue
	}
	if len(c.Points) == 0 {
		return zeroValue
	}
	return c.Points[len(c.Points)-1].Value
}

func (c *Cache[T]) Mean() float64 {
	if len(c.Points) == 0 {
		return 0
	}

	// 使用类型断言检查是否为 float64
	var sum float64
	for _, point := range c.Points {
		if val, ok := any(point.Value).(float64); ok {
			sum += val
		}
	}
	mean := sum / float64(len(c.Points))
	return mean
}

func (c *Cache[T]) Count() int {
	if c == nil || len(c.Points) == 0 {
		return 0
	}

	// 使用 map 来统计唯一值
	uniqueValues := make(map[any]struct{})

	for _, point := range c.Points {
		uniqueValues[point.Value] = struct{}{}
	}

	return len(uniqueValues)
}

// Only for bool type
// newest point is true and second newest point is false
func (c *Cache[T]) Rising() bool {
	if c == nil || len(c.Points) < 2 {
		return false
	}

	if val, ok := any(c.Points[len(c.Points)-1].Value).(bool); ok && val {
		if val, ok := any(c.Points[len(c.Points)-2].Value).(bool); ok && !val {
			return true
		}
	}
	return false
}

func (c *Cache[T]) Falling() bool {
	if c == nil || len(c.Points) < 2 {
		return false
	}

	if val, ok := any(c.Points[len(c.Points)-1].Value).(bool); ok && !val {
		if val, ok := any(c.Points[len(c.Points)-2].Value).(bool); ok && val {
			return true
		}
	}
	return false
}

// Only for []byte type
// []byte value act like a whole bit array
// index is the bit position, starting from 0
// returns true if the bit at the specified index is set
func (c *Cache[T]) Bit(index int) bool {
	if c == nil || len(c.Points) == 0 {
		return false
	}

	if val, ok := any(c.Points[len(c.Points)-1].Value).([]byte); ok && index >= 0 && index < len(val)*8 {
		byteIndex := index / 8
		bitIndex := index % 8
		if byteIndex < len(val) {
			return (val[byteIndex] & (1 << bitIndex)) != 0
		}
	}
	return false
}

func (c *Cache[T]) AddPoint(value T, timestamp *time.Time) {
	if timestamp == nil {
		now := time.Now()
		timestamp = &now
	}

	// 检查并清理过期的点
	c.cleanExpiredPoints()

	c.Points = append(c.Points, Point[T]{Value: value, Timestamp: timestamp})
}

func (c *Cache[T]) cleanExpiredPoints() {
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
