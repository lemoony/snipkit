package stringutil

import (
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/stretchr/testify/assert"
)

func Test_StringSet(t *testing.T) {
	set := StringSet{}
	set.Add("Foo")
	set.Add("Moo")

	assert.Len(t, set, 2)
	assert.True(t, set.Contains("Foo"))
	assert.True(t, set.Contains("Moo"))
	assert.False(t, set.Contains("Another"))

	then.AssertThat(t,
		set.Keys(),
		is.AllOf(is.ArrayContaining("Foo"), is.ArrayContaining("Moo")),
	)
}
