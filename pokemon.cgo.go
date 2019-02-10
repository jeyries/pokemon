package main

import (
	"bufio"
	"fmt"
	"log"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

/*
#cgo CFLAGS: -O3 -march=native
#cgo LDFLAGS:

#include <stdint.h>

#define M 150
#define N 6

void search_i0( uint64_t attacks[M][3], uint8_t answer[N],  int64_t i0, int64_t* best) {

	for (int64_t i1 = 0; i1 < i0; i1++) {
		for (int64_t i2 = 0; i2 < i1; i2 ++) {
			for (int64_t i3 = 0; i3 < i2; i3 ++) {
				for (int64_t i4 = 0; i4 < i3; i4++) {
					for (int64_t i5 = 0; i5 < i4; i5++) {

						uint64_t s0 = 0;
						uint64_t s1 = 0;
						uint64_t s2 = 0;

						s0 |= attacks[i0][0];
						s1 |= attacks[i0][1];
						s2 |= attacks[i0][2];

						s0 |= attacks[i1][0];
						s1 |= attacks[i1][1];
						s2 |= attacks[i1][2];

						s0 |= attacks[i2][0];
						s1 |= attacks[i2][1];
						s2 |= attacks[i2][2];

						s0 |= attacks[i3][0];
						s1 |= attacks[i3][1];
						s2 |= attacks[i3][2];

						s0 |= attacks[i4][0];
						s1 |= attacks[i4][1];
						s2 |= attacks[i4][2];

						s0 |= attacks[i5][0];
						s1 |= attacks[i5][1];
						s2 |= attacks[i5][2];

						int64_t score = __builtin_popcountll(s0) + __builtin_popcountll(s1) + __builtin_popcountll(s2);
						if (score > *best) {
							*best = score;
							answer[0] = (uint8_t) i0;
							answer[1] = (uint8_t) i1;
							answer[2] = (uint8_t) i2;
							answer[3] = (uint8_t) i3;
							answer[4] = (uint8_t) i4;
							answer[5] = (uint8_t) i5;
						}
					}
				}
			}
		}
	}
}
*/
import "C"

const M = 150
const N = 6

var names []string
var attacks [M][3]uint64
var answer [N]uint8

func search_i0(i0 int, best *int, answer *[N]uint8) {

	var s4 [3]uint64

	for i1 := 0; i1 < i0; i1++ {
		for i2 := 0; i2 < i1; i2++ {
			for i3 := 0; i3 < i2; i3++ {
				for i4 := 0; i4 < i3; i4++ {

					s4[0] = 0
					s4[1] = 0
					s4[2] = 0

					s4[0] |= attacks[i0][0]
					s4[1] |= attacks[i0][1]
					s4[2] |= attacks[i0][2]

					s4[0] |= attacks[i1][0]
					s4[1] |= attacks[i1][1]
					s4[2] |= attacks[i1][2]

					s4[0] |= attacks[i2][0]
					s4[1] |= attacks[i2][1]
					s4[2] |= attacks[i2][2]

					s4[0] |= attacks[i3][0]
					s4[1] |= attacks[i3][1]
					s4[2] |= attacks[i3][2]

					s4[0] |= attacks[i4][0]
					s4[1] |= attacks[i4][1]
					s4[2] |= attacks[i4][2]

					for i5 := 0; i5 < i4; i5++ {

						s50 := s4[0]
						s51 := s4[1]
						s52 := s4[2]

						s50 |= attacks[i5][0]
						s51 |= attacks[i5][1]
						s52 |= attacks[i5][2]

						score := bits.OnesCount64(s50)
						score += bits.OnesCount64(s51)
						score += bits.OnesCount64(s52)

						if score > *best {
							*best = score
							answer[0] = uint8(i0)
							answer[1] = uint8(i1)
							answer[2] = uint8(i2)
							answer[3] = uint8(i3)
							answer[4] = uint8(i4)
							answer[5] = uint8(i5)
						}

					}

				}
			}
		}
	}
}

func search_i0_c(i0 int, best *int, answer *[N]uint8) {
	var _attacks = (*[3]C.uint64_t)(unsafe.Pointer(&attacks))
	var _i0 = C.int64_t(i0)
	var _best = (*C.int64_t)(unsafe.Pointer(best))
	var _answer = (*C.uint8_t)(unsafe.Pointer(answer))
	C.search_i0(_attacks, _answer, _i0, _best)
}

func search() {

	for i := 0; i < N; i++ {
		answer[i] = 0
	}

	best := 0

	for _i0 := 0; _i0 < M; _i0++ {
		i0 := M - 1 - _i0

		fmt.Printf("processing i0= %d\n", i0)
		start := time.Now()

		//search_i0(i0, &best, &answer)
		search_i0_c(i0, &best, &answer)

		end := time.Now()
		elapsed := end.Sub(start)
		fmt.Printf("elapsed = %.3f s\n", elapsed.Seconds())
		fmt.Printf("best = %d\n", best)
	}
}

// go run pokemon.go input_3.txt
// go tool compile -S pokemon.go > pokemon_go.s
// go build -o pokemon_go.exe pokemon.go
// ./pokemon_go.exe input_3.txt
func main() {

	if !(len(os.Args) >= 2) {
		log.Fatal("usage: pokemon input.txt")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	k := 0
	s := bufio.NewScanner(f)
	for s.Scan() {
		s := strings.Split(s.Text(), ",")
		names = append(names, s[0])
		for _, v := range s[1:] {
			x, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			attacks[k][x/64] |= uint64(1) << (uint(x) % 64)
		}
		k++
	}
	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %016x %016x %016x\n", names[0], attacks[0][0], attacks[0][1], attacks[0][2])

	search()

	for i := 0; i < len(answer); i++ {
		fmt.Println(names[answer[i]])
	}

}
