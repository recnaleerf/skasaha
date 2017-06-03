package skasaha

import (
	"fmt"
	"time"

	"github.com/KuroiKitsu/go-gbf/scraper"
	"github.com/bwmarrin/discordgo"
)

// cmdEvents handles the events command.
func (s *Skasaha) cmdEvents(ds *discordgo.Session, m *discordgo.MessageCreate) error {
	var (
		err error

		log = s.logger()
	)

	em := &discordgo.MessageEmbed{
		Title: "Events",
		URL:   scraper.WikiEventsURL.String(),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    GBFLogoImageURL,
			Width:  GBFLogoImageWidth,
			Height: GBFLogoImageHeight,
		},
		Provider: &discordgo.MessageEmbedProvider{
			Name: "Granblue Fantasy Wiki",
			URL:  scraper.WikiHomeURL.String(),
		},
		Fields: make([]*discordgo.MessageEmbedField, 0, len(s.Events)),
	}

	now := time.Now()

	for _, event := range s.Events {
		value := "-"

		if event.StartsAt.Before(now) && event.EndsAt.After(now) {
			endsAt := event.EndsAt.UTC().Format("2006-01-02 15:04:05 MST")
			value = fmt.Sprintf("Ends on %s", endsAt)
		} else if event.StartsAt.After(now) {
			startsAt := event.StartsAt.UTC().Format("2006-01-02 15:04:05 MST")
			value = fmt.Sprintf("Begins on %s", startsAt)
		} else {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name:   event.Title,
				Value:  "-",
				Inline: false,
			})
			continue
		}

		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   event.Title,
			Value:  value,
			Inline: false,
		})
	}

	_, err = ds.ChannelMessageSendEmbed(m.ChannelID, em)
	if err != nil {
		log.Print(err)
	}

	return err
}
