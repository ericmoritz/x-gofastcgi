package gofastcgi

import (
	"time"
	)
type TimeIt struct {
	start int64;
}

func (t *TimeIt) Start() {
	t.start = time.Nanoseconds();
}

func (t *TimeIt) End() int64{
	return time.Nanoseconds() - t.start;
}