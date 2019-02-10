#!/usr/bin/env python3
# -*- coding: utf-8 -*-

# See http://primers.xyz/3

import sys
import pdb
import collections
import itertools
import time

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


class Pokemon:

    def __init__(self, name, attacks):
        self.name = name
        self.attacks = attacks

    def __repr__(self):
        return "%s %r" % (self.name, self.attacks)

def score(team):
    s = set()
    for p in team:
        s.update(p.attacks)
    return len(s)

def main():
    pokemons = []
    while True:
        line = finput.readline()
        if not line:
            break
        items = line.split(sep=",")
        pokemons.append( Pokemon( name=items[0], attacks=set(map(int, items[1:]) )) )

    #for p in pokemons:
    #    debug(p)

    n = 3
    total = nCk( len(pokemons), n)
    debug("total", total)

    start = time.time()
    last = start

    count = 0
    best = 0
    for team in itertools.combinations(pokemons, n):

        now = time.time()
        if now - last > 1.0:
            last = now
            debug("done %d %%" % round(100*float(count)/float(total)) )

        #debug(team)
        s = score(team)
        if s > best:
            best = s
            debug(best)

        count += 1

    #
    print(best)


####

if __name__ == '__main__':
    main()