package record

import "io"

type Page struct {
	size   Offset
	list   List
	items  []Item
	cursor Offset
}

type Item = []byte

func (p *Page) Add(i Item) {
	l := len(i)
	p.list.Add(Offset(l))
	// p.items = append(p.items, i)
}

func (p *Page) Size() int {
	return p.list.Size()
}

func (p *Page) Write(w interface {
	// io.WriterAt
	io.Writer
}) error {
	if err := p.list.Write(w); err != nil {
		return err
	}

	// for _, i := range p.items {
	// 	if err := p.writeItem(w, i); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (p *Page) Read(r io.Reader) error {
	return p.list.Read(r)
}

func (p *Page) writeItem(w io.WriterAt, item Item) error {
	if len(p.items) == 0 {
		p.cursor = p.size
	}

	p.cursor -= Offset(len(item))

	_, err := w.WriteAt(item, int64(p.cursor))
	return err
}
