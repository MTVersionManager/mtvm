package shared

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTestingT struct {
	mock.TestingT
	logs   []string
	errors []string
	failed bool
}

func (m *mockTestingT) Logf(format string, args ...interface{}) {
	m.logs = append(m.logs, fmt.Sprintf(format, args...))
}

func (m *mockTestingT) Errorf(format string, args ...interface{}) {
	m.errors = append(m.errors, fmt.Sprintf(format, args...))
}

func (m *mockTestingT) FailNow() {
	m.failed = true
	runtime.Goexit()
}

func newMockTestingT() *mockTestingT {
	return &mockTestingT{
		logs:   []string{},
		errors: []string{},
		failed: false,
	}
}

func TestAssertIsNotFoundError(t *testing.T) {
	normalTestSource := Source{
		File:     "ipsum",
		Function: "dolor",
	}
	alteredTestSource := Source{
		File:     "sit",
		Function: "amet",
	}
	tests := map[string]struct {
		err      error
		thing    string
		source   Source
		testFunc func(mockT *mockTestingT)
	}{
		"matching": {
			err: NotFoundError{
				Thing:  "lorem",
				Source: normalTestSource,
			},
			thing:  "lorem",
			source: normalTestSource,
			testFunc: func(mockT *mockTestingT) {
				assertNotFailedNoLogs(t, *mockT)
				assert.Len(t, mockT.errors, 0, "want no errors, got errors")
			},
		},
		"thing mismatch": {
			err: NotFoundError{
				Thing:  "sit",
				Source: normalTestSource,
			},
			thing:  "lorem",
			source: normalTestSource,
			testFunc: func(mockT *mockTestingT) {
				assertNotFailedNoLogs(t, *mockT)
				getExpectedErrorAndCompare(t, *mockT, func(mockForExpected *mockTestingT) {
					assert.Equal(mockForExpected, "lorem", "sit", "want error to contain thing lorem, got sit")
				})
			},
		},
		"source mismatch": {
			err: NotFoundError{
				Thing:  "lorem",
				Source: alteredTestSource,
			},
			thing:  "lorem",
			source: normalTestSource,
			testFunc: func(mockT *mockTestingT) {
				assertNotFailedNoLogs(t, *mockT)
				getExpectedErrorAndCompare(t, *mockT, func(mockForExpected *mockTestingT) {
					assert.Equalf(mockForExpected, normalTestSource, alteredTestSource, "want error to contain source %v, got %v", normalTestSource, alteredTestSource)
				})
			},
		},
		"wrong error type": {
			err:    errors.New("loremIpsum"),
			thing:  "lorem",
			source: normalTestSource,
			testFunc: func(mockT *mockTestingT) {
				assertFailedNoLogs(t, *mockT)
				getExpectedErrorAndCompare(t, *mockT, func(mockForExpected *mockTestingT) {
					assert.ErrorAs(mockForExpected, errors.New("loremIpsum"), &NotFoundError{})
				})
			},
		},
		"nil error": {
			err:    nil,
			thing:  "lorem",
			source: normalTestSource,
			testFunc: func(mockT *mockTestingT) {
				assertFailedNoLogs(t, *mockT)
				getExpectedErrorAndCompare(t, *mockT, func(mockForExpected *mockTestingT) {
					assert.Error(mockForExpected, nil)
				})
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockT := newMockTestingT()
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				AssertIsNotFoundError(mockT, tt.err, tt.thing, tt.source)
			}()
			wg.Wait()
			// go AssertIsNotFoundError(mockT, tt.err, tt.thing, tt.source)
			tt.testFunc(mockT)
		})
	}
}

func removeErrorTrace(error string) string {
	if i := strings.Index(error, "Error:"); i != -1 {
		return strings.TrimSpace(error[i:])
	}
	return error
}

func getExpectedErrorAndCompare(t *testing.T, mockT mockTestingT, function func(mockForExpected *mockTestingT)) {
	require.Lenf(t, mockT.errors, 1, "want 1 error, got %v errors", len(mockT.errors))
	mockForExpected := newMockTestingT()
	function(mockForExpected)
	require.Lenf(t, mockForExpected.errors, 1, "want 1 error when getting expected, got %v errors", len(mockForExpected.errors))
	assert.Equal(t, removeErrorTrace(mockForExpected.errors[0]), removeErrorTrace(mockT.errors[0]), "unexpected error")
}

func assertNotFailedNoLogs(t *testing.T, mockT mockTestingT) {
	assert.False(t, mockT.failed, "want not failed, got failed")
	assert.Len(t, mockT.logs, 0, "want no logs, got logs")
}

func assertFailedNoLogs(t *testing.T, mockT mockTestingT) {
	assert.True(t, mockT.failed, "want failed, got not failed")
	assert.Len(t, mockT.logs, 0, "want no logs, got logs")
}
