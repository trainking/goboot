package errgroup

import (
	"context"
	"fmt"
	"testing"
)

func TestForIndexGroup(t *testing.T) {
	group := WithContextForIndexGroup(context.Background())

	for i := 0; i < 10; i++ {
		group.Go(func(ctx context.Context, i int) error {
			fmt.Println(i)
			return nil
		}, i)
	}

	if err := group.Wait(); err != nil {
		t.Error(err)
	}
}
