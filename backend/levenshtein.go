package main

import "fmt"

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min3(a, b, c int) int {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	} else {
		return c
	}
}

func levenshteinDistance(query string, item string) int {
	queryLen, itemLen := len(query), len(item)

	M := make([][]int, queryLen+1, queryLen+1)

	for i := 0; i <= queryLen; i++ {
		M[i] = make([]int, itemLen+1, itemLen+1)
		M[i][0] = i
	}
	for j := 0; j <= itemLen; j++ {
		M[0][j] = j
	}

	for i := 1; i <= queryLen; i++ {
		for j := 1; j <= itemLen; j++ {
			var subCost int
			if query[i-1] == item[j-1] {
				subCost = 0
			} else {
				subCost = 1
			}
			M[i][j] = min3(M[i-1][j-1]+subCost, M[i-1][j]+1, M[i][j-1]+1)
		}
	}
	if queryLen < itemLen {
		return M[queryLen][itemLen] - (itemLen - queryLen)
	}
	return M[queryLen][itemLen]
}

func minLevenshteinDistance(query string, course *unParsedCourse) int {
	item := fmt.Sprintf("%s%d", course.SubjCode, course.CrseNum)
	d1 := levenshteinDistance(query, item)
	d2 := levenshteinDistance(query, course.Title)
	if d1 < d2 {
		return d1
	}
	return d2
}
