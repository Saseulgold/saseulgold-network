import hashlib
import json

def hash(input_string):
  encoded_string = input_string.encode()
  sha256_hash = hashlib.sha256()
  sha256_hash.update(encoded_string)

  return sha256_hash.hexdigest()

blockheight = 4
s_ts = 1730335275000000

previous_blockhash = "0625bb12653a8002ab8d24287436848994ae9a351900d945fa21b45e3ea768ae19c0453b41ca0e"

thashs = [
  "0625bb136578288b0330b1a4126decddb85701916fe4482199ee2e6b68277ab5ae6ab0fbfca3d2"
]

uhashs = [
  "12194c0ef66a96758afcf4e7ddd3a0b851bba110c7dd2ffff358cbabd725b3fc0000000000000000000000000000000000000000000001edcb2fc15a26f0f59088588efcdee13430452921b2e033ef4ba3aa960f8d4f",
  "290eed314ce4d91c387028c290936b5b261e06f05d871bad42dfdf7436e89e9c000000000000000000000000000000000000000000008679932ed95880b0b030b53835a555146a66de726113e28de5d8c9f2cc449cc1",
  "724d2935080d38850e49b74927eb0351146c9ee955731f4ef53f24366c5eb9b100000000000000000000000000000000000000000000c32018abfa1a558faa839941f795f5b1491e9aa73cc7c8a50a1ef0d1ef6e783f",
  "87abdca0d3d3be9f71516090a362e5e79546f3183b1793789902c2e5176f62ae0000000000000000000000000000000000000000000054a83d7aa0269f5a997b4dedc21c43c51ed30c54f5c21ae9cbe68c374a29fe91",
  "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770d4a8fd2ebb308370a689c3ef47cb83ba182683def3d4f4e9f2b9553749280c41fb7abea697f0135d0cf252166f8770f162c275e4b7f3b",
  "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770d50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf45d09916fcdb9655ef131b297c5df7f23971efe33bc88008f945081283f9242a",
  "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770d632f1558df90bfccd366eae2f61ce553fc58f682aefb5c9bffc69ef0bd3f9ba1b2820d29a675832f568367465826f59bf33f3230f7a6",
  "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770dc6cbb26499ef5d324adafdd27d5613f1b276d0f6706e5c9bffc69ef0bd3f9ba1b2820d29a675832f568367465826f59bf33f3230f7a6",
  "c5ca2cb405daf22453b559420907bb12d7fb34519ac55d81f47829054374512f4a8fd2ebb308370a689c3ef47cb83ba182683def3d4fc115dc75e7e57deb1fbb596bd6d8de069b43e029bf89b5c9d7ebfc04cf715af9",
  "c8c603ff91a3c59d637c7bda83e732dea6ec74e1001b35600f0ba7831dbfe32900000000000000000000000000000000000000000000a509d3882cdfc3dfdce1d33aec769a69d103b4773959b97d1400938053d7917a",
  "fbab6eb9aa47eeb4d14b9473201b5aedbe0c03ba583be29f5840452ad2f1724200000000000000000000000000000000000000000000f39a06b8be71b6b599ba786972df8f9b7af5f86b5ced026cbe144666a9238989"
]


def merkelRoot(items):
    if len(items) == 0:
       return hash('')

    parent = list( hash(x) for x in items )
    while len(parent) > 1:
      child = []
      i = 0
      while i < len(parent):
        if len(parent) > i + 1:
          child.append(hash(parent[i] + parent[i + 1]))
        else:
          child.append(parent[i])
        i += 2
      parent = child

    return parent[0]

thashs = sorted(thashs)
uhashs = sorted(uhashs)

uh = merkelRoot(uhashs)
th = merkelRoot(thashs)

print("transactionRoot: ", th)
print("updateRoot: ", uh)

blockRoot = hash(th+uh)
print("BlockRoot: ", blockRoot)

i = {
  "height": blockheight,
  "s_timestamp": s_ts,
  "block_root": blockRoot
}

#i = str(i).replace(" ", "")
i = json.dumps(i, separators=(",", ":"))
print(i)
header = hash(i)
print("BlockHeader: ", header)

body = hash(f"{previous_blockhash}{header}")
print("BlockHash: ",  body)




















