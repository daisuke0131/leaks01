package main

import (
	"go.uber.org/goleak"
	"testing"
)

func Test(t *testing.T) {
	defer goleak.VerifyNone(t)
	main()
}
