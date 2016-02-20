package main

type Layout struct {
	pressure    int
	spacing     int
	num_spaces  int
	widget_size Size
}

func (layout Layout) width() int {
	return layout.widget_size.width + layout.num_spaces*layout.spacing
}
