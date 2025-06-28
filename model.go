package edgeexpr

import (
	"crypto/md5"
	"fmt"
	"sort"
)

type Model struct {
	Connections map[string]string    `json:"connections"` // map of connection name to connection type
	Variables   map[string]*Variable `json:"variables"`   // map of variable name to Variable struct
}

func (m *Model) Hash() string {
	hash := md5.New()

	// 对 Connections 排序
	connKeys := make([]string, 0, len(m.Connections))
	for k := range m.Connections {
		connKeys = append(connKeys, k)
	}
	sort.Strings(connKeys)
	for _, k := range connKeys {
		hash.Write([]byte(fmt.Sprintf("%s:%s;", k, m.Connections[k])))
	}

	// 对 Variables 排序
	varKeys := make([]string, 0, len(m.Variables))
	for k := range m.Variables {
		varKeys = append(varKeys, k)
	}
	sort.Strings(varKeys)
	for _, k := range varKeys {
		hash.Write([]byte(fmt.Sprintf("%s:%s;", k, m.Variables[k].Hash())))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
