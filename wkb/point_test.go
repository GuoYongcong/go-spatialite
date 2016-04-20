package wkb

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint(t *testing.T) {
	invalid := map[error][]byte{
		ErrInvalidStorage: {
			0x01,
		}, // header too short
		ErrInvalidStorage: {
			0x01, 0x01, 0x00, 0x00, 0x00, 0x00,
		}, // no payload
		ErrInvalidStorage: {
			0x02, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		}, // invalid endianness
		ErrInvalidStorage: {
			0x01, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		}, // single coordinate only
		ErrUnsupportedValue: {
			0x01, 0x02, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		}, // invalid type
	}

	for expected, b := range invalid {
		p := Point{}
		if err := p.Scan(b); assert.Error(t, err) {
			assert.Exactly(t, expected, err, "Expected point <%s> to fail", hex.EncodeToString(b))
		}
	}

	valid := []byte{
		0x01, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
	}
	p := Point{}
	if assert.NoError(t, p.Scan(valid)) {
		assert.Equal(t, Point{30, 10}, p)
	}
}

func TestMultipoint(t *testing.T) {
	invalid := map[error][]byte{
		ErrInvalidStorage: {
			0x01, 0x04, 0x00, 0x00, 0x00, 0x00,
		}, // no payload
		ErrUnsupportedValue: {
			0x01, 0x42, 0x00, 0x00, 0x00, 0x00,
		}, // invalid type
		ErrInvalidStorage: {
			0x01, 0x04, 0x00, 0x00, 0x00, 0x00,
			0x01, 0x00, 0x00, 0x00, // numpoints - 1
		}, // no points
	}

	for expected, b := range invalid {
		mp := MultiPoint{}
		if err := mp.Scan(b); assert.Error(t, err) {
			assert.Exactly(t, expected, err, "Expected multipoint <%s> to fail", hex.EncodeToString(b))
		}
	}

	valid := []byte{
		0x01, 0x04, 0x00, 0x00, 0x00, // header
		0x04, 0x00, 0x00, 0x00, // numpoints - 4
		0x01, 0x01, 0x00, 0x00, 0x00, // point 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x01, 0x01, 0x00, 0x00, 0x00, // point 2
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x01, 0x01, 0x00, 0x00, 0x00, // point 3
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x01, 0x01, 0x00, 0x00, 0x00, // point 4
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
	}

	mp := MultiPoint{}
	if assert.NoError(t, mp.Scan(valid)) && assert.Equal(t, 4, len(mp)) {
		assert.Equal(t, Point{10, 40}, mp[0])
		assert.Equal(t, Point{40, 30}, mp[1])
		assert.Equal(t, Point{20, 20}, mp[2])
		assert.Equal(t, Point{30, 10}, mp[3])
	}
}
