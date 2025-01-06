package downloader

import "testing"

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
			if progress != 0.5 {
				t.Fatalf("want 0.5 progress, got %v progress", progress)
			}
		case returnedData := <-returnDataChannel:
			if returnedData.error != nil {
				t.Fatal(returnedData.error)
			}
			if returnedData.int != 50 {
				t.Fatalf("want 50 bytes written, got %v bytes written", returnedData.int)
			}
		}
	}
	if dw.downloadedSize != 50 {
		t.Fatalf("want total witten size 50, got %v", dw.downloadedSize)
	}
	if len(dw.downloadedData) != 50 {
		t.Fatalf("want 50 bytes of content, got %v bytes of content", len(dw.downloadedData))
	}
}
