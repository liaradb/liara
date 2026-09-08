package link

const RecordLocatorSize = FilePositionSize + SlotIDSize

type RecordLocator struct {
	block  FilePosition
	slotID SlotID
}

func NewRecordLocator(block FilePosition, slotID SlotID) RecordLocator {
	return RecordLocator{
		block:  block,
		slotID: slotID,
	}
}

func (i RecordLocator) Block() FilePosition { return i.block }
func (i RecordLocator) SlotID() int16       { return i.slotID.Value() }
func (i RecordLocator) Size() int           { return RecordLocatorSize }

func (le RecordLocator) Write(data []byte) ([]byte, bool) {
	data0, ok := le.block.WriteData(data)
	if !ok {
		return nil, false
	}

	return le.slotID.WriteData(data0)
}

func (le *RecordLocator) Read(data []byte) ([]byte, bool) {
	data0, ok := le.block.ReadData(data)
	if !ok {
		return nil, false
	}

	return le.slotID.ReadData(data0)
}
