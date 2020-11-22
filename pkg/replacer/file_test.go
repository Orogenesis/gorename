package replacer

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	testsTrue := [][2]string{
		{"github.com/stretchr/testify", "github.com/stretchr/testify"},
		{"github.com/stretchr/testify", "github.com/stretchr/testify/assert"},
		{"github.com/stretchr/testify", "github.com/stretchr/testify/assert/second"},
		{"github.com/stretchr/testify", "github.com/stretchr/testify/assert/second/third"},
	}

	for _, testTrue := range testsTrue {
		assert.True(t, mustMatchPath(strconv.Quote(testTrue[0]), testTrue[1]))
	}

	testsFalse := [][2]string{
		{"github.com/stretchr/testify", "github.com/stretchr/test"},
		{"github.com/stretchr/testify", "github.com/golang/mod"},
		{"github.com/stretchr/testify", "gitlab.com/golang/mod"},
	}

	for _, testFalse := range testsFalse {
		assert.False(t, mustMatchPath(strconv.Quote(testFalse[0]), testFalse[1]))
	}
}
