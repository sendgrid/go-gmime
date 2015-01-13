package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"runtime/debug"
)

// loops are for manual memory leak check

func TestNewContentDispositionFromString(t *testing.T) {
	loop := 1
	for i := 0; i < loop; i++ {
		contentString := "hola!"
		cd := NewContentDispositionFromString(contentString)
		assert.Equal(t, cd.Disposition(), contentString)
	}
	debug.FreeOSMemory()
}

func TestDispositionDisposition(t *testing.T) {
	loop := 1
	for i := 0; i < loop; i++ {
		contentString := "disposition"
		cd := NewContentDispositionFromString(contentString)

		assert.Equal(t, cd.Disposition(), contentString)
	}
	debug.FreeOSMemory()
}

func TestDispositionIsAttachment(t *testing.T) {
	loop := 1
	for i := 0; i < loop; i++ {
		// Case 1: not an attachment
		contentString := ""
		cd := NewContentDispositionFromString(contentString)
		assert.False(t, cd.IsAttachment())

		// Case 2: is an attachment
		contentString = "attachment"
		cd = NewContentDispositionFromString(contentString)
		assert.True(t, cd.IsAttachment())

		// Case 3: is an inline attachment
		contentString = "inline"
		cd = NewContentDispositionFromString(contentString)
		assert.True(t, cd.IsAttachment())

		// Case 4: anything is an attachment for now
		contentString = "bogus"
		cd = NewContentDispositionFromString(contentString)
		assert.True(t, cd.IsAttachment())
	}
	debug.FreeOSMemory()
}

func TestDispositionToString(t *testing.T) {
	loop := 1
	for i := 0; i < loop; i++ {
		contentString := "ToString()"
		cd := NewContentDispositionFromString(contentString)

		// folding enabled
		assert.Equal(t, cd.ToString(true), contentString+"\n")

		// folding disabled
		assert.Equal(t, cd.ToString(false), contentString)
	}
}
