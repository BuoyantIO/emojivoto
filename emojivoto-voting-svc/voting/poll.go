package voting

import (
	"log"
	"sort"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	counter *prometheus.CounterVec
}

func (p *inMemoryPoll) Vote(choice string) error {
	p.Lock()
	defer p.Unlock()

	if p.votes[choice] > 0 {
		p.votes[choice] = p.votes[choice] + 1
	} else {
		p.votes[choice] = 1
	}
	p.counter.With(prometheus.Labels{"emoji": choice}).Inc()
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

var counter *prometheus.CounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "emojivoto_votes_total",
	Help: "Number of emoji votes",
}, []string{"emoji"})

func NewPoll() Poll {
	poll := &inMemoryPoll{
		votes:   make(map[string]int, 0),
		counter: counter,
	}
	return poll
}
