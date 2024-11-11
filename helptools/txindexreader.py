import os
import json


CHAIN_HEADER_BYTES = 4
CHAIN_HEIGHT_BYTES = 4
DATA_ID_BYTES = 2
SEEK_BYTES = 4
LENGTH_BYTES = 4

HASH_BYTES = 32
HEX_TIME_BYTES = 7


TIME_HASH_BYTES = HEX_TIME_BYTES + HASH_BYTES #39

CHAIN_KEY_BYTES = TIME_HASH_BYTES
TIME_HASH_SIZE = TIME_HASH_BYTES * 2

CHAIN_HEAP_BYTES = TIME_HASH_BYTES + 14
            
def bin2hex(bytes):
    return bytes.hex()

def substr(str, offset, length):
    return str[offset: offset+length]

def bindec(bin_str: str) -> int:
    return int(bin2hex(bin_str), 16)

import sys
DATA_PATH = sys.argv[1]

class MainChainIndexReader:

    def __init__(self):
        self.idx = 0

        with open(f"{DATA_PATH}/index", "rb") as f:
            self.raw = f.read()

        self.h = self.header()
        self.data = self.raw[CHAIN_HEADER_BYTES:]

    @staticmethod
    def readPart(cls, file_id, seek, length):
        length = int(length, 16)
        seek = int(seek, 16)

        with open(f"{DATA_PATH}/{file_id.zfill(4)}", "rb") as f:
            raw = f.read()
            return raw[seek: seek+length]

    def reset(self):
        self.idx = 0

    def part(self, seek, length):
        return self.data[seek: seek+length]
    
    def header(self):
        d = self.data[:CHAIN_HEADER_BYTES]
        print("headerraw: ", d)
        return bin2hex(d)

    @staticmethod
    def chain_key(raw: str):
        return int(substr(raw, 0, CHAIN_KEY_BYTES), 16)

    @staticmethod
    def chain_index(raw) -> list:
        offset = 0
        chainKey = substr(raw, offset, CHAIN_KEY_BYTES)
        chainKey = bin2hex(chainKey)
        offset += CHAIN_KEY_BYTES

        height = substr(raw, offset, CHAIN_HEIGHT_BYTES)
        height = bin2hex(height)
        offset += CHAIN_HEIGHT_BYTES

        file_id = substr(raw, offset, DATA_ID_BYTES)
        file_id = bin2hex(file_id)
        offset += DATA_ID_BYTES

        seek = substr(raw, offset, SEEK_BYTES)
        seek = bin2hex(seek)
        offset += SEEK_BYTES

        length = substr(raw, offset, LENGTH_BYTES)
        length = bin2hex(length)

        return chainKey, [ height, file_id, seek, length]

    def at(self, idx):
        raw = self.data[idx * CHAIN_HEAP_BYTES: (idx + 1) * CHAIN_HEAP_BYTES]
        return self.chain_index(raw)

    def header(self):
        header = self.raw[: CHAIN_HEADER_BYTES]
        count = int(bin2hex(header), 16)
        length = count * CHAIN_HEAP_BYTES
        return bin2hex(header), count, length

if __name__ == "__main__":
    reader = MainChainIndexReader()
    header, count, length = reader.header()

    print(header, count, length)

    j = 0
    while True:
        k, i = reader.at(j)
        print(k)
        j += 1
        data = reader.readPart(*i)
        data = json.loads(data.decode())


