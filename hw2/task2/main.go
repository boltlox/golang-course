package main

import (
	"fmt"
	"sort"
)

type Candidate struct {
	Name  string
	Votes int
}

func countVotes(names []string) []Candidate {

	voteCount := make(map[string]int)

	for _, name := range names {
		voteCount[name]++
	}

	var candidates []Candidate

	for name, votes := range voteCount {
		candidates = append(candidates, Candidate{Name: name, Votes: votes})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Votes > candidates[j].Votes
	})

	return candidates

}

func main() {
	votes := []string{"Ann", "Kate", "Peter", "Kate", "Ann", "Ann", "Helen"}
	result := countVotes(votes)
	fmt.Println(result)

}
