package main

import (
	"hello/pkg/core/structure"
	"testing"
)

func TestParseOrderedMap_ComplexTransaction(t *testing.T) {
	// Test JSON string
	jsonStr := `{
        "transaction": {
            "type": "Send",
            "to": "9c39512b205792e5ab37f2612a3b64ed2ff394fa68ee",
            "amount": "10000",
            "from": "c68c21083d3000ed78476e1c0ca8b621b8f5e6a4b4f6",
            "timestamp": 1733633624081000
        },
        "public_key": "090d8f39b7577bf43fb1a12ff90c15c8b0837a81e8de9d6051f4a63f4cf8b51a",
        "signature": "28dbc58937756cc7caec02e7978a6e08b9f47a23282d0c4d84a1a0adbabfaf8afdee340b4ad08b746b5f3464709db169e75e8e8e00b1a4ed30d4b5a9cb1fe806"
    }`

	// Parse into OrderedMap
	om, err := structure.ParseOrderedMap(jsonStr)
	if err != nil {
		t.Fatalf("Error during parsing: %v", err)
	}

	// Validate top level keys
	expectedTopLevelKeys := []string{"transaction", "public_key", "signature"}
	keys := om.Keys()
	if len(keys) != len(expectedTopLevelKeys) {
		t.Errorf("Number of top level keys does not match. Expected: %d, Got: %d", len(expectedTopLevelKeys), len(keys))
	}

	// Validate transaction object
	transaction, ok := om.Get("transaction")
	if !ok {
		t.Fatal("Could not find transaction key")
	}
	t.Logf("transaction: %v", transaction)
	transactionMap, ok := transaction.(*structure.OrderedMap)
	if !ok {
		t.Fatal("Transaction value is not OrderedMap type")
	}

	// Validate transaction fields
	expectedFields := map[string]interface{}{
		"type":      "Send",
		"to":        "9c39512b205792e5ab37f2612a3b64ed2ff394fa68ee",
		"amount":    "10000",
		"from":      "c68c21083d3000ed78476e1c0ca8b621b8f5e6a4b4f6",
		"timestamp": int64(1733633624081000),
	}

	for field, expectedValue := range expectedFields {
		value, exists := transactionMap.Get(field)
		if !exists {
			t.Errorf("Could not find field: %s", field)
			continue
		}
		if value != expectedValue {
			t.Errorf("Field %s value does not match. Expected: %v, Got: %v", field, expectedValue, value)
		}
	}

	// Validate public_key and signature
	publicKey, ok := om.Get("public_key")
	if !ok || publicKey != "090d8f39b7577bf43fb1a12ff90c15c8b0837a81e8de9d6051f4a63f4cf8b51a" {
		t.Error("public_key value does not match")
	}

	signature, ok := om.Get("signature")
	if !ok || signature != "28dbc58937756cc7caec02e7978a6e08b9f47a23282d0c4d84a1a0adbabfaf8afdee340b4ad08b746b5f3464709db169e75e8e8e00b1a4ed30d4b5a9cb1fe806" {
		t.Error("signature value does not match")
	}
}

func TestParseOrderedMap_FlatData(t *testing.T) {
	// Flat JSON data without nesting
	jsonStr := `{
		"name": "홍길동",
		"age": 30,
		"email": "hong@example.com",
		"active": true
	}`

	om, err := structure.ParseOrderedMap(jsonStr)
	if err != nil {
		t.Fatalf("Error during parsing: %v", err)
	}

	// Validate expected fields and values
	expectedFields := map[string]interface{}{
		"name":   "홍길동",
		"age":    int64(30),
		"email":  "hong@example.com",
		"active": true,
	}

	// Validate key count
	if len(om.Keys()) != len(expectedFields) {
		t.Errorf("Key count does not match. Expected: %d, Got: %d", len(expectedFields), len(om.Keys()))
	}

	// Validate each field value
	for field, expectedValue := range expectedFields {
		value, exists := om.Get(field)
		if !exists {
			t.Errorf("Could not find field: %s", field)
			continue
		}
		if value != expectedValue {
			t.Errorf("Field %s value does not match. Expected: %v, Got: %v", field, expectedValue, value)
		}
	}
}

func TestParseOrderedMap_TripleNestedData(t *testing.T) {
	// Triple nested JSON data
	jsonStr := `{
		"user": {
			"profile": {
				"address": {
					"city": "서울",
					"district": "강남구",
					"street": "테헤란로",
					"zipcode": "06234"
				},
				"name": "김철수",
				"age": 25
			},
			"settings": {
				"notification": true,
				"theme": "dark"
			}
		},
		"created_at": "2024-03-20"
	}`

	om, err := structure.ParseOrderedMap(jsonStr)
	if err != nil {
		t.Fatalf("Error during parsing: %v", err)
	}

	// Validate user object
	user, ok := om.Get("user")
	if !ok {
		t.Fatal("Could not find user key")
	}
	userMap, ok := user.(*structure.OrderedMap)
	if !ok {
		t.Fatal("User value is not OrderedMap type")
	}

	// Validate profile object
	profile, ok := userMap.Get("profile")
	if !ok {
		t.Fatal("Could not find profile key")
	}
	profileMap, ok := profile.(*structure.OrderedMap)
	if !ok {
		t.Fatal("Profile value is not OrderedMap type")
	}

	// Validate address object
	address, ok := profileMap.Get("address")
	if !ok {
		t.Fatal("Could not find address key")
	}
	addressMap, ok := address.(*structure.OrderedMap)
	if !ok {
		t.Fatal("Address value is not OrderedMap type")
	}

	// Validate address fields
	expectedAddressFields := map[string]interface{}{
		"city":     "서울",
		"district": "강남구",
		"street":   "테헤란로",
		"zipcode":  "06234",
	}

	for field, expectedValue := range expectedAddressFields {
		value, exists := addressMap.Get(field)
		if !exists {
			t.Errorf("Could not find address field: %s", field)
			continue
		}
		if value != expectedValue {
			t.Errorf("Address field %s value does not match. Expected: %v, Got: %v", field, expectedValue, value)
		}
	}

	// Validate created_at
	createdAt, ok := om.Get("created_at")
	if !ok || createdAt != "2024-03-20" {
		t.Error("created_at value does not match")
	}
}
