package main

import (
	"encoding/json"
	"sort"
)

type IntMapOrderedItem struct {
	Key   string
	Value int
}

type IntMapOrdered []*IntMapOrderedItem

func (m IntMapOrdered) Len() int {
	return len(m)
}

func (m IntMapOrdered) Less(i, j int) bool {
	return m[i].Value < m[j].Value
}

func (m IntMapOrdered) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func SortMap(m IntMapOrdered, order string) IntMapOrdered {
	if order == "desc" {
		sort.Sort(sort.Reverse(m))
		return m
	}
	sort.Sort(m)
	return m
}

func Json(v interface{}) string {
	if bytes, err := json.Marshal(v); err == nil {
		return string(bytes)
	}
	return "{}"
}

func NewIntMapOrdered(m map[string]int) IntMapOrdered {
	im := make(IntMapOrdered, len(m))
	index := 0
	for key, value := range m {
		im[index] = &IntMapOrderedItem{key, value}
		index++
	}
	return im
}
