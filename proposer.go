package paxos

import (
	"log"
	"time"
)

type proposer struct {
	id int
	// stable
	lastSeq int

	proposalID int64

	value  string
	valueN int64

	acceptors map[int]promise
	nt        network
}

func newProposer(id int, value string, nt network, acceptors ...int) *proposer {
	p := &proposer{id: id, nt: nt, lastSeq: 0, value: value, acceptors: make(map[int]promise)}
	p.proposalID = p.n()
	for _, a := range acceptors {
		p.acceptors[a] = message{}
	}
	return p
}

func (p *proposer) run() {
	var ok bool
	var m message

	// stage 1: do prepare until reach the majority
	for !p.majorityReached() {
		if !ok {
			ms := p.prepare()
			for i := range ms {
				p.nt.send(ms[i])
			}
		}
		m, ok = p.nt.recv(time.Second)
		if !ok {
			// the previous prepare is failed
			// continue to do another prepare
			continue
		}

		switch m.typ {
		case Promise:
			p.receivePromise(m)
		default:
			log.Panicf("proposer: %d unexpected message type: %v", p.id, m.typ)
		}
	}
	log.Printf("proposer: %d promise %d reached majority %d", p.id, p.n(), p.majority())

	// stage 2: do propose
	log.Printf("proposer: %d starts to propose [%d: %s]", p.id, p.n(), p.value)
	ms := p.propose()
	for i := range ms {
		p.nt.send(ms[i])
	}
}

// If the proposer receives the requested responses from a majority of
// the acceptors, then it can issue a proposal with number n and value
// v, where v is the value of the highest-numbered proposal among the
// responses, or is any value selected by the proposer if the responders
// reported no proposals.
func (p *proposer) propose() []message {
	ms := make([]message, p.majority())

	i := 0
	for to, promise := range p.acceptors {
		if promise.number() == p.proposalID {
			ms[i] = message{from: p.id, to: to, typ: Propose, n: p.proposalID, value: p.value}
			i++
		}
		if i == p.majority() {
			break
		}
	}
	return ms
}

// A proposer chooses a new proposal number n and sends a request to
// each member of some set of acceptors, asking it to respond with:
// (a) A promise never again to accept a proposal numbered less than n, and
// (b) The proposal with the highest number less than n that it has accepted, if any.
func (p *proposer) prepare() []message {
	p.proposalID = p.n()
	p.lastSeq++
	ms := make([]message, p.majority())
	i := 0
	for to := range p.acceptors {
		ms[i] = message{from: p.id, to: to, typ: Prepare, n: p.proposalID}
		i++
		if i == p.majority() {
			break
		}
	}
	return ms
}

func (p *proposer) receivePromise(promise message) {
	prevPromise := p.acceptors[promise.from]

	if prevPromise.number() < promise.number() {
		log.Printf("proposer: %d received a new promise %+v", p.id, promise)
		p.acceptors[promise.from] = promise

		//update value to the value with a larger N
		if promise.proposalValue() != "" {
			log.Printf("proposer: %d updated the value [%s] to %s", p.id, p.value, promise.proposalValue())
			p.value = promise.proposalValue()
		}
	}
}

func (p *proposer) majority() int { return len(p.acceptors)/2 + 1 }

func (p *proposer) majorityReached() bool {
	m := 0
	for _, promise := range p.acceptors {
		if promise.number() == p.proposalID {
			m++
		}
	}
	if m >= p.majority() {
		return true
	}
	return false
}

func (p *proposer) n() int64 {
	begin := time.Date(2018, time.November, 17, 22, 0, 0, 0, time.UTC)
	duration := time.Now().Sub(begin)
	return (duration.Milliseconds())*10000 + int64(p.id)
}
