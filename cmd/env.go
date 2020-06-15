package cmd

import "os"

// Environment environment variable management
type Environment struct {
	// values   map[string]string
	cache    map[string]interface{}
	bindings []string
}

// NewEnvironment create new manager
func NewEnvironment() *Environment {
	return &Environment{
		// values: make(map[string]string),
	}
}

// // Add add a key pair
// func (e Environment) Add(key, value string) {
// 	if key != "" && value != "" {
// 		e.values[key] = value
// 	}
// }

// BindEnv bind a env name
func (e *Environment) BindEnv(key string) {
	e.bindings = append(e.bindings, key)
}

// // Remove remove a key
// func (e Environment) Remove(key string) bool {
// 	if _, ok := e.values[key]; ok && key != "" {
// 		delete(e.values, key)
// 		return true
// 	}
// 	return false
// }

// Remove remove a key
// func (e Environment) Remove(key string) bool {
// 	if _, ok := e.values[key]; ok && key != "" {
// 		delete(e.values, key)
// 		return true
// 	}
// 	return false
// }

// Len returns internal map length
// func (e Environment) Len() int {
// 	return len(e.values)
// }

// Len returns internal map length
func (e Environment) Len() int {
	return len(e.bindings)
}

// ToArgs to "key=value" format slice
// func (e Environment) ToArgs() []string {
// 	var result []string
// 	for k, v := range e.values {
// 		result = append(result, k+"="+v)
// 	}
// 	return result
// }

// ToArgs to "key=value" format slice
func (e Environment) ToArgs() []string {
	var result []string
	for _, key := range e.bindings {
		if val := os.Getenv(key); val != "" {
			result = append(result, key+"="+val)
		}
	}
	return result
}
