package downloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
