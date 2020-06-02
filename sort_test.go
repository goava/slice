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
	t.Run("cycle", func(t *testing.T) {
		bundles := []Bundle{cycle}
		_, valid := sortBundles(bundles)
		require.False(t, valid)
	})
}

type FirstBundle struct {
}

func (f FirstBundle) Build(builder ContainerBuilder) {}

type SecondBundle struct {
}

func (d SecondBundle) Build(builder ContainerBuilder) {}

func (d SecondBundle) DependOn() []Bundle {
	return []Bundle{
		&FirstBundle{},
	}
}

type ThirdBundle struct {
}

func (t ThirdBundle) Build(builder ContainerBuilder) {}

func (t ThirdBundle) DependOn() []Bundle {
	return []Bundle{
		&SecondBundle{},
		&FourBundle{},
	}
}

type FourBundle struct {
}

func (t FourBundle) Build(builder ContainerBuilder) {}

func (t FourBundle) DependOn() []Bundle {
	return []Bundle{
		&FirstBundle{},
	}
}

type CycleBundle struct {
}

func (c CycleBundle) Build(builder ContainerBuilder) {
}

func (c CycleBundle) DependOn() []Bundle {
	return []Bundle{
		&CycleBundle{},
	}
}

var (
	first  = &FirstBundle{}
	second = &SecondBundle{}
	third  = &ThirdBundle{}
	four   = &FourBundle{}
	cycle  = &CycleBundle{}
)
