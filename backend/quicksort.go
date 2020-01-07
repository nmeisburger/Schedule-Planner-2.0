package main

import (
	"fmt"
	"math/rand"
	"time"
)

type pair struct {
	Distance int
	Course   *unParsedCourse
}

func quickSort(seq *[]pair) {
	L := len(*seq)
	if L > 1 {
		temp := make([]pair, L)
		for i := 0; i < L; i++ {
			temp[i] = (*seq)[i]
		}
		r := (*seq)[rand.Intn(L)]
		l := 0
		g := L - 1
		for _, t := range temp {
			if t.Distance < r.Distance {
				(*seq)[l] = t
				l++
			} else if t.Distance > r.Distance {
				(*seq)[g] = t
				g--
			}
		}
		same := l
		for _, t := range temp {
			if t.Distance == r.Distance {
				(*seq)[same] = t
				same++
			}
		}
		// for i := l; i <= g; i++ {
		// 	(*seq)[i] = r
		// }
		g++
		if g < L {
			right := (*seq)[g:]
			quickSort(&right)
		}
		if l > 0 {
			left := (*seq)[:l]
			quickSort(&left)
		}
	}
}

func testQuickSort(min, max int) {
	fails := 0
	start := time.Now()
	for i := min; i <= max; i++ {
		s := make([]pair, i, i)
		for j := 0; j < i; j++ {
			s[j] = pair{Distance: rand.Intn(100)}
		}
		quickSort(&s)
		for n := 0; n < i-1; n++ {
			if s[n].Distance > s[n+1].Distance {
				fmt.Println("Failure", i)
				fails++
				break
			}
		}
	}
	end := time.Now()
	fmt.Printf("Tests Complete: %d/%d Passed - Time: %v\n", max-fails-min, max-min, end.Sub(start))
}
