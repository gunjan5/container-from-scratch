package cmd

import (
	"testing"
)

func TestRunFlagParsing(t *testing.T) {
	cases := []struct {
		testArgs        []string
		skipFlagParsing bool
		expectedErr     error
	}{
		{[]string{"run", "ls", "-la"}, false, nil}, // Test normal case
		{[]string{"blah", "blah"}, true, nil},      // Test SkipFlagParsing without any args that look like flags
	}

}
