package orm

import (
	"reflect"

	"github.com/deliveroo/pg-roo/internal"
	"github.com/deliveroo/pg-roo/types"
)

type sliceModel struct {
	Discard
	slice    reflect.Value
	nextElem func() reflect.Value
	scan     func(reflect.Value, types.Reader, int) error
}

var _ Model = (*sliceModel)(nil)

func newSliceModel(slice reflect.Value, elemType reflect.Type) *sliceModel {
	return &sliceModel{
		slice: slice,
		scan:  types.Scanner(elemType),
	}
}

func (m *sliceModel) Init() error {
	if m.slice.IsValid() && m.slice.Len() > 0 {
		m.slice.Set(m.slice.Slice(0, 0))
	}
	return nil
}

func (m *sliceModel) NewModel() ColumnScanner {
	return m
}

func (m *sliceModel) ScanColumn(colIdx int, _ string, rd types.Reader, n int) error {
	if m.nextElem == nil {
		m.nextElem = internal.MakeSliceNextElemFunc(m.slice)
	}
	v := m.nextElem()
	return m.scan(v, rd, n)
}
