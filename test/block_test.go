package main

import (
	_ "fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/core/storage"
	S "hello/pkg/core/structure"
	"testing"
)

func TestBlock_WithMultipleUpdates(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "genesis_test_2"

	// Create Transaction
	txData := S.NewOrderedMap()
	data := S.NewOrderedMap()
	txData.Set("type", "Submit")
	txData.Set("phrase", "062844d209aac040ae0d923501f7ef000567fc639ec0481d0e26f7678a273ccfc4e9d05d37f904")
	txData.Set("nonce", "9916046")
	txData.Set("calculated_hash", "000010adf9537cb677709b57aaf34b3ea19ca5b77d34180436f8a41c0680c03d")
	txData.Set("from", "570802432a9917544300c9a3db0becbab5539f4c54aa")
	txData.Set("timestamp", int64(1733129866477000))
	data.Set("transaction", txData)
	data.Set("public_key", "e3360fb1b094899e82b04e3e81a234af06bd4fee13fb66d6e1b08048bd569e52")
	data.Set("signature", "70d5e58e0618c95644523f3e7fbecba3c04eb63a303b5da6813119a89074e651616de8b9dd0cdbca1db114c8d2e18114410f158d2aacf8c07114d78b0a6cd00b")
	tx, err := NewSignedTransaction(data)
	if err != nil {
		t.Fatalf("Failed to create tx: %v", err)
	}

	// Create a block
	previousBlockhash := "062845bbcea340127a64fb4398a575f8d56ee511c92525b020df2ae24156154c568c2c459dc10e"
	block := NewBlock(1281776, previousBlockhash)
	block.SetTimestamp(1733129867000000)
	block.Difficulty = 3725
	block.RewardAddress = "570802432a9917544300c9a3db0becbab5539f4c54aa"
	block.Vout = "738036000000000000000"
	block.Nonce = "9916046"

	// Universal Updates added
	block.AppendUniversalUpdate(Update{
		Key: "fe38ef5ff626c7e8caeeba0eecb873d8652b03283aaf3bd1721ad58aa073b4b700000000000000000000000000000000000000000000",
		Old: "00000000000000000000000000000000000000000000",
		New: "00000000000000000000000000000000000000000000",
	})

	block.AppendUniversalUpdate(Update{
		Key: "60bca2cdd6712dc6316262cd7196a724ed2f47b59c7a4c6bac63113ddc4dc22300000000000000000000000000000000000000000000",
		Old: "0",
		New: "0",
	})

	block.AppendUniversalUpdate(Update{
		Key: "1b49a8b9ffaf283c1cb5ace27b3cb04ce4dfff10caba0c7228c2e3a54df54b1f00000000000000000000000000000000000000000000",
		Old: "1088403",
		New: "9916046",
	})

	block.AppendUniversalUpdate(Update{
		Key: "54604d2761063cbb07a16c4aa192b2bbab7d754fb7d294d4f3b3f642cb815e8c00000000000000000000000000000000000000000000",
		Old: "3921",
		New: "38",
	})

	block.AppendUniversalUpdate(Update{
		Key: "1efdbeb7a5fffb7e1505672c39778da4101466076bde006530949b0e909e287f00000000000000000000000000000000000000000000",
		Old: 1733129829124809,
		New: 1733129867207691,
	})

	block.AppendUniversalUpdate(Update{
		Key: "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770d570802432a9917544300c9a3db0becbab5539f4c54aa",
		Old: "119699008308270313659337618",
		New: "119699746344270313659337618",
	})

	block.AppendUniversalUpdate(Update{
		Key: "7bc990ac426a5a6c43cab692b4da3046859ad68a162c1eb8e73ed73b651ed7c600000000000000000000000000000000000000000000",
		Old: "062844d209aac040ae0d923501f7ef000567fc639ec0481d0e26f7678a273ccfc4e9d05d37f904",
		New: "062845bbcea340127a64fb4398a575f8d56ee511c92525b020df2ae24156154c568c2c459dc10e",
	})

	block.AppendUniversalUpdate(Update{
		Key: "87abdca0d3d3be9f71516090a362e5e79546f3183b1793789902c2e5176f62ae00000000000000000000000000000000000000000000",
		Old: "3725",
		New: "3790",
	})

	block.AppendUniversalUpdate(Update{
		Key: "fbab6eb9aa47eeb4d14b9473201b5aedbe0c03ba583be29f5840452ad2f1724200000000000000000000000000000000000000000000",
		Old: nil,
		New: "00114ab4c0e970882c4f6f2e56305cf18b8ce6bcdbee2af5990f43f3924e1558",
	})

	block.AppendUniversalUpdate(Update{
		Key: "c8c603ff91a3c59d637c7bda83e732dea6ec74e1001b35600f0ba7831dbfe32900000000000000000000000000000000000000000000",
		Old: "76153662000000000000000",
		New: "738036000000000000000",
	})

	// Local Updates added
	block.AppendLocalUpdate(Update{
		Key: "724d2935080d38850e49b74927eb0351146c9ee955731f4ef53f24366c5eb9b100000000000000000000000000000000000000000000",
		Old: 3561659,
		New: 3561660,
	})

	// Add Transaction
	block.AppendTransaction(tx)

	// Verification logic

	expectedTxRoot := "1790d8a9c5c07a12c5f3c8a7d18512a1dc1d5082345d097dcabeea126e36a12e"
	actualTxRoot := block.TransactionRoot()
	if actualTxRoot != expectedTxRoot {
		t.Errorf("TransactionRoot() = %v; want %v", actualTxRoot, expectedTxRoot)
	}

	expectedUpdateRoot := "eb8d6dc3ca89fc6654a97126f070d5fc5a3a05ddfbc61a62ed87a4adc583ffd3"
	actualUpdateRoot := block.UpdateRoot()
	if actualUpdateRoot != expectedUpdateRoot {
		t.Errorf("UpdateRoot() = %v; want %v", actualUpdateRoot, expectedUpdateRoot)
	}

	expectedBlockRoot := "4061165ff7383021921d8354c21aa43f30254486067539b003c8f6d2cc7c13d5"
	actualBlockRoot := block.BlockRoot()
	if actualBlockRoot != expectedBlockRoot {
		t.Errorf("BlockRoot() = %v; want %v", actualBlockRoot, expectedBlockRoot)
	}

	expectedBlockHeader := "20c60593fafff926da4a179ed654629c5ae892fe8de1d7f48840950bfffb5773"
	actualBlockHeader := block.BlockHeader()
	if actualBlockHeader != expectedBlockHeader {
		t.Errorf("BlockHeader() = %v; want %v", actualBlockHeader, expectedBlockHeader)
	}

	expectedBlockHash := "062845be1278c064126c0dfbbac31850a780508b8f22bab9e4e7d752941356edeed1f78b7ca555"
	actualBlockHash := block.BlockHash()
	if actualBlockHash != expectedBlockHash {
		t.Errorf("BlockHash() = %v; want %v", actualBlockHash, expectedBlockHash)
	}

	chain := GetChainStorageInstance()
	err = chain.Touch()
	if err != nil {
		t.Errorf("Error occurred during touching chain: %v", err)
	}

	err = chain.Write(&block)
	if err != nil {
		t.Errorf("Error occurred during writing block: %v", err)
	}
}
