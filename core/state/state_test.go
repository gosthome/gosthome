package state

import "github.com/gosthome/gosthome/core/entity"

type test struct {
	A int
}

var _ entity.WithState[test] = (*State_[test])(nil)
