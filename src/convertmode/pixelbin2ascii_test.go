package convertmode

import (
	"github.com/brutella/can"
	"testing"
)

func TestPixelBin2AsciiToCan(t *testing.T) {
	input := []byte("1 #00ff00")
	expected_output := [4]byte{0x1, 0x0, 0xff, 0x0}
	output, err := PixelBin2Ascii{}.ToCan(input)

	// Check whether err is nil
	if err != nil {
		t.Fatalf(`PixelBin2AsciiToCan failed, err not nil: %s`, err.Error())
	}

	// Check whether the output has the correct length
	if int(output.Length) != len(expected_output) {
		t.Fatalf(`PixelBin2AsciiToCan failed, expected length  %d, actual length %d`, len(expected_output), output.Length)
	}

	// Check if the output has the correct content
	for i := 0; i < int(output.Length); i++ {
		if output.Data[i] != expected_output[i] {
			t.Fatalf(`PixelBin2AsciiToCan failed, output wrong at byte %d`, i)
		}
	}

	// Check if back and forth conversion leads to original input
	back, err := PixelBin2Ascii{}.ToMqtt(output)
	if err != nil {
		t.Fatalf(`PixelBin2AsciiToMqtt failed, err not nil: %s`, err.Error())
	}

	if string(input) != string(back) {
		t.Fatalf(`PixelBin2AsciiToCan failed, back and forth conversion (%s) did not lead to original input (%s)`, string(back), string(input))
	}
}

func TestPixelBin2AsciiToMqtt(t *testing.T) {
	input := can.Frame{
		ID:     0,
		Length: 4,
		Flags:  0,
		Res0:   0,
		Res1:   0,
		Data:   [8]uint8{0x1, 0x0, 0xff, 0x0},
	}
	expected_output := "1 #00ff00"
	output, err := PixelBin2Ascii{}.ToMqtt(input)

	// Check whether err is nil
	if err != nil {
		t.Fatalf(`PixelBin2AsciiToMqtt failed, err not nil: %s`, err.Error())
	}

	// Check if the output has the correct content
	if string(output) != expected_output {
		t.Fatalf(`PixelBin2AsciiToCan failed, expected output: %s, actual output: %s`, expected_output, string(output))
	}

	// Check if back and forth conversion leads to original input
	back, err := PixelBin2Ascii{}.ToCan(output)
	if err != nil {
		t.Fatalf(`PixelBin2AsciiToCan failed, err not nil: %s`, err.Error())
	}

	if input != back {
		t.Fatalf(`PixelBin2AsciiToCan failed, back and forth conversion did not lead to original input`)
	}
}