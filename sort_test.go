package slice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type bundleNames []Bundle

func (p bundleNames) Names() []string {
	var names []string
	for _, b := range p {
		names = append(names, b.Name)
	}
	return names
}

func Test_sortBundles(t *testing.T) {
	t.Run("dependency bundle added to list in correct order", func(t *testing.T) {
		bundles := []Bundle{third}
		result, err := prepareBundles(bundles)
		require.NoError(t, err)
		require.Len(t, result, 4)
		require.Equal(t, []Bundle{third, four, second, first}, result)
		fmt.Println(bundleNames(result).Names())
	})

	t.Run("duplicate bundles filtered correctly", func(t *testing.T) {
		bundles := []Bundle{first, second, four}
		result, err := prepareBundles(bundles)
		require.NoError(t, err)
		require.Len(t, result, 3)
		require.Equal(t, []Bundle{four, second, first}, result)
		fmt.Println(bundleNames(result).Names())
	})

	t.Run("chaos check", func(t *testing.T) {
		bundles := []Bundle{first, second, third, first, second, third, first, second, first, second}
		result, err := prepareBundles(bundles)
		require.NoError(t, err)
		require.Len(t, result, 4)
		require.Equal(t, []Bundle{third, four, second, first}, result)
		fmt.Println(bundleNames(result).Names())
	})
}

var (
	first = Bundle{
		Name: "1:[]",
	}
	second = Bundle{
		Name: "2:1",
		Bundles: []Bundle{
			first,
		},
	}
	third = Bundle{
		Name: "3:2,4",
		Bundles: []Bundle{
			second,
			four,
		},
	}
	four = Bundle{
		Name: "4:1",
		Bundles: []Bundle{
			first,
		},
	}
)
