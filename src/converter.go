package main

import (
	"fmt"

	"github.com/jaster-prj/can2mqtt/convertmode"
)

type ConverterFactory struct {
	converterMap map[string]ConvertMode
}

func NewConverterFactory() *ConverterFactory {
	converterMap := map[string]ConvertMode{}
	// initialize all convertModes
	converterMap[convertmode.None{}.String()] = convertmode.None{}
	converterMap[convertmode.SixteenBool2Ascii{}.String()] = convertmode.SixteenBool2Ascii{}
	converterMap[convertmode.PixelBin2Ascii{}.String()] = convertmode.PixelBin2Ascii{}
	converterMap[convertmode.ByteColor2ColorCode{}.String()] = convertmode.ByteColor2ColorCode{}
	converterMap[convertmode.MyMode{}.String()] = convertmode.MyMode{}
	// Dynamically create int and uint convertmodes
	for _, bits := range []uint{8, 16, 32, 64} {
		for _, instances := range []uint{1, 2, 4, 8} {
			if bits*instances <= 64 {
				// int
				cmi, _ := convertmode.NewInt2Ascii(instances, bits)
				converterMap[cmi.String()] = cmi
				// uint
				cmu, _ := convertmode.NewUint2Ascii(instances, bits)
				converterMap[cmu.String()] = cmu
			}
		}
	}
	return &ConverterFactory{
		converterMap: converterMap,
	}
}

func (c *ConverterFactory) GetConverter(converter string) (ConvertMode, error) {
	conv, ok := c.converterMap[converter]
	if !ok {
		return nil, fmt.Errorf("converter not initialized: %s", converter)
	}
	return conv, nil
}
