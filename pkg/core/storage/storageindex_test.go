package storage

import (
	"encoding/json"
	. "hello/pkg/core/model"
	"os"
	"path/filepath"
	"testing"
)

func TestStorageIndex(t *testing.T) {
	// 테스트용 임시 디렉토리 생성
	tempDir, err := os.MkdirTemp("", "storage_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("인덱스 생성 및 읽기 테스트", func(t *testing.T) {
		// 테스트 데이터 준비
		testKey := "testkey123456789"
		fileID := "0001"
		seek := int64(1000)
		length := int64(500)

		// 인덱스 파일 경로
		indexPath := filepath.Join(tempDir, "test_index")

		// 새 인덱스 생성
		index := NewStorageCursor(testKey, fileID, seek, length)
		if index.FileID != fileID {
			t.Errorf("FileID가 일치하지 않음: got %v, want %v", index.FileID, fileID)
		}
		if index.Seek != seek {
			t.Errorf("Seek이 일치하지 않음: got %v, want %v", index.Seek, seek)
		}
		if index.Length != length {
			t.Errorf("Length가 일치하지 않음: got %v, want %v", index.Length, length)
		}

		// 인덱스 데이터 직렬화
		indexData := IndexRaw(testKey, fileID, seek, length)

		// 파일에 쓰기
		err = os.WriteFile(indexPath, indexData, 0644)
		if err != nil {
			t.Fatal(err)
		}

		// 인덱스 읽기
		indexes := ReadStatusStorageIndex(indexPath, true)

		// 결과 검증
		readIndex, exists := indexes[testKey]
		if !exists {
			t.Error("인덱스를 찾을 수 없음")
		}
		if readIndex.FileID != fileID {
			t.Errorf("FileID가 일치하지 않음: got %v, want %v", readIndex.FileID, fileID)
		}
	})

	t.Run("데이터 크기 변경 시나리오 테스트", func(t *testing.T) {
		sf := GetStatusFileInstance()
		sf.Reset()

		// 초기 데이터 생성
		testKey := "testkey123456789"
		initialData := map[string]interface{}{
			"field1": "initial value",
			"field2": 123,
		}

		// UpdateMap 생성
		updates := make(UpdateMap)
		updates[testKey] = Update{New: initialData}

		// 첫 번째 쓰기
		err := sf.WriteUniversal(updates)
		if err != nil {
			t.Fatal(err)
		}
		if err = sf.WriteTasks(); err != nil {
			t.Fatal(err)
		}
		if err = sf.Commit(); err != nil {
			t.Fatal(err)
		}

		// 첫 번째 데이터의 인덱스 확인
		index1, exists := sf.CachedUniversalIndexes[testKey]
		if !exists {
			t.Error("인덱스를 찾을 수 없음")
		}
		initialLength := index1.Length

		// 더 큰 데이터로 업데이트
		largerData := map[string]interface{}{
			"field1": "much longer value than before",
			"field2": 123,
			"field3": "additional field with more data",
			"field4": []int{1, 2, 3, 4, 5},
		}
		updates[testKey] = Update{New: largerData}

		// 두 번째 쓰기
		err = sf.WriteUniversal(updates)
		if err != nil {
			t.Fatal(err)
		}
		if err = sf.WriteTasks(); err != nil {
			t.Fatal(err)
		}
		if err = sf.Commit(); err != nil {
			t.Fatal(err)
		}

		// 두 번째 데이터의 인덱스 확인
		index2, exists := sf.CachedUniversalIndexes[testKey]
		if !exists {
			t.Error("인덱스를 찾을 수 없음")
		}

		// 길이 비교
		if index2.Length <= initialLength {
			t.Error("새 데이터의 길이가 더 커야 함")
		}

		// 실제 데이터 확인
		data1, _ := json.Marshal(initialData)
		data2, _ := json.Marshal(largerData)
		if len(data2) <= len(data1) {
			t.Error("직렬화된 데이터의 길이도 더 커야 함")
		}
	})

	t.Run("동일 키에 대한 덮어쓰기 테스트", func(t *testing.T) {
		sf := GetStatusFileInstance()
		sf.Reset()

		testKey := "testkey123456789"
		initialData := map[string]interface{}{"value": "initial"}
		updates := make(UpdateMap)
		updates[testKey] = Update{New: initialData}

		// 첫 번째 쓰기
		err := sf.WriteUniversal(updates)
		if err != nil {
			t.Fatal(err)
		}
		if err = sf.WriteTasks(); err != nil {
			t.Fatal(err)
		}
		if err = sf.Commit(); err != nil {
			t.Fatal(err)
		}

		// 원래 위치 저장
		originalIndex := sf.CachedUniversalIndexes[testKey]

		// 같은 크기의 데이터로 업데이트
		sameData := map[string]interface{}{"value": "updated"}
		updates[testKey] = Update{New: sameData}

		// 두 번째 쓰기
		err = sf.WriteUniversal(updates)
		if err != nil {
			t.Fatal(err)
		}
		if err = sf.WriteTasks(); err != nil {
			t.Fatal(err)
		}
		if err = sf.Commit(); err != nil {
			t.Fatal(err)
		}

		// 새 인덱스 확인
		newIndex := sf.CachedUniversalIndexes[testKey]

		// 같은 위치에 저장되었는지 확인
		if originalIndex.FileID != newIndex.FileID {
			t.Errorf("FileID가 일치하지 않음: got %v, want %v", newIndex.FileID, originalIndex.FileID)
		}
		if originalIndex.Seek != newIndex.Seek {
			t.Errorf("Seek이 일치하지 않음: got %v, want %v", newIndex.Seek, originalIndex.Seek)
		}
		if originalIndex.Length != newIndex.Length {
			t.Errorf("Length가 일치하지 않음: got %v, want %v", newIndex.Length, originalIndex.Length)
		}
	})
}
