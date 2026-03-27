package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitAndTrim(t *testing.T) {

	t.Run("Basic split and trim", func(t *testing.T) {
		input := "tag1, tag2, tag3"
		expected := []string{"tag1", "tag2", "tag3"}
		result := SplitAndTrim(input, ",")
		assert.Equal(t, expected, result)
	})

	t.Run("Handles extra spaces and empty tags", func(t *testing.T) {
		input := "  tag1 , , tag2 ,  tag3  , "
		expected := []string{"tag1", "tag2", "tag3"}
		result := SplitAndTrim(input, ",")
		assert.Equal(t, expected, result)
	})

	t.Run("Empty input string", func(t *testing.T) {
		input := ""
		expected := []string{}
		result := SplitAndTrim(input, ",")
		assert.Equal(t, expected, result)
	})

	t.Run("No separator in input", func(t *testing.T) {
		input := "singleTag"
		expected := []string{"singleTag"}
		result := SplitAndTrim(input, ",")
		assert.Equal(t, expected, result)
	})
}
