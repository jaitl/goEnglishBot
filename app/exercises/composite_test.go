package exercises

import (
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompositePuzzle(t *testing.T) {
	phrases := []phrase.Phrase{
		{EnglishText: "look it"},
		{EnglishText: "I get"},
		{EnglishText: "Check it"},
	}
	composite := NewComposite(phrases, PuzzleMode, false)

	assert.Len(t, composite.phrases, 3)
	assert.Equal(t, 0, composite.curPos)

	// 0
	result := composite.Next()
	assert.Len(t, result.Result.Variants, 2)
	assert.False(t, result.IsFinish)

	result = composite.HandleAnswer([]string{"look"})
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 1)

	result = composite.HandleAnswer([]string{"it"})
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 0)

	// 1
	result = composite.Next()
	assert.Len(t, result.Result.Variants, 2)
	assert.False(t, result.IsFinish)
	assert.Equal(t, 1, composite.curPos)

	result = composite.HandleAnswer([]string{"i"})
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 1)

	result = composite.HandleAnswer([]string{"get"})
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 0)

	// 2
	result = composite.Next()
	assert.Len(t, result.Result.Variants, 2)
	assert.False(t, result.IsFinish)
	assert.Equal(t, 2, composite.curPos)

	result = composite.HandleAnswer([]string{"check"})
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 1)

	result = composite.HandleAnswer([]string{"it"})
	assert.True(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 0)

	// correct
	result = composite.Next()
	assert.Len(t, result.Result.Variants, 0)
	assert.True(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.Equal(t, 2, composite.curPos)
}

func TestCompositeWrite(t *testing.T) {
	phrases := []phrase.Phrase{
		{EnglishText: "look it"},
		{EnglishText: "I get"},
		{EnglishText: "Check it"},
	}
	composite := NewComposite(phrases, WriteMode, false)

	assert.Len(t, composite.phrases, 3)
	assert.Equal(t, 0, composite.curPos)

	// 0
	result := composite.Next()
	assert.False(t, result.IsFinish)
	assert.Equal(t, 0, composite.curPos)

	result = composite.HandleAnswer([]string{"look", "it"})
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)

	// 1
	result = composite.Next()
	assert.False(t, result.IsFinish)
	assert.Equal(t, 1, composite.curPos)

	result = composite.HandleAnswer([]string{"i", "get"})
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)

	// 2
	result = composite.Next()
	assert.False(t, result.IsFinish)
	assert.Equal(t, 2, composite.curPos)

	result = composite.HandleAnswer([]string{"check", "it"})
	assert.True(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)

	// correct
	result = composite.Next()
	assert.True(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.Equal(t, 2, composite.curPos)
}