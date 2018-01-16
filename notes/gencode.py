#!/usr/bin/env python3
#

ABC = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'

def cycle(v): return ((v * 0x1481) + 0x3039) & 0x7fff

def gencode(mix, key):
   mix3 = list(mix)
   mix12 = []
   for i in range(4):
       mix3 = [cycle(v) for v in mix3]
       mix12.extend(mix3)
   return "".join([ABC[(v ^ key) % 52] for v in mix12])

key = 0x1092
mix = [0x4f26, 0x319b, 0x1c58]

if __name__ == '__main__':
    print("".join(gencode(mix, key)))
