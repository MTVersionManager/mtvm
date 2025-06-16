package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/MTVersionManager/mtvm/shared"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/spf13/afero"

	"github.com/MTVersionManager/mtvm/components/downloadProgress"
	tea "github.com/charmbracelet/bubbletea"
)

type downloadWriter struct {
	totalSize       int64
	downloadedSize  int64
	file            afero.File
	progressChannel chan float64
	resp            *http.Response
	copyDone        chan bool
	downloadedData  []byte
}

type DownloadStartedMsg struct {
	contentLengthKnown bool
	Cancel             context.CancelFunc
}

type DownloadCancelledMsg bool

func (dw *downloadWriter) Start() {
	var err error
	if dw.file == nil {
		_, err = io.Copy(dw, dw.resp.Body)
	} else {
		_, err = io.Copy(dw.file, io.TeeReader(dw.resp.Body, dw))
	}
	// This sends a signal to the update function that it is safe to close the response body
	dw.copyDone <- true
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Println("Error from copying")
		log.Fatal(err)
	}
}

func (dw *downloadWriter) Write(p []byte) (int, error) {
	dw.downloadedSize += int64(len(p))
	if dw.file == nil {
		dw.downloadedData = append(dw.downloadedData, p...)
	}
	if dw.totalSize > 0 && dw.progressChannel != nil {
		dw.progressChannel <- float64(dw.downloadedSize) / float64(dw.totalSize)
	}
	return len(p), nil
}

type Model struct {
	url                string
	downloader         downloadProgress.Model
	progress           float64
	writer             *downloadWriter
	contentLengthKnown bool
	spinner            spinner.Model
	cancel             context.CancelFunc
	Canceled           bool
}

type Option func(Model) Model

func WriteToFs(filePath string, fs afero.Fs) Option {
	return func(model Model) Model {
		file, err := fs.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		model.writer.file = file
		return model
	}
}

func UseTitle(title string) Option {
	return func(model Model) Model {
		model.downloader.Title = title
		return model
	}
}

func New(url string, opts ...Option) Model {
	progressChannel := make(chan float64)
	downloader := downloadProgress.New(progressChannel)
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	model := Model{
		url:        url,
		downloader: downloader,
		progress:   0,
		writer: &downloadWriter{
			progressChannel: progressChannel,
			copyDone:        make(chan bool),
		},
		spinner: spin,
	}
	for _, opt := range opts {
		model = opt(model)
	}
	return model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.startDownload, downloadProgress.WaitForProgress(m.writer.progressChannel), waitForResponseFinish(m.writer.copyDone), m.spinner.Tick)
}

func (m Model) startDownload() tea.Msg {
	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.url, nil)
	if err != nil {
		cancel()
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		cancel()
		return err
	}
	if resp.StatusCode != http.StatusOK {
		cancel()
		return fmt.Errorf("%v %v", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	contentLengthKnown := true
	if resp.ContentLength <= 0 {
		if resp.ContentLength == -1 {
			contentLengthKnown = false
		} else {
			cancel()
			return errors.New("error when getting content length")
		}
	}
	m.writer.totalSize = resp.ContentLength
	m.writer.resp = resp
	if contentLengthKnown && m.writer.file == nil {
		m.writer.downloadedData = make([]byte, 0, m.writer.totalSize)
	}
	go m.writer.Start()
	return DownloadStartedMsg{
		contentLengthKnown: contentLengthKnown,
		Cancel:             cancel,
	}
}

func waitForResponseFinish(doneChan chan bool) tea.Cmd {
	return func() tea.Msg {
		<-doneChan
		return shared.SuccessMsg("download")
	}
}

// GetDownloadedData returns the data that was downloaded.
func (m Model) GetDownloadedData() []byte {
	return m.writer.downloadedData
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case DownloadStartedMsg:
		m.contentLengthKnown = msg.contentLengthKnown
		m.cancel = msg.Cancel
	case shared.SuccessMsg:
		if msg == "download" {
			m.cancel()
			err := m.writer.resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			if m.writer.file != nil {
				err = m.writer.file.Close()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	case DownloadCancelledMsg:
		m.Canceled = true
	}
	var cmd tea.Cmd
	m.downloader, cmd = m.downloader.Update(msg)
	cmds = append(cmds, cmd)
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.contentLengthKnown {
		return m.downloader.View()
	}
	spinnerMsg := "Downloading..."
	if m.downloader.Title != "" {
		spinnerMsg = m.downloader.Title
	}
	return fmt.Sprintf("%v %v\n", m.spinner.View(), spinnerMsg)
}

func (m Model) StopDownload() tea.Cmd {
	if m.cancel == nil {
		return nil
	}
	return func() tea.Msg {
		m.cancel()
		return DownloadCancelledMsg(true)
	}
}

func (m Model) GetUrl() string {
	return m.url
}
