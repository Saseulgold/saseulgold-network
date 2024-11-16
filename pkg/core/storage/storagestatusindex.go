package storage

import (
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	F "hello/pkg/util"
	"sort"
	"sync"
)

type StatusIndex struct {
	localIndexes     map[string]map[string]StorageIndexCursor
	universalIndexes map[string]map[string]StorageIndexCursor
}

var statusIndexinstance *StatusIndex
var statusIndexonce sync.Once

func GetStatusIndexInstance() *StatusIndex {
	statusIndexonce.Do(func() {
		statusIndexinstance = &StatusIndex{
			localIndexes:     make(map[string]map[string]StorageIndexCursor),
			universalIndexes: make(map[string]map[string]StorageIndexCursor),
		}
	})
	return statusIndexinstance
}

func (s *StatusIndex) Load() {
	chainStorage := GetChainStorageInstance()

	// bundling
	statusFile := GetStatusFileInstance()
	statusFile.Cache()

	fixedHeight := LastHeight()
	bundleHeight := statusFile.BundleHeight()

	for i := bundleHeight + 1; i <= fixedHeight; i++ {
		block, err := chainStorage.GetBlock(i)
		if err != nil {
			panic(err)
		}
		statusFile.Write(block)

		if i%256 == 0 || i == fixedHeight {
			DebugLog(fmt.Sprintf("Commit Bundle: Bundling.. %d", i))
		}
	}

	statusFile.Flush()

	// update indexes
	statusFile.CopyBundles()

	// read
	localIndexes := ReadStatusStorageIndex(statusFile.LocalBundleIndex(), true)
	universalIndexes := ReadStatusStorageIndex(statusFile.UniversalBundleIndex(), true)

	// updates
	lastHeight := LastHeight()
	bundleHeight = statusFile.BundleHeight()

	DebugLog(fmt.Sprintf("Bundle Height: %d", bundleHeight))
	DebugLog(fmt.Sprintf("Last Main Block Height: %d", lastHeight))

	for i := bundleHeight + 1; i <= lastHeight; i++ {
		block, err := chainStorage.GetBlock(i)
		if err != nil {
			panic(err)
		}
		localIndexes = statusFile.UpdateLocal(localIndexes, block.LocalUpdates)
		universalIndexes = statusFile.UpdateUniversal(universalIndexes, block.UniversalUpdates)

		if i%256 == 0 || i == lastHeight {
			DebugLog(fmt.Sprintf("Update Status Datas... Height: %d", i))
		}
	}

	// cache
	s.AddLocalIndexes(localIndexes)
	s.AddUniversalIndexes(universalIndexes)
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
	prefix := key[:C.STATUS_PREFIX_SIZE]
	suffix := key[C.STATUS_PREFIX_SIZE:]
	return prefix, suffix
}

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

func (s *StatusIndex) SearchUniversalIndexes(item []interface{}) map[string]StorageIndexCursor {
	indexes := make(map[string]StorageIndexCursor)
	if len(item) == 0 {
		return indexes
	}

	prefix := item[0].(string)
	if _, ok := s.universalIndexes[prefix]; ok {
		page := 0
		count := 50
		if len(item) > 1 {
			page = item[1].(int)
		}
		if len(item) > 2 {
			count = item[2].(int)
		}

		offset := page * count
		keys := make([]string, 0, len(s.universalIndexes[prefix]))
		for k := range s.universalIndexes[prefix] {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		end := offset + count
		if end > len(keys) {
			end = len(keys)
		}
		if offset < len(keys) {
			for _, k := range keys[offset:end] {
				indexes[prefix+k] = s.universalIndexes[prefix][k]
			}
		}
	}

	return indexes
}

func (s *StatusIndex) CountUniversalIndexes(prefix string) int {
	if indexes, ok := s.universalIndexes[prefix]; ok {
		return len(indexes)
	}
	return 0
}
