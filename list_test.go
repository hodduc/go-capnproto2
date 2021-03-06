package capnp

import (
	"testing"
)

func TestToListDefault(t *testing.T) {
	msg := &Message{Arena: SingleSegment([]byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		42, 0, 0, 0, 0, 0, 0, 0,
	})}
	seg, err := msg.Segment(0)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		ptr  Pointer
		def  []byte
		list List
	}{
		{nil, nil, List{}},
		{Struct{}, nil, List{}},
		{Struct{seg: seg, off: 0}, nil, List{}},
		{List{}, nil, List{}},
		{
			ptr: List{
				seg:    seg,
				off:    8,
				length: 1,
				size:   ObjectSize{DataSize: 8},
			},
			list: List{
				seg:    seg,
				off:    8,
				length: 1,
				size:   ObjectSize{DataSize: 8},
			},
		},
	}

	for _, test := range tests {
		list, err := ToListDefault(test.ptr, test.def)
		if err != nil {
			t.Errorf("ToListDefault(%#v, % 02x) error: %v", test.ptr, test.def, err)
			continue
		}
		if !deepPointerEqual(list, test.list) {
			t.Errorf("ToListDefault(%#v, % 02x) = %#v; want %#v", test.ptr, test.def, list, test.list)
		}
	}
}

func TestListValue(t *testing.T) {
	_, seg, err := NewMessage(SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		list  List
		paddr Address
		val   rawPointer
	}{
		{
			list:  List{},
			paddr: 0,
			val:   0,
		},
		{
			list:  List{seg: seg, length: 3, size: ObjectSize{}},
			paddr: 16,
			val:   0x00000018fffffff5,
		},
		{
			list:  List{seg: seg, off: 24, length: 15, flags: isBitList},
			paddr: 16,
			val:   0x0000007900000001,
		},
		{
			list:  List{seg: seg, off: 40, length: 15, size: ObjectSize{DataSize: 1}},
			paddr: 16,
			val:   0x0000007a00000009,
		},
		{
			list:  List{seg: seg, off: 40, length: 15, size: ObjectSize{DataSize: 2}},
			paddr: 16,
			val:   0x0000007b00000009,
		},
		{
			list:  List{seg: seg, off: 40, length: 15, size: ObjectSize{DataSize: 4}},
			paddr: 16,
			val:   0x0000007c00000009,
		},
		{
			list:  List{seg: seg, off: 40, length: 15, size: ObjectSize{DataSize: 8}},
			paddr: 16,
			val:   0x0000007d00000009,
		},
		{
			list:  List{seg: seg, off: 40, length: 15, size: ObjectSize{PointerCount: 1}},
			paddr: 16,
			val:   0x0000007e00000009,
		},
		{
			list:  List{seg: seg, off: 40, length: 7, size: ObjectSize{DataSize: 16, PointerCount: 1}, flags: isCompositeList},
			paddr: 16,
			val:   0x000000af00000005,
		},
	}
	for _, test := range tests {
		if val := test.list.value(test.paddr); val != test.val {
			t.Errorf("%+v.value(%v) = %#v; want %#v", test.list, test.paddr, val, test.val)
		}
	}
}
