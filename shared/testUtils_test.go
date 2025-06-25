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
				assert.False(t, mockT.failed, "want not failed, got failed")
				assert.Len(t, mockT.logs, 0, "want no logs, got logs")
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
				assert.False(t, mockT.failed, "want not failed, got failed")
				assert.Len(t, mockT.logs, 0, "want no logs, got logs")
				require.Lenf(t, mockT.errors, 1, "want 1 error, got %v errors", len(mockT.errors))
				var expected string
				{
					mockForExpected := newMockTestingT()
					assert.Equal(mockForExpected, "lorem", "sit", "want error to contain thing lorem, got sit")
					require.Lenf(t, mockForExpected.errors, 1, "want 1 error when getting expected, got %v errors", len(mockForExpected.errors))
					expected = removeErrorTrace(mockForExpected.errors[0])
				}
				assert.Equal(t, expected, removeErrorTrace(mockT.errors[0]), "unexpected error")
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
				assert.False(t, mockT.failed, "want not failed, got failed")
				assert.Len(t, mockT.logs, 0, "want no logs, got logs")
				require.Lenf(t, mockT.errors, 1, "want 1 error, got %v errors", len(mockT.errors))
				var expected string
				{
					mockForExpected := newMockTestingT()
					assert.Equalf(mockForExpected, normalTestSource, alteredTestSource, "want error to contain source %v, got %v", normalTestSource, alteredTestSource)
					require.Lenf(t, mockForExpected.errors, 1, "want 1 error when getting expected, got %v errors", len(mockForExpected.errors))
					expected = removeErrorTrace(mockForExpected.errors[0])
				}
				assert.Equal(t, expected, removeErrorTrace(mockT.errors[0]), "unexpected error")
			},
		},
		"wrong error type": {
			err:    errors.New("loremIpsum"),
			thing:  "lorem",
			source: normalTestSource,
			testFunc: func(mockT *mockTestingT) {
				assert.True(t, mockT.failed, "want failed, got not failed")
				assert.Len(t, mockT.logs, 0, "want no logs, got logs")
				require.Lenf(t, mockT.errors, 1, "want 1 error, got %v errors", len(mockT.errors))
				var expected string
				{
					mockForExpected := newMockTestingT()
					assert.ErrorAs(mockForExpected, errors.New("loremIpsum"), &NotFoundError{})
					require.Lenf(t, mockForExpected.errors, 1, "want 1 error when getting expected, got %v errors", len(mockForExpected.errors))
					expected = removeErrorTrace(mockForExpected.errors[0])
				}
				assert.Equal(t, expected, removeErrorTrace(mockT.errors[0]), "unexpected error")
			},
		},
		"nil error": {
			err:    nil,
			thing:  "lorem",
			source: normalTestSource,
			testFunc: func(mockT *mockTestingT) {
				assert.True(t, mockT.failed, "want failed, got not failed")
				assert.Len(t, mockT.logs, 0, "want no logs, got logs")
				require.Lenf(t, mockT.errors, 1, "want 1 error, got %v errors", len(mockT.errors))
				var expected string
				{
					mockForExpected := newMockTestingT()
					assert.NotNil(mockForExpected, nil, "want not nil error, got nil")
					require.Lenf(t, mockForExpected.errors, 1, "want 1 error when getting expected, got %v errors", len(mockForExpected.errors))
					expected = removeErrorTrace(mockForExpected.errors[0])
				}
				assert.Equal(t, expected, removeErrorTrace(mockT.errors[0]), "unexpected error")
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
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
