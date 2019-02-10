#!/usr/bin/env python3
# -*- coding: utf-8 -*-

# See http://primers.xyz/3

import array
import sys
import pdb
import collections
import itertools
import time
import random

from functools import reduce
from operator import mul    # or mul=lambda x,y:x*y
from fractions import Fraction

def nCk(n,k):
    return int( reduce(mul, (Fraction(n-i, i+1) for i in range(k)), 1) )


DEBUG = 2

finput = sys.stdin
if len(sys.argv) >= 2:
    finput = open(sys.argv[1], "rt")

def readline():
    return finput.readline().rstrip('\n')

def readints():
    return [int(x) for x in readline().split()]

def debug(*args):
    if DEBUG >= 1:
        print(*args, file=sys.stderr)

def dump(matrix, printer=str, file=sys.stderr):
    if DEBUG >= 2:
        for row in matrix:
            print("".join(map(printer, row)), file=file)

please_exit = False

import signal
import sys
def signal_handler(sig, frame):
    global please_exit
    please_exit = True

signal.signal(signal.SIGINT, signal_handler)
print('Press Ctrl+C')
#signal.pause()


class Pokemon:

    def __init__(self, name, attacks):
        self.name = name
        self.attacks = list(attacks)

    def setbits(self, bits, offset):
        for i in range(3):
            bits[offset + i] = 0
        for a in self.attacks:
            bits[offset + a//64] |= 1 << (a % 64)

    def __repr__(self):
        return "%s %r" % (self.name, self.attacks)

def score(pokemons, indices):
    s = set()
    for i in indices:
        p = pokemons[i]
        s.update(p.attacks)
    return len(s)

#http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetTable
#POPCOUNT_TABLE16 = [0] * (1 << 16)
POPCOUNT_TABLE16 = array.array('B', (0 for _ in range(1 << 16)))
for index in range(len(POPCOUNT_TABLE16)):
    POPCOUNT_TABLE16[index] = (index & 1) + POPCOUNT_TABLE16[index >> 1]

def popcount32(v):
    return (POPCOUNT_TABLE16[ v        & 0xffff] +
            POPCOUNT_TABLE16[(v >> 16) & 0xffff])

def popcount64(v):
    return (POPCOUNT_TABLE16[ v        & 0xffff] +
            POPCOUNT_TABLE16[(v >> 16) & 0xffff] +
            POPCOUNT_TABLE16[(v >> 32) & 0xffff] +
            POPCOUNT_TABLE16[(v >> 48) & 0xffff])

def _score(bits, indices):
    s0 = 0
    s1 = 0
    s2 = 0
    for i in indices:
        s0 |= bits[i*3+0]
        s1 |= bits[i*3+1]
        s2 |= bits[i*3+2]
    return popcount64(s0) + popcount64(s1) + popcount64(s2)

def main():
    pokemons = []
    while True:
        line = finput.readline()
        if not line:
            break
        items = line.split(sep=",")
        pokemons.append( Pokemon( name=items[0], attacks=list(map(int, items[1:]) )) )

    #bits = [0] * (3*len(pokemons))
    bits = array.array('Q', (0 for _ in range(3*len(pokemons))))
    for i in range(len(pokemons)):
        pokemons[i].setbits(bits, i*3)

    #for p in pokemons:
    #    debug(p)

    n = 6
    total = nCk( len(pokemons), n)
    debug("total", total)

    start = time.time()

    count = 0
    infocount = 0
    infotime = start

    best = 0
    #answer = list(range(n))
    answer = array.array('B', range(n))
    choosen = [False] * len(pokemons)
    for i in answer:
        choosen[i] = True

    previous = None

    while not please_exit:

        infocount += 1
        if infocount >= 1000:
            infocount = 0
            now = time.time()
            if now - infotime > 1.0:
                #debug("done %d %%" % round(100*float(count)/float(total)) )
                speed = float(count) / (now - infotime)
                debug("speed: %.1f Ksims per second" % (speed*1e-3) )
                count = 0
                infotime = now

        value = random.randrange(0, len(pokemons) )
        if choosen[value]:
            continue

        pos = random.randrange(0, n)
        previous = answer[pos]
        choosen[previous] = False
        choosen[value] = True
        answer[pos] = value

        #debug(team)
        s = _score(bits, answer)
        #assert( s == score(pokemons, answer) )
        #if s != score(pokemons, answer):
        #    print("s", s, "score", score(pokemons, answer))
        #    pdb.set_trace()

        if s > best:
            best = s
            debug("best", best)
        else:
            choosen[value] = False
            choosen[previous] = True
            answer[pos] = previous

        count += 1

    #
    debug("best", best)
    for i in answer:
        print(pokemons[i].name)


####

if __name__ == '__main__':
    main()