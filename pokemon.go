package main

import (
	"bufio"
	"fmt"
	"log"
	"math/bits"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const M = 150
const N = 6

var names []string
var attacks [M][3]uint64

type Result struct {
	best   int
	answer [N]uint8
}

func (result *Result) search_i0(i0 int) {

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

						if score > result.best {
							result.best = score
							result.answer[0] = uint8(i0)
							result.answer[1] = uint8(i1)
							result.answer[2] = uint8(i2)
							result.answer[3] = uint8(i3)
							result.answer[4] = uint8(i4)
							result.answer[5] = uint8(i5)
						}

					}

				}
			}
		}
	}
}

func worker(w int, jobs <-chan int, results chan<- Result) {
	for i0 := range jobs {
		result := Result{}

		log.Printf("[w %d][i0 %d] processing\n", w, i0)
		start := time.Now()

		result.search_i0(i0)

		end := time.Now()
		elapsed := end.Sub(start)

		log.Printf("[w %d][i0 %d] elapsed = %.3f s\n", w, i0, elapsed.Seconds())
		log.Printf("[w %d][i0 %d] best = %d\n", w, i0, result.best)

		results <- result
	}
}

func search() Result {

	jobs := make(chan int, M)
	results := make(chan Result, M)

	for i0 := M - 1; i0 >= 0; i0-- {
		jobs <- i0
	}
	close(jobs)

	//numWorkers := 2
	numWorkers := runtime.GOMAXPROCS(0)
	for w := 0; w < numWorkers; w++ {
		go worker(w, jobs, results)
	}

	current := Result{}
	//for result := range results {
	for i0 := M - 1; i0 >= 0; i0-- {
		result := <-results
		if result.best > current.best {
			current.best = result.best
			for i := 0; i < N; i++ {
				current.answer[i] = result.answer[i]
			}
		}
		log.Printf("[current] best = %d\n", current.best)
	}

	return current
}

// go run pokemon.go input_3.txt
// go tool compile -S pokemon.go > pokemon_go.s
// go build -o pokemon_go.exe pokemon.go
// ./pokemon_go.exe input_3.txt
func main() {

	if !(len(os.Args) >= 2) {
		log.Fatal("usage: pokemon input.txt")
	}

	log.Printf("NumCPU= %d", runtime.NumCPU())
	//runtime.GOMAXPROCS(2)
	log.Printf("GOMAXPROCS= %d", runtime.GOMAXPROCS(0))

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
		s := strings.Split(strings.TrimSpace(s.Text()), ",")
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

	start := time.Now()

	result := search()

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("elapsed = %.3f s\n", elapsed.Seconds())
	log.Printf("best = %d\n", result.best)

	for i := 0; i < N; i++ {
		fmt.Println(names[result.answer[i]])
	}

}
