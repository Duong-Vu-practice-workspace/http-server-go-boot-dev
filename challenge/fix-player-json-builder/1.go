package main

import "encoding/json"

type player struct {
	Name   string `json:"name"`
	Level  int    `json:"level"`
	Online bool   `json:"online"`
}

func buildPlayerJSON(name string, level int, online bool) string {
	p := player{
		Name:   name,
		Level:  level,
		Online: online,
	}

	var result string
	res, err := json.Marshal(p)
	result = string(res)
	if err != nil {
		return ""
	}
	return result
}
