package webtransport

import (
	"fmt"
	"sync"
	"bytes"
	"context"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/quicvarint"
)

// A StreamID in QUIC
type StreamID int64

type DatagramMap struct {
	mutex sync.RWMutex
	conn quic.Connection
	datagrammers map[StreamID]*SessionWithDatagram
}

func (d *DatagramMap) RunReceiving(){
	for {
		data, err := d.conn.ReceiveDatagram(context.Background())
		if err != nil {
			fmt.Println("Stop receiving datagram: %s", err)
			return
		}
		buf := bytes.NewBuffer(data)
		quarterStreamID, err := quicvarint.Read(buf)
		if err != nil {
			fmt.Println("Reading datagram Quarter Stream ID failed: %s", err)
			continue
		}
		streamID := quarterStreamID * 4
		d.mutex.RLock()
		stream, ok := d.datagrammers[StreamID(streamID)]
		d.mutex.RUnlock()
		if !ok {
			fmt.Println("Received datagram for unknown stream: %d", streamID)
			continue
		}
		stream.handleDatagram(buf.Bytes())
	}
}

type SessionWithDatagram struct {
	Session
	qconn quic.Connection
	datagramCh chan []byte
}

func (s *SessionWithDatagram) SendDatagram([]byte) error {

}

func (d *SessionWithDatagram) ReceiveDatagram(context.Context) ([]byte, error) {

}