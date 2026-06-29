package service

import "encoding/json"

// unmarshalJSON 反序列化 JSON（辅助函数，集中管理 import）
func unmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// marshalJSON 序列化为 JSON
func marshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalTags 将标签切片序列化为 JSON（供 controller 调用）
func MarshalTags(tags []string) ([]byte, error) {
	return json.Marshal(tags)
}
