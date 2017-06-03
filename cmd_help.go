package skasaha

import (
	"github.com/bwmarrin/discordgo"
)

const (
	helpURL = "https://github.com/KuroiKitsu/skasaha#commands"
)

// cmdEvents handles the help command.
func (s *Skasaha) cmdHelp(ds *discordgo.Session, m *discordgo.MessageCreate, topic string) error {
	var (
		err error

		log = s.logger()
	)

	em := &discordgo.MessageEmbed{
		Title:       "Help",
		URL:         helpURL,
		Description: "Inline help is unavailable at the moment. Please visit the repository for more information.",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    SkasahaAvatarImageURL,
			Width:  SkasahaAvatarImageWidth,
			Height: SkasahaAvatarImageHeight,
		},
	}

	_, err = ds.ChannelMessageSendEmbed(m.ChannelID, em)
	if err != nil {
		log.Print(err)
	}

	return err
}
