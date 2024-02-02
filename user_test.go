package sdump

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCounter_TakeN(t *testing.T) {
	tt := []struct {
		name          string
		initialState  int64
		itemsToTake   int64
		expectedValue int64
		hasError      bool
	}{
		{
			name:          "can take from non zero counter",
			initialState:  10,
			itemsToTake:   1,
			expectedValue: 9,
			hasError:      false,
		},
		{
			name:          "cannot take from zero counter",
			initialState:  0,
			itemsToTake:   1,
			expectedValue: 0,
			hasError:      true,
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			c := Counter(v.initialState)

			err := c.TakeN(v.itemsToTake)

			if v.hasError {
				require.Error(t, err)
				return
			}

			require.Equal(t, Counter(v.expectedValue), c)
		})
	}
}

func TestCounter_Add(t *testing.T) {
	tt := []struct {
		name          string
		initialState  int64
		expectedValue int64
	}{
		{
			name:          "zero couner can be increased",
			initialState:  0,
			expectedValue: 1,
		},
		{
			name:          "non zero counter can be increased",
			initialState:  1,
			expectedValue: 2,
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			c := Counter(v.initialState)

			c.Add()

			require.Equal(t, Counter(v.expectedValue), c)
		})
	}
}
