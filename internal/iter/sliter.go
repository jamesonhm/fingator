package iter

type SliceIter struct {
	arrSize   int
	batchSize int
	Start     int
	End       int
}

func NewSliceIter(arrsize, batchsize int) *SliceIter {
	return &SliceIter{
		arrSize:   arrsize,
		batchSize: batchsize,
		Start:     0,
		End:       0,
	}
}

func (si *SliceIter) Next() bool {
	if si.End >= si.arrSize-1 {
		return false
	}
	si.Start = si.End
	si.End = si.Start + si.batchSize
	if si.End > si.arrSize-1 {
		si.End = si.arrSize - 1
	}
	return true
}
