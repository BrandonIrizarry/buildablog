package rss

import "time"

func (item Item) Date() time.Time {
	return item.date
}
