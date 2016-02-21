package main

import (
	"fmt"
	"math"
	"time"

	"github.com/nsf/termbox-go"
)

type CursorMode int
type TimeSinceEpochUnit int

const MAX_INTEGER_WIDTH = 8
const MIN_INTEGER_WIDTH = 1
const MAX_FLOATING_POINT_WIDTH = 8
const MIN_FLOATING_POINT_WIDTH = 4
const MAX_BIT_PATTERN_WIDTH = 2
const MIN_BIT_PATTERN_WIDTH = 1

const (
	StringMode CursorMode = iota + 1
	BitPatternMode
	IntegerMode
	FloatingPointMode
)

const (
	SecondsSinceEpoch TimeSinceEpochUnit = iota + 1
	DaysSinceEpoch
)

type ByteRange struct {
	pos    int
	length int
}

type Cursor struct {
	pos         int
	max_pos     int
	mark_pos    int
	mark_active bool
	int_length  int
	fp_length   int
	bit_length  int
	mode        CursorMode
	unsigned    bool
	big_endian  bool
	hex_mode    bool
	epoch_time  time.Time
	epoch_unit  TimeSinceEpochUnit
}

func (cursor *Cursor) c_type() string {
	if cursor.mode == IntegerMode {
		if cursor.unsigned {
			return fmt.Sprintf("uint%d_t", cursor.int_length*8)
		} else {
			return fmt.Sprintf(" int%d_t", cursor.int_length*8)
		}
	} else if cursor.mode == FloatingPointMode {
		if cursor.fp_length == 4 {
			return " float"
		} else if cursor.fp_length == 8 {
			return " double"
		}
	} else if cursor.mode == BitPatternMode {
		return " char"
	} else if cursor.mode == StringMode {
		return " char *"
	}
	return ""
}

func (cursor *Cursor) length() int {
	if cursor.mode == IntegerMode {
		return cursor.int_length
	}
	if cursor.mode == FloatingPointMode {
		return cursor.fp_length
	}
	if cursor.mode == BitPatternMode {
		return cursor.bit_length
	}
	return 1
}

func (cursor *Cursor) maximumLength() int {
	if cursor.mode == IntegerMode {
		return MAX_INTEGER_WIDTH
	}
	if cursor.mode == FloatingPointMode {
		return MAX_FLOATING_POINT_WIDTH
	}
	if cursor.mode == BitPatternMode {
		return MAX_BIT_PATTERN_WIDTH
	}
	return 1
}

func (cursor *Cursor) minimumLength() int {
	if cursor.mode == IntegerMode {
		return MIN_INTEGER_WIDTH
	}
	if cursor.mode == FloatingPointMode {
		return MIN_FLOATING_POINT_WIDTH
	}
	if cursor.mode == BitPatternMode {
		return MIN_BIT_PATTERN_WIDTH
	}
	return 1
}

func (cursor *Cursor) grow() {
	if cursor.length() < cursor.maximumLength() {
		if cursor.mode == IntegerMode {
			cursor.int_length *= 2
		} else if cursor.mode == FloatingPointMode {
			cursor.fp_length *= 2
		} else if cursor.mode == BitPatternMode {
			cursor.bit_length *= 2
		}
	}
}

func (cursor *Cursor) shrink() {
	if cursor.length() > cursor.minimumLength() {
		if cursor.mode == IntegerMode {
			cursor.int_length /= 2
		} else if cursor.mode == FloatingPointMode {
			cursor.fp_length /= 2
		} else if cursor.mode == BitPatternMode {
			cursor.bit_length /= 2
		}
	}
}

func (cursor *Cursor) color(style Style) termbox.Attribute {
	if cursor.mode == IntegerMode {
		return style.int_cursor_hex_bg
	}
	if cursor.mode == FloatingPointMode {
		return style.fp_cursor_hex_bg
	}
	if cursor.mode == BitPatternMode {
		return style.bit_cursor_hex_bg
	}
	return style.text_cursor_hex_bg
}

func (cursor *Cursor) toggleMark () {
	cursor.mark_active = !cursor.mark_active
	if cursor.mark_active {
		cursor.mark_pos = cursor.pos
	}
}

func (cursor *Cursor) highlightRange(data []byte) ByteRange {
	var hilite ByteRange

	if cursor.mark_active {
		hilite.pos = Min(cursor.mark_pos, cursor.pos)
		hilite.length = Abs(cursor.mark_pos - cursor.pos) + 1
		return hilite
	}

	if cursor.mode != StringMode || !isPrintable(data[cursor.pos]) {
		return hilite
	}

	hilite.pos = cursor.pos
	for ; hilite.pos > 0 && isPrintable(data[hilite.pos-1]); hilite.pos-- {
	}
	for ; hilite.pos+hilite.length < len(data) && isPrintable(data[hilite.pos+hilite.length]); hilite.length++ {
	}

	return hilite
}

func (cursor *Cursor) interpretBytesAsInteger(data []byte) uint64 {
	var integer uint64
	if cursor.big_endian {
		for i := 0; i < len(data); i++ {
			integer = (integer * 256) + uint64(data[i])
		}
	} else {
		for i := len(data) - 1; i >= 0; i-- {
			integer = (integer * 256) + uint64(data[i])
		}
	}
	return integer
}

func (cursor *Cursor) formatBytesAsNumber(data []byte) string {
	integer := cursor.interpretBytesAsInteger(data)
	if cursor.mode == IntegerMode {
		if cursor.int_length == 1 {
			if cursor.unsigned {
				return fmt.Sprintf("%d", uint8(integer))
			} else {
				return fmt.Sprintf("%d", int8(integer))
			}
		} else if cursor.int_length == 2 {
			if cursor.unsigned {
				return fmt.Sprintf("%d", uint16(integer))
			} else {
				return fmt.Sprintf("%d", int16(integer))
			}
		} else if cursor.int_length == 4 {
			if cursor.unsigned {
				return fmt.Sprintf("%d", uint32(integer))
			} else {
				return fmt.Sprintf("%d", int32(integer))
			}
		} else if cursor.int_length == 8 {
			if cursor.unsigned {
				return fmt.Sprintf("%d", uint64(integer))
			} else {
				return fmt.Sprintf("%d", int64(integer))
			}
		}
	} else if cursor.mode == FloatingPointMode {
		if cursor.fp_length == 4 {
			var integer32 uint32 = uint32(integer)
			return fmt.Sprintf("%.5g", math.Float32frombits(integer32))
		} else if cursor.fp_length == 8 {
			return fmt.Sprintf("%g", math.Float64frombits(integer))
		}
	} else if cursor.mode == BitPatternMode {
		str := ""
		for i := 0; i < cursor.bit_length; i++ {
			str = fmt.Sprintf("%s%02x", str, data[i])
		}
		return str
	}

	return string(data)
}

func (cursor *Cursor) interpretBytesAsTime(data []byte) time.Time {
	integer := cursor.interpretBytesAsInteger(data)
	var date_time time.Time
	if cursor.mode == IntegerMode {
		if cursor.epoch_unit == SecondsSinceEpoch {
			date_time = cursor.epoch_time.Add(time.Duration(integer) * time.Second)
		} else if cursor.epoch_unit == DaysSinceEpoch {
			date_time = cursor.epoch_time.Add(time.Duration(integer) * 24 * time.Hour)
		}
	} else if cursor.mode == FloatingPointMode {
		var float float64
		if cursor.fp_length == 4 {
			float = float64(math.Float32frombits(uint32(integer)))
		} else {
			float = math.Float64frombits(integer)
		}
		if cursor.epoch_unit == SecondsSinceEpoch {
			date_time = cursor.epoch_time.Add(time.Duration(float * float64(time.Second)))
		} else {
			date_time = cursor.epoch_time.Add(time.Duration(float * 24 * float64(time.Hour)))
		}
	} else {
		date_time = cursor.epoch_time
	}
	return date_time.UTC()
}

func (cursor *Cursor) setPos (pos int) {
	if pos < 0 {
		pos = 0
	} else if pos + cursor.length() > cursor.max_pos {
		pos = cursor.max_pos - cursor.length()
	}
	cursor.pos = pos
}

func (cursor *Cursor) move (delta int) {
	cursor.setPos (cursor.pos + delta)
}
