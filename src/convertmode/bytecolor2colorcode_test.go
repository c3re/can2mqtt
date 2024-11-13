package convertmode

import (
	"github.com/brutella/can"
	"reflect"
	"testing"
)

func TestByteColor2ColorCode_String(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{{"default", "bytecolor2colorcode"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			by := ByteColor2ColorCode{}
			if got := by.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteColor2ColorCode_ToCan(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		args    args
		want    can.Frame
		wantErr bool
	}{
		{
			"empty input",
			args{input: make([]byte, 0)},
			can.Frame{},
			true,
		},
		{
			"regular",
			args{input: []byte("#00ff00")},
			can.Frame{Length: 3, Data: [8]uint8{0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
			false,
		},
		{
			"mixed caps",
			args{input: []byte("#AbFf0E")},
			can.Frame{Length: 3, Data: [8]uint8{0xab, 0xff, 0xe, 0x0, 0x0, 0x0, 0x0, 0x0}},
			false,
		},
		{
			"missing no sign",
			args{input: []byte("AbFf0E")},
			can.Frame{Length: 3, Data: [8]uint8{0xab, 0xff, 0xe, 0x0, 0x0, 0x0, 0x0, 0x0}},
			false,
		},
		{
			"multiple no sign",
			args{input: []byte("###AbFf0E")},
			can.Frame{},
			true,
		},
		{
			"input too long",
			args{input: []byte("#AAAAAAE")},
			can.Frame{},
			true,
		},
		{
			"input too short",
			args{input: []byte("#AAf0E")},
			can.Frame{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			by := ByteColor2ColorCode{}
			got, err := by.ToCan(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToCan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToCan() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteColor2ColorCode_ToMqtt(t *testing.T) {
	type args struct {
		input can.Frame
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"good1",
			args{can.Frame{Length: 3,
				Data: [8]uint8{0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}},
			[]byte("#00ff00"),
			false,
		},
		{
			"good2",
			args{can.Frame{Length: 3,
				Data: [8]uint8{0xab, 0xff, 0x12, 0x0, 0x0, 0x0, 0x0, 0x0}}},
			[]byte("#abff12"),
			false,
		},
		{
			"good3",
			args{can.Frame{Length: 3,
				Data: [8]uint8{0xab, 0xcd, 0xef, 0x0, 0x0, 0x0, 0x0, 0x0}}},
			[]byte("#abcdef"),
			false,
		},
		{
			"frame too short",
			args{can.Frame{Length: 2,
				Data: [8]uint8{0xab, 0xcd, 0xef, 0x0, 0x0, 0x0, 0x0, 0x0}}},
			[]byte{},
			true,
		},
		{
			"frame too long",
			args{can.Frame{Length: 4,
				Data: [8]uint8{0xab, 0xcd, 0xef, 0x0, 0x0, 0x0, 0x0, 0x0}}},
			[]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			by := ByteColor2ColorCode{}
			got, err := by.ToMqtt(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToMqtt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMqtt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
