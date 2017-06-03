package skasaha_test

import (
	"testing"

	"github.com/KuroiKitsu/skasaha"
)

func New() *skasaha.Skasaha {
	return &skasaha.Skasaha{}
}

func TestSkasaha_SyncEvents(t *testing.T) {
	s := New()

	err := s.SyncEvents()
	if err != nil {
		t.Fatal(err)
	}

	if len(s.Events) == 0 {
		t.Fatal("no events")
	}

	for i, event := range s.Events {
		if event.Title == "" {
			t.Errorf("empty title at event %d", i)
		}
	}
}
