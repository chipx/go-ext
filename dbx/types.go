package dbx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonStringArray []string

func (s JsonStringArray) Value() (driver.Value, error) {
	j, err := json.Marshal(s)

	return j, err
}

func (s *JsonStringArray) Scan(src interface{}) error {
	if src == nil {
		*s = make([]string, 0)
		return nil
	}
	source, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Type assertion .([]byte) failed. ")
	}

	var i []string
	err := json.Unmarshal(source, &i)
	if err != nil {
		return fmt.Errorf("Unmarshal data to []string failed: %s ", err)
	}

	*s = i
	return nil
}

type KeyValueStruct struct {
	Key   string
	Value string
}

type JsonKeyValueStructArray []KeyValueStruct

func (s JsonKeyValueStructArray) GetValueOrEmpty(key string) string {
	for i := 0; i < len(s); i++ {
		if s[i].Key == key {
			return s[i].Value
		}
	}

	return ""
}

func (s *JsonKeyValueStructArray) AddValue(key string, value string) *JsonKeyValueStructArray {
	*s = append(*s, KeyValueStruct{
		Key:   key,
		Value: value,
	})

	return s
}

func (s JsonKeyValueStructArray) Value() (driver.Value, error) {
	j, err := json.Marshal(s)

	return j, err
}

func (s *JsonKeyValueStructArray) Scan(src interface{}) error {
	if src == nil {
		*s = make([]KeyValueStruct, 0)
		return nil
	}
	source, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Type assertion .([]byte) failed. ")
	}

	var i []KeyValueStruct
	err := json.Unmarshal(source, &i)
	if err != nil {
		return fmt.Errorf("Unmarshal data to []KeyValueStruct failed: %s ", err)
	}

	*s = i
	return nil
}

type JsonMapInterface map[string]interface{}

func (m *JsonMapInterface) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Type assertion .([]byte) failed. ")
	}

	if err := json.Unmarshal(source, m); err != nil {
		return fmt.Errorf("Unmarshal data to map[string]interface failed: %s ", err)
	}

	return nil
}

func (m *JsonMapInterface) Value() (driver.Value, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return string(data), nil
}
