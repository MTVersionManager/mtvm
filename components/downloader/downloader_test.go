package downloader

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/stretchr/testify/assert"
)

var hundredByteLongLoremIpsum = []byte("Lorem ipsum dolor sit amet consectetur adipiscing elit. Quisque faucibus ex sapien vitae pellentesqu")

func TestDownloadWriter_Write(t *testing.T) {
	dw := downloadWriter{
		totalSize:       100,
		progressChannel: make(chan float64),
	}
	returnDataChannel := make(chan struct {
		int
		error
	})
	dataToWrite := make([]byte, 50)
	go func() {
		written, err := dw.Write(dataToWrite)
		returnDataChannel <- struct {
			int
			error
		}{int: written, error: err}
	}()
	for i := 0; i < 2; i++ {
		select {
		case progress := <-dw.progressChannel:
			assert.Equalf(t, 0.5, progress, "want 0.5 progress, got %v progress", progress)
		case returnedData := <-returnDataChannel:
			if returnedData.error != nil {
				t.Fatal(returnedData.error)
			}
			assert.Equalf(t, 50, returnedData.int, "want 50 bytes written, got %v bytes written", returnedData.int)
		}
	}
	assert.Equalf(t, int64(50), dw.downloadedSize, "want total witten size 50, got %v", dw.downloadedSize)
	assert.Equalf(t, 50, len(dw.downloadedData), "want 50 bytes of content, got %v bytes of content", len(dw.downloadedData))
}

func TestModel_StartDownload(t *testing.T) {
	tests := map[string]struct {
		headers              map[string]string
		options              []Option
		statusCode           int
		testFuncBeforeFinish func(t *testing.T, model Model, msg tea.Msg)
		testFuncAfterFinish  func(t *testing.T, model Model, msg tea.Msg)
		shouldDownload       bool
		chunked              bool
	}{
		"content length present status 200": {
			headers: map[string]string{
				"Content-Length": "100",
			},
			statusCode: 200,
			testFuncBeforeFinish: func(t *testing.T, model Model, msg tea.Msg) {
				if err, ok := msg.(error); ok {
					assert.NoError(t, err)
				}
				assert.IsType(t, DownloadStartedMsg{}, msg)
				downloadStartedMsg := msg.(DownloadStartedMsg)
				assert.True(t, downloadStartedMsg.contentLengthKnown, "want content length to be known, got not known")
				assert.EqualValuesf(t, 100, model.writer.totalSize, "want 100 bytes total size, got %v bytes total size", model.writer.totalSize)
			},
			testFuncAfterFinish: func(t *testing.T, model Model, msg tea.Msg) {
				assert.Equal(t, hundredByteLongLoremIpsum, model.writer.downloadedData)
			},
			shouldDownload: true,
			chunked:        false,
		},
		"content length not present status 200": {
			statusCode:     200,
			shouldDownload: true,
			chunked:        true,
			headers: map[string]string{
				"Transfer-Encoding": "chunked",
			},
			testFuncBeforeFinish: func(t *testing.T, model Model, msg tea.Msg) {
				if err, ok := msg.(error); ok {
					assert.NoError(t, err)
				}
				assert.IsType(t, DownloadStartedMsg{}, msg)
				downloadStartedMsg := msg.(DownloadStartedMsg)
				assert.False(t, downloadStartedMsg.contentLengthKnown, "want content length to be not known, got known")
			},
			testFuncAfterFinish: func(t *testing.T, model Model, msg tea.Msg) {
				assert.Equal(t, hundredByteLongLoremIpsum, model.writer.downloadedData)
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				for header, value := range tt.headers {
					writer.Header().Set(header, value)
				}
				writer.WriteHeader(tt.statusCode)
				var written int
				if tt.chunked {
					writtenChunk1, err := writer.Write(hundredByteLongLoremIpsum[:50])
					assert.NoError(t, err)
					assert.Equalf(t, 50, writtenChunk1, "want 50 bytes written from server for chunk 1, got %v bytes written", writtenChunk1)
					written += writtenChunk1
					writtenChunk2, err := writer.Write(hundredByteLongLoremIpsum[50:])
					assert.NoError(t, err)
					assert.Equalf(t, 50, writtenChunk2, "want 50 bytes written from server for chunk 2, got %v bytes written", writtenChunk2)
					written += writtenChunk2
				} else {
					var err error
					written, err = writer.Write(hundredByteLongLoremIpsum)
					assert.NoError(t, err)
				}
				assert.Equalf(t, 100, written, "want 100 bytes written from server, got %v bytes written", written)
			}))
			defer server.Close()
			model := New(server.URL, tt.options...)
			msg := model.startDownload()
			require.NotNil(t, tt.testFuncBeforeFinish)
			tt.testFuncBeforeFinish(t, model, msg)
			if tt.shouldDownload {
				go func() {
					for range model.writer.progressChannel {
						<-model.writer.progressChannel
					}
				}()
				waitForResponseFinish(model.writer.copyDone)()
				err := model.writer.resp.Body.Close()
				assert.NoErrorf(t, err, "want no error closing response body, got %v", err)
				assert.EqualValuesf(t, 100, model.writer.downloadedSize, "want 100 bytes downloaded, got %v bytes downloaded", model.writer.downloadedSize)
				require.NotNil(t, tt.testFuncAfterFinish)
				tt.testFuncAfterFinish(t, model, msg)
			}
		})
	}
}
