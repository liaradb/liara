package freespace

import (
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/storage/link"
)

const RootSize = 1
const MidSize = 1
const FanOut = 4000

const (
	Level0  = 1 << 0
	Level1  = 1 << 1
	Level2  = 1 << 2
	Level3  = 1 << 3
	Level4  = 1 << 4
	Level5  = 1 << 5
	Level6  = 1 << 6
	Level7  = 1 << 7
	Level8  = 1 << 8
	Level9  = 1 << 9
	Level10 = 1 << 10
)

type MidID int64

func MidIDPosition(i MidID) int64 {
	return RootSize + int64(i)*(FanOut+1)
}

func PageIDPosition(i PageID) int64 {
	return RootSize + MidSize + int64(i)/FanOut + int64(i)
}

type PageID int64

type Page struct {
	id       PageID
	b        *buffer.Buffer
	data     []byte
	tree     []byte
	level0   []byte
	level1   []byte
	level2   []byte
	level3   []byte
	children []byte
}

func GetPage(data []byte, value byte, fanout int16) (int16, bool) {
	if !RootSpace(data, value) {
		return 0, false
	}

	var i int16 = 1
	for i < fanout {
		i = GetDirection(data, i, value)
	}
	return i, true
}

func RootSpace(data []byte, value byte) bool {
	return value >= data[1]
}

func GetDirection(data []byte, i int16, value byte) int16 {
	l := NextLeft(i)
	left := data[l]
	if value <= left {
		return l
	}
	return l + 1
}

func NextLeft(i int16) int16 {
	return i << 1
}

func nextPrent(i int16) int16 {
	return i >> 1
}

func SetPage(data []byte, i int16, value byte, fanout int16) {
	data[i] = value
	p := nextPrent(i)
	for p > 0 {
		if data[p] < value {
			data[p] = value
		}
		p = nextPrent(p)
	}
}

func NewPage(id PageID, data []byte) *Page {
	var (
		tree   = data[1:2048]
		level0 = tree[0:Level0]
		level1 = level0[Level0 : Level0+Level1]
		level2 = level1[Level1 : Level1+Level2]
		level3 = level2[Level2 : Level2+Level3]
	)
	for i := range level0 {
		level0[i] = 1
	}
	for i := range level1 {
		level1[i] = 2
	}
	for i := range level2 {
		level2[i] = 3
	}
	for i := range level3 {
		level3[i] = 4
	}

	return &Page{
		id:       id,
		data:     data,
		tree:     tree,
		level0:   level0,
		level1:   level1,
		level2:   level2,
		level3:   level3,
		children: data[2048:],
	}
}

func (p *Page) Block(i link.BlockID) {
	//    1
	//    2
	//    4
	//    8
	//   16
	//   32
	//   64
	//  128
	//  256
	//  512
	// 1024
	// 2048
	// 4096

}
