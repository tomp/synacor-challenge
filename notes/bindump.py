#!/usr/bin/env python3
#
import array
import struct

SOURCEFILE = "challenge.bin"
WIDTH = 16

buffer = array.array('H')

with open(SOURCEFILE, 'rb') as fp:
    try:
        buffer.fromfile(fp, 30050)
    except EOFError:
        pass
size = len(buffer)
buffer.fromlist([0]*WIDTH)

def char(c):
    if c >= 0x20 and c < 0x7f:
        return chr(c)
    else:
        return '.'

def write_row(addr, words):
    print("{:04x}: {}  {}".format(addr,
        " ".join(["{:04x}".format(c) for c in words]),
        "".join([char(c) for c in words])))

def decrypt_bulk(word, key1, key2=0x4154):
    return 0x7fff & (word ^ (key1 * key1) ^ key2)

def decrypt_string(addr, key):
    nchar = buffer[addr]
    print("decrypt {:04x}({:x}) w/ {:04x}".format(addr, nchar, key))
    for loc in range(addr+1,addr+nchar+1):
        buffer[loc] = buffer[loc] ^ key

for addr in range(0x17b4, 0x7562):
    buffer[addr] = decrypt_bulk(buffer[addr], addr)

ops = [0x0001, 0x8000, 0x0001, 0x8001, 0x0009, 0x8002, 0x0011]

cur = -1
try:
    step = buffer[cur+1:].index(0x05b2)
    while step > 0:
        cur += step + 1
        b = buffer[cur-11:cur+1].tolist()
        if b[:2] + b[3:5] + b[6:8] + [b[10]] == ops:
            addr, sub, key1, key2 = b[2], b[5], b[8], b[9]
            assert sub == 0x05fb
            key = (key1 + key2) & 0x7fff
            decrypt_string(addr, key)
        step = buffer[cur+1:].index(0x05b2)
except ValueError:
    pass

addr = 0
while addr < size:
    words = buffer[addr:addr+WIDTH]
    write_row(addr, words)
    addr += WIDTH

