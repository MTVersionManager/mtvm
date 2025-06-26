package shared

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type mockMemMapFs struct {
	*afero.MemMapFs
	openError error
}

func (m *mockMemMapFs) Stat(name string) (os.FileInfo, error) {
	if m.openError != nil {
		return nil, m.openError
	}
	return m.MemMapFs.Stat(name)
}

func newMockMemMapFs(openError error) afero.Fs {
	return &mockMemMapFs{
		MemMapFs:  &afero.MemMapFs{},
		openError: openError,
	}
}

func TestIsVersionInstalled(t *testing.T) {
	viper.Reset()
	viper.Set("installDir", "/test/install/")
	err := viper.Unmarshal(&Configuration)
	assert.NoError(t, err)
	tests := map[string]struct {
		tool, version string
		fs            afero.Fs
		want          bool
		errorCheck    func(t *testing.T, err error)
	}{
		"installed": {
			tool:    "test",
			version: "0.0.0",
			want:    true,
			errorCheck: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				err := fs.MkdirAll("/test/install/test/0.0.0", 0o777)
				assert.NoError(t, err)
				return fs
			}(),
		},
		"not installed": {
			tool:    "test",
			version: "0.0.0",
			want:    false,
			errorCheck: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			fs: afero.NewMemMapFs(),
		},
		"permission error": {
			tool:    "test",
			version: "0.0.0",
			want:    false,
			errorCheck: func(t *testing.T, err error) {
				assert.Error(t, err)
				if !os.IsPermission(err) {
					t.Errorf("want permission error, got %v", err)
				}
			},
			fs: func() afero.Fs {
				fs := newMockMemMapFs(os.ErrPermission)
				err := fs.MkdirAll("/test/install/test/0.0.0", 0o000)
				assert.NoError(t, err)
				return fs
			}(),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.NotNil(t, tt.fs)
			installed, err := IsVersionInstalled(tt.tool, tt.version, tt.fs)
			tt.errorCheck(t, err)
			assert.Equal(t, tt.want, installed)
		})
	}
}
