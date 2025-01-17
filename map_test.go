package godash_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/godash"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("support primitive types", func(t *testing.T) {
		in := []int{1, 2, 3}
		out := make([]int, 0)

		err := godash.Map(in, &out, func(element int) int {
			return element * element
		})

		expected := []int{1, 4, 9}
		assert.NoError(t, err)
		assert.Equal(t, expected, out)
	})

	t.Run("support structs", func(t *testing.T) {
		type person struct {
			name string
		}

		in := []person{
			{name: "john"},
			{name: "doe"},
		}
		out := make([]string, 0)
		expected := []string{"john", "doe"}

		err := godash.Map(in, &out, func(p person) string {
			return p.name
		})

		assert.NoError(t, err)
		assert.Equal(t, expected, out)
	})

	squared := func(element int) int {
		return element * element
	}

	t.Run("should not panic if output is nil", func(t *testing.T) {
		in := []int{1, 2, 3}
		{
			var out []int

			err := godash.Map(in, out, squared)

			assert.EqualError(t, err, "output is nil. Pass a reference to set output")
		}

		{
			err := godash.Map(in, nil, squared)

			assert.EqualError(t, err, "output is nil. Pass a reference to set output")
		}
	})

	t.Run("should not panic if output is not a slice", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out int

		err := godash.Map(in, &out, squared)

		assert.EqualError(t, err, "output should be a slice for input of type slice")
	})

	t.Run("should not accept mapper function that are not functions", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out []int

		err := godash.Map(in, &out, 7)

		assert.EqualError(t, err, "mapperFn has to be a function")
	})

	t.Run("should not accept mapper function that do not take exactly one argument", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out []int

		{
			err := godash.Map(in, &out, func() int { return 0 })
			assert.EqualError(t, err, "mapper function has to take only one argument")
		}

		{
			err := godash.Map(in, &out, func(int, int) int { return 0 })
			assert.EqualError(t, err, "mapper function has to take only one argument")
		}
	})

	t.Run("should not accept mapper function that do not return exactly one value", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out []int

		{
			err := godash.Map(in, &out, func(int) {})
			assert.EqualError(t, err, "mapper function should return only one return value")
		}

		{
			err := godash.Map(in, &out, func(int) (int, int) { return 0, 0 })
			assert.EqualError(t, err, "mapper function should return only one return value")
		}
	})

	t.Run("should accept mapper function whose argument's kind should be slice's element kind", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out []int

		{
			err := godash.Map(in, &out, func(string) string { return "" })
			assert.EqualError(t, err, "mapper function's first argument (string) has to be (int)")
		}

		{
			err := godash.Map(in, &out, func(int) int { return 0 })
			assert.NoError(t, err)
		}
	})

	t.Run("should accept mapper function whose return's kind should be  output slice's element kind", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out []string

		{
			err := godash.Map(in, &out, func(int) int { return 0 })
			assert.EqualError(t, err, "mapper function's return type has to be (int) but is (string)")
		}

		{
			err := godash.Map(in, &out, func(int) string { return "" })
			assert.NoError(t, err)
		}
	})
}

func ExampleMap() {
	input := []int{0, 1, 2, 3, 4}
	var output []string

	_ = godash.Map(input, &output, func(num int) string {
		return fmt.Sprintf("%d", num*num)
	})

	fmt.Println(output)

	// Output: [0 1 4 9 16]
}
