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

func TestNewSetFromSlice(t *testing.T) {
	t.Run("should return a slice ordered by key ( lexicographic )", func(t *testing.T) {
		actual := NewSetFromSlice("z", "g", "aa")

		assert.NotNil(t, actual)
		assert.NotEmpty(t, actual.List())
		DeepEqual(t, []string{"aa", "g", "z"}, actual.List())
	})

	t.Run("should return a slice of UNIQUE ordered key", func(t *testing.T) {
		actual := NewSetFromSlice("c", "a", "a", "b", "b")

		assert.NotNil(t, actual)
		assert.NotEmpty(t, actual.List())
		assert.Len(t, actual.List(), 3)
		DeepEqual(t, []string{"a", "b", "c"}, actual.List())
	})
}

func TestSet_Remove(t *testing.T) {
	t.Run("should remove an existing key in slice", func(t *testing.T) {
		actual := NewSetFromSlice("c", "a", "a", "b", "b")
		assert.NotNil(t, actual)
		assert.NotEmpty(t, actual.List())
		assert.Len(t, actual.List(), 3)

		actual.Remove("a")
		assert.Len(t, actual.List(), 2)
		DeepEqual(t, []string{"b", "c"}, actual.List())
	})
}

func TestSet_Contains(t *testing.T) {
	t.Run("should contain an existing key", func(t *testing.T) {
		actual := NewSetFromSlice("a", "b", "c")
		assert.NotNil(t, actual)
		assert.NotEmpty(t, actual.List())
		assert.Len(t, actual.List(), 3)

		assert.True(t, actual.Contains("a"), true)
	})
}
