package casper

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type EventListener struct {
	stream    io.ReadCloser
	bufReader *bufio.Reader

	apiVersion string
}

func NewEventListener(baseURL *url.URL, eventPath string, startFrom *uint64) (*EventListener, error) {
	streamURL, err := url.Parse(strings.Join([]string{baseURL.String(), eventPath}, ""))
	if err != nil {
		return nil, err
	}

	if startFrom != nil {
		streamURL.RawQuery = url.Values{
			"start_from": []string{strconv.FormatUint(*startFrom, 10)},
		}.Encode()
	}

	resp, err := connectStream(streamURL)
	if err != nil {
		return nil, err
	}

	listener := EventListener{
		stream:    resp.Body,
		bufReader: bufio.NewReader(resp.Body),
	}

	if err = listener.preFetchAPIVersion(); err != nil {
		return nil, err
	}

	return &listener, nil
}

// APIVersion return api version of provided stream
func (sc *EventListener) APIVersion() string {
	return sc.apiVersion
}

func (sc *EventListener) ReadEvent() (RawEventData, error) {
	// events from node are coming as a pair of events, e.g.
	// data:{"BlockAdded":{...}}
	// id:23388765

	dataEventType, data, err := sc.readLine()
	if err != nil {
		return RawEventData{}, err
	}

	eventIDType, eventIDData, err := sc.readLine()
	if err != nil {
		return RawEventData{}, err
	}

	if eventIDType != EventIDEventType {
		return RawEventData{}, ErrInvalidEventType
	}

	return NewRawEventData(dataEventType, eventIDData, data)
}

// ReadLine synchronously reads bufReader until new line and return RawEventData
func (sc *EventListener) readLine() (EventType, EventData, error) {
	for {
		line, err := sc.bufReader.ReadBytes('\n')
		if err != nil {
			log.Printf("Read error received: %s\n", err.Error())
			return 0, nil, err
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 || bytes.Equal(line, []byte(":")) {
			continue
		}

		return sc.PreProcessLine(string(line))
	}
}

// PreProcessLine create new RawEventData based from line
func (sc *EventListener) PreProcessLine(line string) (EventType, EventData, error) {
	var eventType EventType

	switch {
	case strings.Contains(line, "ApiVersion"):
		eventType = APIVersionEventType
	case strings.Contains(line, "BlockAdded"):
		eventType = BlockAddedEventType
	case strings.Contains(line, "DeployProcessed"):
		eventType = DeployProcessedEventType
	case strings.Contains(line, "DeployAccepted"):
		eventType = DeployAcceptedEventType
	case strings.Contains(line, "DeployExpired"):
		eventType = DeployExpiredEventType
	case strings.Contains(line, "Step"):
		eventType = StepEventType
	case strings.Contains(line, "Fault"):
		eventType = FaultEventType
	case strings.Contains(line, "FinalitySignature"):
		eventType = FinalitySignatureType
	case strings.HasPrefix(line, "id:"):
		return EventIDEventType, bytes.TrimPrefix([]byte(line), []byte("id:")), nil
	default:
		return 0, nil, fmt.Errorf("invalid event received - %s", line)
	}

	return eventType, bytes.TrimPrefix([]byte(line), []byte("data:")), nil
}

// preFetchAPIVersion read first event form the stream
func (sc *EventListener) preFetchAPIVersion() error {
	eventType, eventData, err := sc.readLine()
	if err != nil {
		return err
	}

	if eventType != APIVersionEventType {
		return errors.Wrap(ErrInvalidEventType, "should be - APIVersionEventType")
	}

	apiVersionEvent, err := eventData.ParseAsAPIVersionEvent()
	if err != nil {
		return errors.Wrap(err, "cant receive APIVersion event")
	}
	sc.apiVersion = apiVersionEvent.APIVersion
	return nil
}

// Close send event to closeCh to stop reading stream
func (sc *EventListener) Close() {
	sc.stream.Close()
}

func connectStream(streamURL *url.URL) (*http.Response, error) {
	_, err := net.DialTimeout("tcp", streamURL.Host, time.Second*30)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, streamURL.String(), nil)
	if err != nil {
		log.Printf("Error: %s could not create request to stream: %s\n", err.Error(), streamURL)
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 30,
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error: %s could not connect to stream: %s\n", err.Error(), streamURL)
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Println("Error invalid connect response code")
		return nil, err
	}

	return resp, nil
}
