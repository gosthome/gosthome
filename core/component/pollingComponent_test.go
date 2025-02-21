package component

import (
	"context"
	"errors"
)

type testPoll struct {
	pc *PollingComponent[testPoll, *testPoll]
	c  int
}

func newTestComponent(ctx context.Context, cfg *PollingComponentConfig[testPoll, *testPoll]) (ret *testPoll, err error) {
	ret = &testPoll{}
	ret.pc, err = NewPollingComponent(ctx, ret, cfg)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Setup implements Poller.
func (t *testPoll) Setup() {
	t.pc.Setup()
}

// Poll implements Poller.
func (t *testPoll) Poll() {
	t.c += 1
}

// Close implements Poller.
func (t *testPoll) Close() error {
	return errors.Join(t.pc.Close())
}

// InitializationPriority implements Poller.
func (t *testPoll) InitializationPriority() InitializationPriority {
	return InitializationPriorityProcessor
}

var _ Poller = (*testPoll)(nil)
var _ Component = (*PollingComponent[testPoll, *testPoll])(nil)
