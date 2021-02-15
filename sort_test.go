package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sortBundles(t *testing.T) {
	t.Run("dependency bundle added to list in correct order", func(t *testing.T) {
		bundles := []Bundle{third}
		result, valid := sortBundles(bundles)
		require.True(t, valid)
		require.Len(t, result, 4)
		require.Equal(t, []Bundle{first, second, four, third}, result)
	})

	t.Run("duplicate bundles filtered correctly", func(t *testing.T) {
		bundles := []Bundle{first, second, four}
		result, valid := sortBundles(bundles)
		require.True(t, valid)
		require.Len(t, result, 3)
		require.Equal(t, []Bundle{first, second, four}, result)
	})

	t.Run("chaos check", func(t *testing.T) {
		bundles := []Bundle{first, second, third, first, second, third, first, second, first, second}
		result, valid := sortBundles(bundles)
		require.True(t, valid)
		require.Len(t, result, 4)
		require.Equal(t, []Bundle{first, second, four, third}, result)
	})
}

var (
	first = Bundle{
		Name: "first-bundle",
	}
	second = Bundle{
		Name: "second-bundle",
		Bundles: []Bundle{
			first,
		},
	}
	third = Bundle{
		Name: "third-bundle",
		Bundles: []Bundle{
			second,
			four,
		},
	}
	four = Bundle{
		Name: "four-bundle",
		Bundles: []Bundle{
			first,
		},
	}
)
