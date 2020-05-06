package paxos

import (
	"testing"
	"time"
)

func TestPaxosWithSingleProposer(t *testing.T) {
	// 1, 2, 3 are acceptors
	// 1001 is a proposer
	pn := newPaxosNetwork(1, 2, 3, 1001, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	p := newProposer(1001, "hello world", pn.agentNetwork(1001), 1, 2, 3)
	go p.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()
	if value != "hello world" {
		t.Errorf("value = %s, want %s", value, "hello world")
	}
}

func TestPaxosWithTwoProposers(t *testing.T) {
	// 1, 2, 3 are acceptors
	// 1001,1002 is a proposer
	pn := newPaxosNetwork(1, 2, 3, 1001, 1002, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	p1 := newProposer(1001, "hello world", pn.agentNetwork(1001), 1, 2, 3)
	go p1.run()

	time.Sleep(time.Millisecond)
	p2 := newProposer(1002, "bad day", pn.agentNetwork(1002), 1, 2, 3)
	go p2.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()
	if value != "hello world" {
		t.Errorf("value = %s, want %s", value, "hello world")
	}
	time.Sleep(time.Millisecond)
}

func TestPaxosWithTwoProposers2(t *testing.T) {
	// 1, 2, 3 are acceptors
	// 1001,1002 is a proposer
	pn := newPaxosNetwork(1, 2, 3, 1001, 1002, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	p1 := newProposer(1002, "hello world", pn.agentNetwork(1002), 1, 2, 3)
	go p1.run()

	p2 := newProposer(1001, "bad day", pn.agentNetwork(1001), 1, 2, 3)
	go p2.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()
	if value != "hello world" {
		t.Errorf("value = %s, want %s", value, "hello world")
	}
	time.Sleep(time.Millisecond)
}

func TestPaxosWithTwoProposers3(t *testing.T) {
	// 1, 2, 3 are acceptors
	// 1001,1002 is a proposer
	pn := newPaxosNetwork(1, 2, 3, 1001, 1002, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	p11 := newProposer(1002, "hello world", pn.agentNetwork(1002), 1, 2, 3)
	p12 := newProposer(1002, "hello world2", pn.agentNetwork(1002), 1, 2, 3)
	go p11.run()
	go p12.run()

	p2 := newProposer(1001, "bad day", pn.agentNetwork(1001), 1, 2, 3)
	go p2.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()
	// only one value will be chosen on greater server ID.
	if value == "bad day" {
		t.Errorf("value = %s, want %s", value, "hello world or hello world2")
	}
	time.Sleep(time.Millisecond)
}
