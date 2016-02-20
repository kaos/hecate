package main

import (
	"github.com/nsf/termbox-go"
)

type Widget interface {
	sizeForLayout(layout Layout) Size
	drawAtPoint(layout Layout, point Point, style Style) Size
}

type WidgetSlice []Widget

func (widgets WidgetSlice) sizeForLayout(layout Layout) Size {
	total_widget_width := 0
	max_widget_height := 0
	for _, widget := range widgets {
		widget_size := widget.sizeForLayout(layout)
		total_widget_width += widget_size.width
		if widget_size.height > max_widget_height {
			max_widget_height = widget_size.height
		}
	}
	return Size{total_widget_width, max_widget_height}
}

func (widgets WidgetSlice) numberVisibleForLayout(layout Layout) int {
	count := 0
	for _, widget := range widgets {
		widget_size := widget.sizeForLayout(layout)
		if widget_size.width > 0 {
			count++
		}
	}
	return count
}

func (widgets WidgetSlice) layout() Layout {
	width, _ := termbox.Size()
	layout := Layout{pressure: 0, spacing: 4, num_spaces: 0, widget_size: Size{0, 0}}
	padding := 2
	for ; layout.pressure < 10; layout.pressure++ {
		layout.spacing = 4
		layout.widget_size = widgets.sizeForLayout(layout)
		layout.num_spaces = widgets.numberVisibleForLayout(layout) - 1
		for ; layout.width() > (width-2*padding) && layout.spacing > 2; layout.spacing-- {
		}
		if layout.width() <= (width - 2*padding) {
			break
		}
	}
	return layout
}

func (tab *DataTab) createWidgets() {
	tab.widgets = WidgetSlice{
		NavigationWidget{tab},
		CursorWidget{tab},
		OffsetWidget{tab},
	}
}

func (tab *DataTab) heightOfWidgets() int {
	layout := tab.widgets.layout()
	return layout.widget_size.height
}

func (tab *DataTab) drawWidgets(style Style) Layout {
	width, height := termbox.Size()
	padding := 2
	layout := tab.widgets.layout()
	start_x := (width-2*padding-layout.width())/2 + padding
	start_y := height - layout.widget_size.height
	point := Point{start_x, start_y}
	for _, widget := range tab.widgets {
		widget_size := widget.drawAtPoint(layout, point, style)
		point.x += widget_size.width
		if widget_size.width > 0 {
			point.x += layout.spacing
		}
	}

	return layout
}
