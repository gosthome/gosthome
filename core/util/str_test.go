package util_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gosthome/gosthome/core/util"
)

func TestTranslit(t *testing.T) {
	fmt.Fprint(os.Stderr, util.CleanString("Привет"))
}
