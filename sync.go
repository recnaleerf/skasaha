package skasaha

import (
	"sort"
	"strconv"

	"github.com/KuroiKitsu/go-gbf"
	"github.com/KuroiKitsu/go-gbf/scraper"
	"github.com/bwmarrin/discordgo"
)

func (s *Skasaha) Sync() error {
	s.initEmbeds()

	var (
		err error
	)

	err = s.SyncEvents()
	if err != nil {
		return err
	}

	err = s.SyncCharacters()
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

func (s *Skasaha) SyncCharacters() error {
	var (
		err error

		scrapedCharacters  []*scraper.Character

		characters = make([]gbf.Character, 0)
		log    = s.logger()
	)

	log.Print("synchronizing characters")

	scrapedCharacters, err = scraper.Characters()
	if err != nil {
		return err
	}

	for _, basicChar := range scrapedCharacters {
		var (
			char   gbf.Character
			details *scraper.CharacterDetails
		)

		char.URL = basicChar.URL
		char.Number = basicChar.Number

		switch basicChar.Tier {
		case "R":
			char.Tier = gbf.R
		case "SR":
			char.Tier = gbf.SR
		case "SSR":
			char.Tier = gbf.SSR
		}

		char.Name = basicChar.Name

		switch basicChar.Element {
		case "Fire":
			char.Element = gbf.Fire
		case "Water":
			char.Element = gbf.Water
		case "Earth":
			char.Element = gbf.Earth
		case "Wind":
			char.Element = gbf.Wind
		case "Light":
			char.Element = gbf.Light
		case "Dark":
			char.Element = gbf.Dark
		}

		switch basicChar.Style {
		case "Attack":
			char.Style = gbf.Attack
		case "Balanced":
			char.Style = gbf.Balanced
		case "Defense":
			char.Style = gbf.Defense
		case "Heal":
			char.Style = gbf.Heal
		case "Special":
			char.Style = gbf.Special
		}

		switch basicChar.Race {
		case "None":
			char.Race = gbf.NoRace
		case "Draph":
			char.Race = gbf.Draph
		case "Erune":
			char.Race = gbf.Erune
		case "Harvin":
			char.Race = gbf.Harvin
		case "Human":
			char.Race = gbf.Human
		case "Primal":
			char.Race = gbf.Primal
		}

		switch basicChar.Sex {
		case "Other":
			char.Sex = gbf.OtherSex
		case "Male":
			char.Sex = gbf.Male
		case "Female":
			char.Sex = gbf.Female
		}

		char.Stars = basicChar.Stars
		char.HP = basicChar.HP
		char.Attack = basicChar.ATK
		char.ExtendedMastery = basicChar.EM

		switch basicChar.Weapon {
		case "Sabre":
			char.WeaponType = gbf.Sabre
		case "Dagger":
			char.WeaponType = gbf.Dagger
		case "Spear":
			char.WeaponType = gbf.Spear
		case "Axe":
			char.WeaponType = gbf.Axe
		case "Staff":
			char.WeaponType = gbf.Staff
		case "Gun":
			char.WeaponType = gbf.Gun
		case "Melee":
			char.WeaponType = gbf.Melee
		case "Bow":
			char.WeaponType = gbf.Bow
		case "Harp":
			char.WeaponType = gbf.Harp
		case "Katana":
			char.WeaponType = gbf.Katana
		case "Boost":
			char.WeaponType = gbf.Boost
		}

		if char.URL != "" {
			log.Printf("fetching details for %#v", char.Name)

			details, err = scraper.CharacterDetailsURL(char.URL)
			if err != nil {
				return err
			}

			char.Description = details.Description

			if len(details.ImagesArtURL) > 0 {
				char.ImageArtURL = details.ImagesArtURL[0]
			}

			if len(details.ImagesSpriteURL) > 0 {
				char.ImageSpriteURL = details.ImagesSpriteURL[0]
			}
		}

		characters = append(characters, char)
	}

	if len(s.embedsCharacters) > 0 {
		for _, id := range s.embedsCharacters {
			delete(s.embeds, id)
		}
	}

	s.Characters = characters

	s.embedsCharacters = make([]Snowflake, 0, len(s.Characters))
	for _, char := range s.Characters {
		id := NewSnowflake()

		em := &discordgo.MessageEmbed{
			Title:       char.Name,
			Description: char.Description,
			URL:         char.URL,
			Provider: &discordgo.MessageEmbedProvider{
				Name: "Granblue Fantasy Wiki",
				URL:  scraper.WikiHomeURL.String(),
			},
			Fields: make([]*discordgo.MessageEmbedField, 0),
		}

		if char.ImageSpriteURL == "" {
			em.Thumbnail = &discordgo.MessageEmbedThumbnail{
				URL:    GBFLogoImageURL,
				Width:  GBFLogoImageWidth,
				Height: GBFLogoImageHeight,
			}
		} else {
			em.Thumbnail = &discordgo.MessageEmbedThumbnail{
				URL:    char.ImageSpriteURL,
			}
		}

		if char.ImageArtURL != "" {
			em.Image = &discordgo.MessageEmbedImage{
				URL: char.ImageArtURL,
			}
		}

		if char.Number > 0 {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name:   "No",
				Value:  strconv.Itoa(char.Number),
				Inline: true,
			})
		}

		if char.Tier != gbf.UnknownTier {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name:   "Tier",
				Value:  char.Tier.String(),
				Inline: true,
			})
		}

		if char.Element != gbf.UnknownElement {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name: "Element",
				Value: char.Element.String(),
				Inline: true,
			})
		}

		if char.Style != gbf.UnknownStyle {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name: "Style",
				Value: char.Style.String(),
				Inline: true,
			})
		}

		if char.Race != gbf.UnknownRace {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name: "Race",
				Value: char.Race.String(),
				Inline: true,
			})
		}

		if char.Sex != gbf.UnknownSex {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name: "Sex",
				Value: char.Sex.String(),
				Inline: true,
			})
		}

		if char.Stars > 0 {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name: "Stars",
				Value: strconv.Itoa(char.Stars),
				Inline: true,
			})
		}

		if char.HP > 0 {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name: "HP",
				Value: strconv.Itoa(char.HP),
				Inline: true,
			})
		}

		if char.Attack > 0 {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name: "Attack",
				Value: strconv.Itoa(char.Attack),
				Inline: true,
			})
		}

		var extendedMasteryStr string
		if char.ExtendedMastery {
			extendedMasteryStr = "Yes"
		} else {
			extendedMasteryStr = "No"
		}

		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name: "Extended Mastery",
			Value: extendedMasteryStr,
			Inline: true,
		})

		if char.WeaponType != gbf.UnknownWeaponType {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name: "Weapon Type",
				Value: char.WeaponType.String(),
				Inline: true,
			})
		}

		s.embeds[id] = em
		s.embedsCharacters = append(s.embedsCharacters, id)
	}

	log.Printf("%d characters found", len(s.Characters))

	return err
}
