package common

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

/*
10458
10470
10551
10473

20011
20127
20103
20039
*/
func TestOpened(t *testing.T) {

	var total uint64 = 0

	background := context.Background()
	ctx, cancelFunc := context.WithCancel(background)
	go func(can context.CancelFunc) {
		timer := time.NewTimer(10 * time.Second)
		<-timer.C
		t.Logf("speed: %dkB/s", total/1024/10)
		can()
	}(cancelFunc)
	for Opened(ctx) {
		total += 1
		fmt.Sprintf("%d", total)
	}
}

/*
10933
10936
10907
10854

20845
20778
20795
20869
*/
func TestNonOpened(t *testing.T) {
	var total uint64 = 0

	go func(can func(code int), total *uint64) {
		timer := time.NewTimer(10 * time.Second)
		<-timer.C
		t.Logf("speed: %dkB/s", *total/1024/10)
		can(1)
	}(os.Exit, &total)
	for {
		total += 1
		fmt.Sprintf("%d", total)
	}
}
