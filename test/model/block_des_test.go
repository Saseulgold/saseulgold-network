package model

import (
	"hello/pkg/core/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockDeserialization(t *testing.T) {
	// 테스트할 블록 JSON 문자열
	blockJson := `{"height":100,"s_timestamp":1734240706687230,"previous_blockhash":"previous_hash_test","blockhash":"0629486146f0feeac9ca9f9df435954edafc958474f7721d07ee77fb7f6ae3f69ea10ab944f81f","difficulty":4,"reward_address":"reward_address_test","vout":"vout_test","nonce":"nonce_test","transactions":{"0629486146f0cdb9408d4a64f8faee013d5d196b9ef1706354677da8c1b269bc215473cbf586a0":{"transaction":{"type":"Send","to":"50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf","from":"60c3a6cd858c90574bcdc35b2da5dbc7225275f50efd","amount":"1000","timestamp":1734240706687181},"public_key":"test_public_key","signature":"test_signature"}},"universal_updates":{},"local_updates":{"old":null,"new":null}}`

	// 블록 파싱
	block, err := storage.ParseBlock([]byte(blockJson))
	assert.NoError(t, err)

	// 기본 필드 검증
	assert.Equal(t, 100, block.Height, "블록 높이 검증")
	assert.Equal(t, int64(1734240706687230), block.Timestamp_s, "타임스탬프 검증")
	assert.Equal(t, "previous_hash_test", block.PreviousBlockhash, "이전 블록 해시 검증")
	assert.Equal(t, 4, block.Difficulty, "난이도 검증")
	assert.Equal(t, "reward_address_test", block.RewardAddress, "보상 주소 검증")
	assert.Equal(t, "vout_test", block.Vout, "Vout 검증")
	assert.Equal(t, "nonce_test", block.Nonce, "Nonce 검증")
}
