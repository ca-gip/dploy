package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet_Concat(t *testing.T) {

	t.Run("with 2 lists should contains 2 items ordered", func(t *testing.T) {
		left := NewSetFromSlice("l1")
		right := NewSetFromSlice("r1")
		concat := left.Concat(right.List())

		assert.NotNil(t, concat)
		assert.NotEmpty(t, concat.List())
		DeepEqual(t, []string{"l1", "r1"}, concat.List())
	})

	t.Run("with 2 lists unordered should contains 2 items ordered", func(t *testing.T) {
		left := NewSetFromSlice("l1")
		right := NewSetFromSlice("r1")
		concat := right.Concat(left.List())

		assert.NotNil(t, concat)
		assert.NotEmpty(t, concat.List())
		DeepEqual(t, []string{"l1", "r1"}, concat.List())
	})

	t.Run("with 2 lists should be ordered", func(t *testing.T) {
		left := NewSetFromSlice("l1")
		right := NewSetFromSlice("r1")
		concat := right.Concat(left.List())

		assert.NotNil(t, concat)
		assert.NotEmpty(t, concat.List())
		DeepEqual(t, []string{"l1", "r1"}, concat.List())
	})

	t.Run("with 3 lists should be ordered", func(t *testing.T) {
		left := NewSetFromSlice("l1")
		right := NewSetFromSlice("r1")
		center := NewSetFromSlice("c1")
		right.Concat(left.List())
		right.Concat(center.List())

		assert.NotNil(t, right)
		assert.NotEmpty(t, right.List())
		DeepEqual(t, []string{"c1", "l1", "r1"}, right.List())
	})

}

func TestSet_Add(t *testing.T) {

	t.Run("with 1 item should be ordered", func(t *testing.T) {
		actual := NewSet()
		actual.Add("s1")

		assert.NotNil(t, actual)
		assert.NotEmpty(t, actual.List())
		DeepEqual(t, []string{"s1"}, actual.List())
	})

	t.Run("with 3 item should be ordered", func(t *testing.T) {
		actual := NewSet()
		actual.Add("s1")
		actual.Add("s2")
		actual.Add("s3")

		assert.NotNil(t, actual)
		assert.NotEmpty(t, actual.List())
		DeepEqual(t, []string{"s1", "s2", "s3"}, actual.List())
	})
}
