package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

type trie interface {
	isLeaf() bool
	getChild(byte) trie
	setChild(byte, trie)
	getValue() string
	setLayer(int)
	insert(s string)
	toString() string
	toJSON() string
}

type arrayMappedTrie struct {
	children [26]trie
	layer    int
}

func (a *arrayMappedTrie) isLeaf() bool {
	return false
}

func (a *arrayMappedTrie) getChild(n byte) trie {
	return a.children[n-'A']
}

func (a *arrayMappedTrie) setChild(n byte, t trie) {
	t.setLayer(a.layer + 1)
	a.children[n-'A'] = t
}

func (a *arrayMappedTrie) getValue() string {
	return ""
}

func (a *arrayMappedTrie) setLayer(i int) {
	a.layer = i
}

func (a *arrayMappedTrie) insert(s string) {
	index := s[a.layer]
	child := a.getChild(index)
	if child == nil {
		a.setChild(index, leaf{value: s})
	} else if child.isLeaf() {
		a.setChild(index, &arrayMappedTrie{})
		a.getChild(index).insert(child.getValue())
		a.getChild(index).insert(s)
	} else {
		child.insert(s)
	}
}

func (a *arrayMappedTrie) toString() string {
	s := "[\n"
	for _, v := range a.children {
		if v != nil {
			s += v.toString()
		}
	}
	s += "]\n"
	return s
}

func (a *arrayMappedTrie) toJSON() string {
	s := "["
	for _, v := range a.children {
		if v != nil {
			s += v.toJSON() + ","
		}
	}
	return s[:len(s)-1] + "]"
}

type leaf struct {
	value string
	layer int
}

func (l leaf) isLeaf() bool {
	return true
}

func (l leaf) getChild(n byte) trie {
	return nil
}

func (l leaf) setChild(n byte, t trie) {}

func (l leaf) getValue() string {
	return l.value
}

func (l leaf) setLayer(i int) {
	l.layer = i
}

func (l leaf) insert(s string) {}

func (l leaf) toString() string {
	return l.value + "\n"
}

func (l leaf) toJSON() string {
	return l.value
}

func generateSubjectTrie(db *sql.DB) (trie, error) {
	var t trie = &arrayMappedTrie{}
	subjects, err := db.Query("SELECT DISTINCT c.SubjCode FROM courses c")
	if err != nil {
		return nil, err
	}
	defer subjects.Close()
	for subjects.Next() {
		var subject string
		subjects.Scan(&subject)
		if len(subject) != 4 {
			return nil, errors.New("Invalid Subject Code found")
		}
		t.insert(subject)
	}
	return t, nil
}

type trieJSON struct {
	Children []trieJSON `json:"children"`
	Value    string     `json:"value"`
}

func insertMap(x1 *map[string]map[string]map[string]map[string]string, s string) {
	_, ok := (*x1)[string(s[0])]
	if !ok {
		x4 := make(map[string]string)
		x3 := make(map[string]map[string]string)
		x2 := make(map[string]map[string]map[string]string)
		x4[string(s[3])] = s
		x3[string(s[2])] = x4
		x2[string(s[1])] = x3
		(*x1)[string(s[0])] = x2
		return
	}
	x2 := (*x1)[string(s[0])]
	_, ok = x2[string(s[1])]
	if !ok {
		x4 := make(map[string]string)
		x3 := make(map[string]map[string]string)
		x4[string(s[3])] = s
		x3[string(s[2])] = x4
		x2[string(s[1])] = x3
		return
	}
	x3 := x2[string(s[1])]
	_, ok = x3[string(s[2])]
	if !ok {
		x4 := make(map[string]string)
		x4[string(s[3])] = s
		x3[string(s[2])] = x4
		return
	}
	x4 := x3[string(s[2])]
	_, ok = x4[string(s[3])]
	if !ok {
		x4[string(s[3])] = s
		return
	}
}

func handleTrie(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]map[string]map[string]map[string]string)
	subjects, err := db.Query("SELECT DISTINCT c.SubjCode FROM courses c")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer subjects.Close()
	for subjects.Next() {
		var subject string
		subjects.Scan(&subject)
		if len(subject) != 4 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		insertMap(&m, subject)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}
