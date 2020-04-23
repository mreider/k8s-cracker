package main

type IndexedPassword struct {
	Data                []byte
	maxLetterIndex      byte
	maxFirstLetterIndex byte
}

func NewIndexedPassword(initialPassword []byte, maxLetter, maxFirstLetter byte) *IndexedPassword {
	if maxFirstLetter > maxLetter {
		maxFirstLetter = maxLetter
	}

	return &IndexedPassword{
		Data:                initialPassword,
		maxLetterIndex:      maxLetter,
		maxFirstLetterIndex: maxFirstLetter,
	}
}

func (p *IndexedPassword) Increment() bool {
	for i := len(p.Data) - 1; i >= 0; i-- {
		if i > 0 {
			if p.Data[i] < p.maxLetterIndex {
				p.Data[i]++
				return true
			}
			p.Data[i] = 0
		} else {
			if p.Data[i] < p.maxFirstLetterIndex {
				p.Data[i]++
				return true
			}
			for j := 1; j < len(p.Data); j++ {
				p.Data[j] = p.maxLetterIndex
			}
			return false
		}
	}
	return false
}
