import os
import sys

class File:
    @staticmethod
    def read_part(file_path, start, length):
        with open(file_path, 'rb') as f:
            f.seek(start)
            return f.read(length)

class Parser:
    @staticmethod
    def bindec(binary_str):
        return int(binary_str, 2)
      
    @staticmethod
    def bin2hex(bytes):
        return binascii.hexlify(bytes)

class Protocol:

    CHAIN_KEY_BYTES = 39
    CHAIN_HEIGHT_BYTES = 4
    CHAIN_HEAP_BYTES = 4

    DATA_ID_BYTES = 4 
    SEEK_BYTES = 4
    LENGTH_BYTES = 4

    def chain_key(raw: bytes) -> str:
        # Convert CHAIN_KEY_BYTES as much as CHAIN_KEY_BYTES to hexadecimal from the first
        return raw[:Protocol.CHAIN_KEY_BYTES].hex()

    @staticmethod
    def chain_index(raw: bytes) -> list:
        offset = Protocol.CHAIN_KEY_BYTES

        # Convert data as much as CHAIN_HEIGHT_BYTES from binary to integer
        height = Parser.bindec(raw[offset:offset + Protocol.CHAIN_HEIGHT_BYTES].decode('utf-8'))
        offset += Protocol.CHAIN_HEIGHT_BYTES

        # Converting data as much as DATA_ID_BYTES into hexadecimal strings
        file_id = raw[offset:offset + Protocol.DATA_ID_BYTES].hex()
        offset += Protocol.DATA_ID_BYTES

        # Convert as much data as SEEK_BYTES from binary to integer
        seek = Parser.bindec(raw[offset:offset + Protocol.SEEK_BYTES].decode('utf-8'))
        offset += Protocol.SEEK_BYTES

        # Convert as much data as LENGTH_BYTES from binary to integer
        length = Parser.bindec(raw[offset:offset + Protocol.LENGTH_BYTES].decode('utf-8'))

        return [height, file_id, seek, length] 

    @classmethod
    def read_chain_indexes(cls, directory):
        indexes = {}

        index_file = os.path.join(directory, 'index')
        header = File.read_part(index_file, 0, Protocol.CHAIN_KEY_BYTES)
        count = Parser.bindec(header)
        print(count)
        exit(0)

        length = count * Protocol.CHAIN_HEAP_BYTES

        data = File.read_part(index_file, Protocol.CHAIN_HEADER_BYTES, length)
        data = [data[i:i + Protocol.CHAIN_HEAP_BYTES] for i in range(0, len(data), Protocol.CHAIN_HEAP_BYTES)]

        for raw in data:
            if len(raw) == Protocol.CHAIN_HEAP_BYTES:
                key = Protocol.chain_key(raw)
                index = Protocol.chain_index(raw)

                indexes[key] = index

        return indexes

if __name__ == "__main__":
    directory = sys.argv[1]
    Protocol.read_chain_indexes(directory)
