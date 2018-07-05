package voting

import (
	"log"
	"sort"
	"sync"
)

type Result struct {
	Shortcode string `json:"shortcode"`
	NumVotes  int    `json:"votes"`
}

type ByVotes []*Result

func (s ByVotes) Len() int      { return len(s) }
func (s ByVotes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByVotes) Less(i, j int) bool {
	return s[i].NumVotes > s[j].NumVotes
}

type Poll interface {
	Vote(choice string) error
	Results() ([]*Result, error)
}

type inMemoryPoll struct {
	votes map[string]int
	sync.RWMutex
}

func (p *inMemoryPoll) Vote(choice string) error {
	p.Lock()
	defer p.Unlock()

	if p.votes[choice] > 0 {
		p.votes[choice] = p.votes[choice] + 1
	} else {
		p.votes[choice] = 1
	}
	log.Printf("Voted for [%s], which now has a total of [%d] votes", choice, p.votes[choice])
	return nil
}

func (p *inMemoryPoll) Results() ([]*Result, error) {
	p.RLock()
	defer p.RUnlock()

	results := make([]*Result, 0)

	for emoji, numVotes := range p.votes {
		results = append(results, &Result{emoji, numVotes})
	}

	sort.Sort(ByVotes(results))

	return results, nil
}

func NewPoll() Poll {
	return &inMemoryPoll{
		votes: make(map[string]int, 0),
	}
}
