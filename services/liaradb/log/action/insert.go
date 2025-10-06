package action

type Insert struct {
	collectionID CollectionID
	itemID       ItemID
	data         []byte
}

func NewInsert(
	collectionID CollectionID,
	itemID ItemID,
	data []byte,
) *Insert {
	return &Insert{
		collectionID: collectionID,
		itemID:       itemID,
		data:         data,
	}
}

func (i *Insert) CollectionID() CollectionID { return i.collectionID }
func (i *Insert) ItemID() ItemID             { return i.itemID }
func (i *Insert) Data() []byte               { return i.data }
