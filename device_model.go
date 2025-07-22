package edgeexpr

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/expr-lang/expr"
)

type DeviceModel struct {
	Connections map[string]string    `json:"connections"` // map of connection name to connection type
	Variables   map[string]*Variable `json:"variables"`   // map of variable name to Variable struct
}

func (m *DeviceModel) UnmarshalJSON(data []byte) error {
	type Alias DeviceModel
	aux := (*Alias)(m)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	// Initialize maps if they are nil
	if m.Connections == nil {
		m.Connections = make(map[string]string)
	}
	if m.Variables == nil {
		m.Variables = make(map[string]*Variable)
	}

	env := make(map[string]any)
	for key, variable := range m.Variables {
		if variable.Cache != nil {
			env[key] = variable.Cache
		}
	}

	keyRegex := regexp.MustCompile(`^\w+$`)

	var errs []string
	for key, variable := range m.Variables {
		if !keyRegex.MatchString(key) {
			errs = append(errs, fmt.Sprintf("Invalid variable key: %s", key))
		}
		if key != variable.Key {
			errs = append(errs, fmt.Sprintf("Variable key mismatch: %s != %s", key, variable.Key))
		}
		if variable.Connection == "" && variable.Script != "" {
			program, err := expr.Compile(variable.Script, expr.Env(env))
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %v", key, err).Error())
			} else {
				variable.Program = program
			}
		}
	}
	if len(errs) > 0 {
		sort.Strings(errs)
		return fmt.Errorf("Script errors:\n%s", strings.Join(errs, "\n"))
	}

	return nil
}

func (m *DeviceModel) Hash() string {
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
