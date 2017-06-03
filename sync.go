package skasaha

import (
	"sort"

	"github.com/KuroiKitsu/go-gbf"
	"github.com/KuroiKitsu/go-gbf/scraper"
	"github.com/bwmarrin/discordgo"
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

func (s *Skasaha) Sync() error {
	s.initEmbeds()

	var (
		err error
	)

	err = s.SyncEvents()
	if err != nil {
		return err
	}

	err = s.rebuildIndex()

	return err
}

func (s *Skasaha) SyncEvents() error {
	var (
		err error

		current  []*scraper.Event
		upcoming []*scraper.Event

		events = make([]gbf.Event, 0)
		log    = s.logger()
	)

	log.Print("synchronizing events")

	current, err = scraper.CurrentEvents()
	if err != nil {
		return err
	}

	upcoming, err = scraper.UpcomingEvents()
	if err != nil {
		return err
	}

	scrapedEvents := make([]*scraper.Event, 0, len(current)+len(upcoming))

	scrapedEvents = append(scrapedEvents, current...)
	scrapedEvents = append(scrapedEvents, upcoming...)

	for _, scraperEvent := range scrapedEvents {
		var (
			event   gbf.Event
			details *scraper.EventDetails
		)

		event.Title = scraperEvent.Title
		event.URL = scraperEvent.URL

		if scraperEvent.URL != "" {
			details, err = scraper.EventDetailsURL(scraperEvent.URL)
			if err != nil {
				return err
			}

			event.ImageURL = details.ImageURL
			event.Description = details.Description
			event.StartsAt = details.StartsAt
			event.EndsAt = details.EndsAt
		}

		events = append(events, event)
	}

	sort.Sort(EventSlice(events))

	if len(s.embedsEvents) > 0 {
		for _, id := range s.embedsEvents {
			delete(s.embeds, id)
		}
	}

	s.Events = events

	s.embedsEvents = make([]Snowflake, 0, len(s.Events))
	for _, event := range s.Events {
		id := NewSnowflake()

		em := &discordgo.MessageEmbed{
			Title:       event.Title,
			Description: event.Description,
			URL:         event.URL,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL:    GBFLogoImageURL,
				Width:  GBFLogoImageWidth,
				Height: GBFLogoImageHeight,
			},
			Provider: &discordgo.MessageEmbedProvider{
				Name: "Granblue Fantasy Wiki",
				URL:  scraper.WikiHomeURL.String(),
			},
			Fields: make([]*discordgo.MessageEmbedField, 0),
		}

		startsAt := "N/A"
		endsAt := "N/A"

		if !event.StartsAt.IsZero() {
			startsAt = event.StartsAt.UTC().Format("Jan 02")
		}
		if !event.EndsAt.IsZero() {
			endsAt = event.EndsAt.UTC().Format("Jan 02")
		}

		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   "Starts",
			Value:  startsAt,
			Inline: true,
		})

		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   "Ends",
			Value:  endsAt,
			Inline: true,
		})

		if event.ImageURL != "" {
			em.Image = &discordgo.MessageEmbedImage{
				URL: event.ImageURL,
			}
		}

		s.embeds[id] = em
		s.embedsEvents = append(s.embedsEvents, id)
	}

	log.Printf("%d events found", len(s.Events))

	return err
}
