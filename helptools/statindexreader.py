import os

STATUS_KEY_BYTES = 64
STATUS_HEAP_BYTES = STATUS_KEY_BYTES + 10

DATA_ID_BYTES = 2
SEEK_BYTES = 4
LENGTH_BYTES = 4

DATA_DIR = "/home/ec2-user/qcn/data/status_bundle"

def bin2hex(bytes):
    return bytes.hex()

def substr(str, offset, length):
    return str[offset: offset+length]

def bindec(bin_str: str) -> int:
    return int(bin2hex(bin_str), 16)

class StatIndexReader:

    def __init__(self):
        self.idx = 0

        with open( os.path.join(DATA_DIR, "ubundle_index"), "rb") as f:
            self.data = f.read()

    @staticmethod
    def readPart(cls, file_id, seek, length):

        length = int(length, 16)
        seek = int(seek, 16)

        with open( os.path.join(DATA_DIR, f"universals-{file_id}"), "rb") as f:
            raw = f.read()
            return raw[seek: seek+length]

    def reset(self):
        self.idx = 0

    def at(self, idx):
        raw = self.data[idx * STATUS_HEAP_BYTES: (idx + 1) * STATUS_HEAP_BYTES]
        return raw

    @classmethod
    def parseRaw(cls, raw):
        k = raw[: STATUS_KEY_BYTES]
        k = bin2hex(k)

        offset = STATUS_KEY_BYTES
        file_id = raw[offset: offset+DATA_ID_BYTES]

        offset += DATA_ID_BYTES
        seek = raw[offset:offset+SEEK_BYTES] 
        print(seek)
        seek = int.from_bytes(seek, byteorder='little')

        offset += SEEK_BYTES
        length = raw[offset: offset+LENGTH_BYTES]
        print(length)
        length = int.from_bytes(length, byteorder='little')
        return k, file_id, seek, length

if __name__ == "__main__":
    reader = StatIndexReader()
   
    j = 0

    while True:
        raw = reader.at(j)
        k, _, seek, length = reader.parseRaw(raw)
        #data = reader.readPart(*k)

        if k == "":
            break
        print(k, seek, length)
        j += 1
