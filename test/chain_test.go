package main

import (
	. "hello/pkg/core/storage"
	"os"
	"testing"
)

func TestChainStorageWriteAndRead(t *testing.T) {
	// 테스트 디렉토리 설정
	testDir := "test_chain_storage"
	defer os.RemoveAll(testDir)

	chain := &ChainStorage{}

	// 초기화
	err := chain.Touch(testDir)
	if err != nil {
		t.Fatalf("초기화 실패: %v", err)
	}

	// 테스트 데이터
	testCases := []struct {
		height int
		key    string
		data   []byte
	}{
		{1, "test1", []byte("hello world 1")},
		{2, "test2", []byte("hello world 2")},
		{3, "test3", []byte("hello world 3")},
	}

	// 데이터 쓰기 테스트
	for _, tc := range testCases {
		err := chain.WriteData(testDir, tc.height, tc.key, tc.data)
		if err != nil {
			t.Errorf("높이 %d 데이터 쓰기 실패: %v", tc.height, err)
		}
	}

	// 데이터 읽기 테스트
	for _, tc := range testCases {
		// 인덱스로 데이터 조회
		indices, err := chain.Index(testDir, tc.height)
		if err != nil {
			t.Errorf("높이 %d 인덱스 조회 실패: %v", tc.height, err)
			continue
		}

		// 데이터 읽기
		data, err := chain.ReadData(testDir, indices)
		if err != nil {
			t.Errorf("높이 %d 데이터 읽기 실패: %v", tc.height, err)
			continue
		}

		// 데이터 검증
		if string(data) != string(tc.data) {
			t.Errorf("높이 %d 데이터 불일치\n원본: %s\n읽은값: %s",
				tc.height, string(tc.data), string(data))
		}
	}

	// 마지막 인덱스 테스트
	lastIdx := chain.LastIdx(testDir)
	if lastIdx != len(testCases) {
		t.Errorf("마지막 인덱스 불일치. 예상: %d, 실제: %d",
			len(testCases), lastIdx)
	}

	// 키로 검색 테스트
	for _, tc := range testCases {
		indices, err := chain.Index(testDir, tc.key)
		if err != nil {
			t.Errorf("키 %s 검색 실패: %v", tc.key, err)
			continue
		}

		data, err := chain.ReadData(testDir, indices)
		if err != nil {
			t.Errorf("키 %s 데이터 읽기 실패: %v", tc.key, err)
			continue
		}

		if string(data) != string(tc.data) {
			t.Errorf("키 %s 데이터 불일치\n원본: %s\n읽은값: %s",
				tc.key, string(tc.data), string(data))
		}
	}
}
