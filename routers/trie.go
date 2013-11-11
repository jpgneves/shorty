package routers

import "fmt"

type TrieNode interface {
	IsPattern() bool
	PatternName() *string
	Value() interface{}
	Lookup(r rune) *TrieNode
	Find(s string) *RouteMatch
	Insert(key string, value interface{})
}

type Trie struct {
	pattern_name *string
	value        interface{}
	children     map[rune]*Trie
}

func CreateTrie() *Trie {
	return &Trie{value: nil, children: make(map[rune]*Trie)}
}

func (t Trie) IsPattern() bool {
	return t.pattern_name != nil
}

func (t Trie) PatternName() *string {
	return t.pattern_name
}

func (t Trie) Lookup(r rune) *Trie {
	if child, ok := t.children[r]; ok {
		return child
	}
	return nil
}

func (t Trie) Find(s string) *RouteMatch {
	matched_patterns := make(map[string]string)
	var current_pattern *string
	var current_value []rune
	node := &t
	for _, r := range s {
		n := node.Lookup(r)
		if n != nil {
			node = n
			if current_pattern != nil {
				matched_patterns[*current_pattern] = string(current_value)
				current_pattern = nil
			}
		} else if current_pattern != nil && r != '/' {
			current_value = append(current_value, r)
		} else if current_pattern != nil && r == '/' {
			matched_patterns[*current_pattern] = string(current_value)
			current_pattern = nil
		} else {
			for key := range node.children {
				child := node.children[key]
				if child.IsPattern() {
					current_pattern = child.PatternName()
					current_value = []rune{r}
					node = child
					break
				}
			}
			if current_pattern != nil {
				continue
			}
			return &RouteMatch{nil, nil}
		}
	}
	if current_pattern != nil {
		matched_patterns[*current_pattern] = string(current_value)
	}
	return &RouteMatch{node.value, matched_patterns}
}

func (t Trie) Insert(key string, value interface{}) {
	in_pattern := false
	var last_rune rune
	var pattern_name []rune
	node := &t
	for _, r := range key {
		if r == ':' && !in_pattern {
			in_pattern = true
			pattern_name = []rune{}
		} else if r != '/' && in_pattern {
			pattern_name = append(pattern_name, r)
		} else if r == '/' && in_pattern {
			in_pattern = false
			new_node := CreateTrie()
			pattern_name_str := string(pattern_name)
			new_node.pattern_name = &pattern_name_str
			node.children[last_rune] = new_node
			node = new_node
		} else {
			n := node.Lookup(r)
			if n != nil {
				node = n
			} else {
				new_node := CreateTrie()
				node.children[r] = new_node
				node = new_node
			}
			last_rune = r
		}
	}
	node.value = value
}

func (t Trie) Remove(key string) {
	in_pattern := false
	node := &t
	for _, r := range key {
		if r == ':' && !in_pattern {
			in_pattern = true
		} else if r != '/' && in_pattern {
		} else if r == '/' && in_pattern {
			in_pattern = false
		} else {
			n := node.Lookup(r)
			if n != nil {
				node = n
			} else {
				return
			}
		}
	}
	node.value = nil
}

func (t Trie) Print() {
	for key, _ := range t.children {
		if t.children[key].IsPattern() {
			fmt.Printf("{%+v} - %T %s\n", *(t.children[key].PatternName()), t.children[key].value, t.children[key].children)
		} else {
			fmt.Printf("%#U - %T %s\n", key, t.children[key].value, t.children[key].children)
		}
		t.children[key].Print()
	}
}
