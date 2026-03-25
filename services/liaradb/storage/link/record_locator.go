package link

const RecordLocatorSize = FilePositionSize + RecordPositionSize

type RecordLocator struct {
	block    FilePosition
	position RecordPosition
}

func NewRecordLocator(block FilePosition, position RecordPosition) RecordLocator {
	return RecordLocator{
		block:    block,
		position: position,
	}
}

func (i RecordLocator) Block() FilePosition { return i.block }
func (i RecordLocator) Position() int16     { return i.position.Value() }
func (i RecordLocator) Size() int           { return RecordLocatorSize }

func (le RecordLocator) Write(data []byte) ([]byte, bool) {
	data0, ok := le.block.WriteData(data)
	if !ok {
		return nil, false
	}

	return le.position.WriteData(data0)
}

func (le *RecordLocator) Read(data []byte) ([]byte, bool) {
	data0, ok := le.block.ReadData(data)
	if !ok {
		return nil, false
	}

	return le.position.ReadData(data0)
}
