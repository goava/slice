package slice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/goava/slice"
)

func TestContext(t *testing.T) {
	b := slice.NewContext(context.TODO())
	b.Set("key", "value")
	ctx := context.WithValue(context.Background(), "ctx-key", "ctx-value")
	joined := b.Join(ctx)
	require.Equal(t, joined.Value("key"), "value")
	require.Equal(t, joined.Value("ctx-key"), "ctx-value")
}
