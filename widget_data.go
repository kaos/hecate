package main

import (
	"crypto/md5"
	"fmt"
)

type DataWidget struct {
	*DataTab
}

func (widget DataWidget) sizeForLayout(layout Layout) Size {
	if layout.pressure > 0 {
		return Size{0, 0}
	}
	return Size{20, 4}
}

func (widget DataWidget) drawAtPoint(layout Layout, point Point, style Style) Size {
	fg := style.default_fg
	bg := style.default_bg
	x_pos := point.x
	if widget.hilite.length > 0 {
		drawStringAtPoint(
			fmt.Sprintf("MD5  %x", md5.Sum(
				widget.bytes[widget.hilite.pos : widget.hilite.pos + widget.hilite.length])),
			x_pos, point.y, fg, bg)
	}
	width := drawStringAtPoint(fmt.Sprintf("%#v", widget.hilite), point.x, point.y + 2, fg, bg)
	return Size{width, 4}
}
