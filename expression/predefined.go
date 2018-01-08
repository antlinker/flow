package expression

import (
	"encoding/json"
)

type predefined struct {
	data []pairs
}

//
func (c *predefined) PredefinedVar(key string, value string) {
	c.data = append(c.data, pairs{Key: key, Value: value})
}
func (c *predefined) PredefinedJson(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.PredefinedVar(key, string(data))
	return nil
}
