package skasaha

import (
	"github.com/KuroiKitsu/go-gbf"
)

type EventSlice []gbf.Event

func (es EventSlice) Len() int {
	return len(es)
}

func (es EventSlice) Less(i, j int) bool {
	si := es[i].StartsAt
	sj := es[j].StartsAt

	ei := es[i].EndsAt
	ej := es[j].EndsAt

	if si.IsZero() && sj.IsZero() {
		if ej.IsZero() {
			return true
		} else if ei.IsZero() {
			return false
		}

		return ei.Before(ej)
	}

	if sj.IsZero() {
		return true
	} else if si.IsZero() {
		return false
	}

	return si.Before(sj)
}

func (es EventSlice) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}
