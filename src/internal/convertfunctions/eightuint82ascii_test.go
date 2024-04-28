package convertfunctions

import (
	"github.com/brutella/can"
	"testing"
)

func TestEightUint82AsciiToCan(t *testing.T) {
	input := "Gladys"
	output, err := EightUint82AsciiToCan(input)

	// Check whether err is nil
	if err != nil {
		t.Fatalf(`EightUint82AsciiToCan failed, err not nil: %s`, err.Error())
	}

	// Check whether the output has the correct length
	if output.Length != (uint8(len(input))) {
		t.Fatalf(`EightUint82AsciiToCan failed, expected length  %d, actual length %d`, len(input), output.Length)
	}

	// Check if the output has the correct content
	for i := 0; i < len(input); i++ {
		if output.Data[i] != input[i] {
			t.Fatalf(`EightUint82AsciiToCan failed, output wrong at byte %d`, i)
		}
	}

	// Check if back and forth conversion leads to original input
	back, err := EightUint82AsciiToMqtt(output)
	if err != nil {
		t.Fatalf(`EightUint82AsciiToMqtt failed, err not nil: %s`, err.Error())
	}

	if input != back {
		t.Fatalf(`EightUint82AsciiToCan failed, back and forth conversion did not lead to original input`)
	}
}

func TestEightUint82AsciiToMqtt(t *testing.T) {
	input := can.Frame{
		ID:     0,
		Length: 6,
		Flags:  0,
		Res0:   0,
		Res1:   0,
		Data:   [8]uint8{'G', 'l', 'a', 'd', 'y', 's'},
	}
	output, err := EightUint82AsciiToMqtt(input)

	// Check whether err is nil
	if err != nil {
		t.Fatalf(`EightUint82AsciiToMqtt failed, err not nil: %s`, err.Error())
	}

	// Check whether the output has the correct length
	if uint8(len(output)) != input.Length {
		t.Fatalf(`EightUint82AsciiToMqtt failed, expected length  %d, actual length %d`, input.Length, uint8(len(output)))
	}

	// Check if the output has the correct content
	for i := uint8(0); i < input.Length; i++ {
		if output[i] != input.Data[i] {
			t.Fatalf(`EightUint82AsciiToCan failed, output wrong at byte %d`, i)
		}
	}

	// Check if back and forth conversion leads to original input
	back, err := EightUint82AsciiToCan(output)
	if err != nil {
		t.Fatalf(`EightUint82AsciiToCan failed, err not nil: %s`, err.Error())
	}

	if input != back {
		t.Fatalf(`EightUint82AsciiToCan failed, back and forth conversion did not lead to original input`)
	}
}

func FuzzEightUint82AsciiToCan(f *testing.F) {
	f.Fuzz(func(t *testing.T, input string) {
		output, err := EightUint82AsciiToCan(input)
		if err != nil {
			t.Fatalf("%v: decode: %v", input, err)
		}

		if len(input) > 8 {
			t.Logf("input (%s) is larger than 8 bytes (%d), only checking first 8 byte", input, len(input))
			// Check whether the output has the correct length
			if output.Length != 8 {
				t.Fatalf(`EightUint82AsciiToCan failed, expected length  %d, actual length %d`, len(input), output.Length)
			}
			// Check if the output has the correct content
			for i := 0; i < 8; i++ {
				if output.Data[i] != input[i] {
					t.Fatalf(`EightUint82AsciiToCan failed, output wrong at byte %d`, i)
				}
			}
			// Check if back and forth conversion leads to original input
			back, err := EightUint82AsciiToMqtt(output)
			if err != nil {
				t.Fatalf(`EightUint82AsciiToMqtt failed, err not nil: %s`, err.Error())
			}

			// only first 8 bytes are important
			if input[:8] != back {
				t.Fatalf(`EightUint82AsciiToCan failed, back and forth conversion did not lead to original input`)
			}
		} else {
			// Check whether the output has the correct length
			if output.Length != (uint8(len(input))) {
				t.Fatalf(`EightUint82AsciiToCan failed, expected length  %d, actual length %d`, len(input), output.Length)
			}
			// Check if the output has the correct content
			for i := 0; i < len(input); i++ {
				if output.Data[i] != input[i] {
					t.Fatalf(`EightUint82AsciiToCan failed, output wrong at byte %d`, i)
				}
			}

			// Check if back and forth conversion leads to original input
			back, err := EightUint82AsciiToMqtt(output)
			if err != nil {
				t.Fatalf(`EightUint82AsciiToMqtt failed, err not nil: %s`, err.Error())
			}

			if input != back {
				t.Fatalf(`EightUint82AsciiToCan failed, back and forth conversion did not lead to original input`)
			}
		}

	})
}

func FuzzEightUint82AsciiToMqtt(f *testing.F) {
	f.Fuzz(func(t *testing.T, inputString string) {
		var input can.Frame
		if len(inputString) > 8 {
			input.Length = 8
		} else {
			input.Length = uint8(len(inputString))
		}
		for i := uint8(0); i < input.Length; i++ {
			input.Data[i] = inputString[i]
		}
		output, err := EightUint82AsciiToMqtt(input)
		if err != nil {
			t.Fatalf("%v: decode: %v", input, err)
		}

		if input.Length > 8 {
			t.Logf("input.Length (%d) is larger than 8 bytes, only checking first 8 byte", input.Length)
			// Check whether the output has the correct length
			if len(output) != 8 {
				t.Fatalf(`EightUint82AsciiToMqtt failed, expected length  %d, actual length %d`, input.Length, len(output))
			}
			// Check if the output has the correct content
			for i := 0; i < 8; i++ {
				if output[i] != input.Data[i] {
					t.Fatalf(`EightUint82AsciiToMqtt failed, output wrong at byte %d`, i)
				}
			}
			// Check if back and forth conversion leads to original input
			back, err := EightUint82AsciiToCan(output)
			if err != nil {
				t.Fatalf(`EightUint82AsciiToMqtt failed, err not nil: %s`, err.Error())
			}

			// only first 8 bytes are important
			for i := uint8(0); i < 8; i++ {
				if input.Data[i] != back.Data[i] {
					t.Fatalf(`EightUint82AsciiToCan failed, back and forth conversion did not lead to original input`)
				}
			}
		} else {
			// Check whether the output has the correct length
			if uint8(len(output)) != input.Length {
				t.Fatalf(`EightUint82AsciiToMqtt failed, expected length  %d, actual length %d`, input.Length, len(output))
			}
			// Check if the output has the correct content
			for i := uint8(0); i < input.Length; i++ {
				if output[i] != input.Data[i] {
					t.Fatalf(`EightUint82AsciiToMqtt failed, output wrong at byte %d`, i)
				}
			}

			// Check if back and forth conversion leads to original input
			back, err := EightUint82AsciiToCan(output)
			if err != nil {
				t.Fatalf(`EightUint82AsciiToMqtt failed, err not nil: %s`, err.Error())
			}

			if input != back {
				t.Fatalf(`EightUint82AsciiToMqtt failed, back and forth conversion did not lead to original input`)
			}
		}

	})
}