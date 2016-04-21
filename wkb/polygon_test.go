package wkb

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	rawPolygon = []byte{
		0x01, 0x03, 0x00, 0x00, 0x00, // header
		0x01, 0x00, 0x00, 0x00, // numlinearring - 1
		0x05, 0x00, 0x00, 0x00, // numpoints - 5
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // point 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // point 2
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // point 3
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // point 4
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // point 5
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
	}
	rawMultiPolygon = []byte{
		0x01, 0x06, 0x00, 0x00, 0x00, // header
		0x02, 0x00, 0x00, 0x00, // numpolygon - 2
		0x01, 0x03, 0x00, 0x00, 0x00, // polygon 1
		0x01, 0x00, 0x00, 0x00, // numlinearring - 1
		0x04, 0x00, 0x00, 0x00, // numpoints - 4,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x01, 0x03, 0x00, 0x00, 0x00, // polygon 2
		0x01, 0x00, 0x00, 0x00, // numlinearring - 1
		0x05, 0x00, 0x00, 0x00, // numpoints - 5
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40,
	}
)

func TestPolygon(t *testing.T) {
	invalid := []struct {
		err error
		b   []byte
	}{
		// invalid type
		{
			ErrUnsupportedValue,
			[]byte{0x01, 0x42, 0x00, 0x00, 0x00},
		},
		// no payload
		{
			ErrInvalidStorage,
			[]byte{0x01, 0x03, 0x00, 0x00, 0x00},
		},
		// no elements
		{
			ErrInvalidStorage,
			[]byte{
				0x01, 0x03, 0x00, 0x00, 0x00, // header
				0x01, 0x00, 0x00, 0x00, // numlinearring - 1
			},
		},
	}

	for _, e := range invalid {
		p := Polygon{}
		if err := p.Scan(e.b); assert.Error(t, err) {
			assert.Exactly(t, e.err, err)
		}
	}

	p := Polygon{}
	if err := p.Scan(rawPolygon); assert.NoError(t, err) {
		assert.Equal(t, Polygon{
			LinearRing{{30, 10}, {40, 40}, {20, 40}, {10, 20}, {30, 10}},
		}, p)
	}

	if raw, err := p.Value(); assert.NoError(t, err) {
		assert.Equal(t, rawPolygon, raw)
	}
}

func TestMultiPolygon(t *testing.T) {
	invalid := []struct {
		err error
		b   []byte
	}{
		// invalid type
		{
			ErrUnsupportedValue,
			[]byte{0x01, 0x42, 0x00, 0x00, 0x00},
		},
		// no payload
		{
			ErrInvalidStorage,
			[]byte{0x01, 0x06, 0x00, 0x00, 0x00},
		},
		// no elements
		{
			ErrInvalidStorage,
			[]byte{
				0x01, 0x06, 0x00, 0x00, 0x00, // header
				0x01, 0x00, 0x00, 0x00, // numpolygon - 2
			},
		},
		// invalid element type
		{
			ErrUnsupportedValue,
			[]byte{
				0x01, 0x06, 0x00, 0x00, 0x00, // header
				0x01, 0x00, 0x00, 0x00, // numpolygon - 2
				0x01, 0x42, 0x00, 0x00, 0x00, // polygon 1
			},
		},
		// no element payload
		{
			ErrInvalidStorage,
			[]byte{
				0x01, 0x06, 0x00, 0x00, 0x00, // header
				0x02, 0x00, 0x00, 0x00, // numpolygon - 2
				0x01, 0x03, 0x00, 0x00, 0x00, // polygon 1
			},
		},
	}

	for _, e := range invalid {
		mp := MultiPolygon{}
		if err := mp.Scan(e.b); assert.Error(t, err) {
			assert.Exactly(t, e.err, err, "Expected MultiPolygon <%v> to fail", hex.EncodeToString(e.b))
		}
	}

	mp := MultiPolygon{}
	if err := mp.Scan(rawMultiPolygon); assert.NoError(t, err) {
		assert.Equal(t, MultiPolygon{
			Polygon{
				LinearRing{{30, 20}, {45, 40}, {10, 40}, {30, 20}},
			},
			Polygon{
				LinearRing{{15, 5}, {40, 10}, {10, 20}, {5, 10}, {15, 5}},
			},
		}, mp)
	}

	if raw, err := mp.Value(); assert.NoError(t, err) {
		assert.Equal(t, rawMultiPolygon, raw)
	}
}
