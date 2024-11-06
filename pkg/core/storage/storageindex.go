package storage

import (
	F "hello/pkg/util"
	"sort"
)

type StatusIndex struct {
	localIndexes     map[string]map[string]StorageIndexCursor
	universalIndexes map[string]map[string]StorageIndexCursor
}

func NewStatusIndex() *StatusIndex {
	return &StatusIndex{
		localIndexes:     make(map[string]map[string]StorageIndexCursor),
		universalIndexes: make(map[string]map[string]StorageIndexCursor),
	}
}

func (s *StatusIndex) LocalIndexes(keys []string) map[string]StorageIndexCursor {
	indexes := make(map[string]StorageIndexCursor)

	for _, key := range keys {
		key = F.FillHash(key)
		prefix, suffix := s.Split(key)

		if prefixMap, ok := s.localIndexes[prefix]; ok {
			if cursor, exists := prefixMap[suffix]; exists {
				indexes[key] = cursor
			}
		}
	}

	return indexes
}

func (s *StatusIndex) UniversalIndexes(keys []string) map[string]StorageIndexCursor {
	indexes := make(map[string]StorageIndexCursor)

	for _, key := range keys {
		key = F.FillHash(key)
		prefix, suffix := s.Split(key)

		if prefixMap, ok := s.universalIndexes[prefix]; ok {
			if cursor, exists := prefixMap[suffix]; exists {
				indexes[key] = cursor
			}
		}
	}

	return indexes
}

func (s *StatusIndex) AddLocalIndexes(indexes map[string]StorageIndexCursor) bool {
	for key, _ := range indexes {
		key = F.FillHash(key)
		prefix, suffix := s.Split(key)

		if s.localIndexes[prefix] == nil {
			s.localIndexes[prefix] = make(map[string]StorageIndexCursor)
		}

		cursor := StorageIndexCursor{
			Key: key,
			// 다른 필드들도 필요에 따라 설정
		}
		s.localIndexes[prefix][suffix] = cursor
	}
	return true
}

func (s *StatusIndex) AddUniversalIndexes(indexes map[string]StorageIndexCursor) bool {
	for key, _ := range indexes {
		key = F.FillHash(key)
		prefix, suffix := s.Split(key)

		if s.universalIndexes[prefix] == nil {
			s.universalIndexes[prefix] = make(map[string]StorageIndexCursor)
		}

		cursor := StorageIndexCursor{
			Key: key,
			// 다른 필드들도 필요에 따라 설정
		}
		s.universalIndexes[prefix][suffix] = cursor
	}
	return true
}

func (s *StatusIndex) Split(key string) (string, string) {
	prefix := key[:STATUS_PREFIX_SIZE]
	suffix := key[STATUS_PREFIX_SIZE:]
	return prefix, suffix
}

// SearchLocalIndexes는 페이지네이션을 지원하는 로컬 인덱스 검색
func (s *StatusIndex) SearchLocalIndexes(item []interface{}) map[string]StorageIndexCursor {
	indexes := make(map[string]StorageIndexCursor)
	if len(item) == 0 {
		return indexes
	}

	prefix := item[0].(string)
	if _, ok := s.localIndexes[prefix]; ok {
		page := 0
		count := 50
		if len(item) > 1 {
			page = item[1].(int)
		}
		if len(item) > 2 {
			count = item[2].(int)
		}

		offset := page * count
		keys := make([]string, 0, len(s.localIndexes[prefix]))
		for k := range s.localIndexes[prefix] {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		end := offset + count
		if end > len(keys) {
			end = len(keys)
		}
		if offset < len(keys) {
			for _, k := range keys[offset:end] {
				indexes[prefix+k] = s.localIndexes[prefix][k]
			}
		}
	}

	return indexes
}

// SearchUniversalIndexes도 SearchLocalIndexes와 유사한 방식으로 구현
