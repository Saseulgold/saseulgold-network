import os

STATUS_KEY_BYTES = 64
STATUS_HEAP_BYTES = STATUS_KEY_BYTES + 10

def bin2hex(bytes):
    return bytes.hex()

def substr(str, offset, length):
    return str[offset: offset+length]

def bindec(bin_str: str) -> int:
    return int(bin2hex(bin_str), 16)

class StatIndexReader:

    def __init__(self):
        self.idx = 0

        with open("../data/status_bundle/lbundle_index", "rb") as f:
            self.data = f.read()

    @staticmethod
    def readPart(cls, file_id, seek, length):

        length = int(length, 16)
        seek = int(seek, 16)

        with open(f"../data/main_chain/{file_id.zfill(4)}", "rb") as f:
            raw = f.read()
            return raw[seek: seek+length]

    def reset(self):
        self.idx = 0

    def at(self, idx):
        raw = self.data[idx * STATUS_HEAP_BYTES: (idx + 1) * STATUS_HEAP_BYTES]
        return self.parseRaw(raw)

    @classmethod
    def parseRaw(cls, raw):
        k = raw[: STATUS_KEY_BYTES]
        return k

if __name__ == "__main__":
    reader = StatIndexReader()

    for j in range(3):
        raw = reader.at(j)
        k = reader.parseRaw(raw)
        #data = reader.readPart(*k)
        print(k)
        break
