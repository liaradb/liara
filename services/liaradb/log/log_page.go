package log

// # Log Records
//
// ## Common to all
// - prevLSN
// - transID
// - type
//
// ## Update records
// - pageID
// - length
// - offset
// - beforeImage
// - afterImage
//
// # Transaction table
// - pageID
// - recLSN
//
// # Dirty page table
// - transID
// - lastLSN

type (
	TransactionID uint64
	PageID        uint64
)

const (
	BlockSize   uint64 = 1024
	SegmentSize uint64 = 1024
)

type LogPage struct {
	Magic           LogMagic
	ID              PageID
	LengthRemaining int
	Records         []*LogRecord
}

type LogRecord struct {
	LogSequenceNumber LogSequenceNumber
	TransactionID     TransactionID
	Length            int
	CRC               CRC
	Data              []byte
	Reverse           []byte
}

func (lr *LogRecord) Write([]byte) error {
	return nil
}

func (lr *LogRecord) Value() []byte {
	data := make([]byte, lr.Size())
	return data
}

func (lr *LogRecord) Size() int {
	return 0
}
