import hashlib

def op_hash_many(vars):
    if not isinstance(vars, list):
        raise Exception("OpHashMany: input is not an array")

    # Join strings with commas
    result = ','.join(str(v) for v in vars if isinstance(v, str)) #+ ","

    # Calculate SHA256 hash and return as hex string
    print(result)
    hash_result = hashlib.sha256(result.encode()).hexdigest()

    # Log the operation
    print(f"OpHashMany input: {vars} string: {result} result: {hash_result}")

    return hash_result

# Test
# test_input = ["hello", "world", "test"]
test_input = ["blk-1", "4975" ]

print(op_hash_many(test_input))
