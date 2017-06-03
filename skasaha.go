package skasaha

import (
	"errors"
	"strings"

	"github.com/KuroiKitsu/go-gbf"
	"github.com/bwmarrin/discordgo"
)

type Skasaha struct {
	Prefix   string
	Token    string
	EmojiDir string

	Logger     Logger
	Events     []gbf.Event
	Characters []gbf.Character

	session *discordgo.Session

	index            *Index
	embeds           map[Snowflake]*discordgo.MessageEmbed
	embedsEvents     []Snowflake
	embedsCharacters []Snowflake
}

func (s *Skasaha) initEmbeds() {
	if s.embeds == nil {
		s.embeds = make(map[Snowflake]*discordgo.MessageEmbed)
	}

	if s.embedsEvents == nil {
		s.embedsEvents = make([]Snowflake, 0)
	}
}

func (s *Skasaha) logger() Logger {
	if s.Logger == nil {
		return quietLoggerSingleton
	}

	return s.Logger
}

// OnMessageCreate is the event handler for discordgo library.
func (s *Skasaha) OnMessageCreate(ds *discordgo.Session, m *discordgo.MessageCreate) {
	if s.Prefix == "" {
		return
	}
	if !strings.HasPrefix(m.Content, s.Prefix) {
		return
	}

	var (
		err error

		head string
		tail string

		log = s.logger()
	)

	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) == 0 {
		return
	}

	head = strings.TrimPrefix(parts[0], s.Prefix)
	if len(parts) > 1 {
		tail = parts[1]
	}

	if tail == "" && strings.HasPrefix(head, s.Prefix) {
		tail = strings.TrimPrefix(head, s.Prefix)
		head = "emoji"

		if tail == "" {
			head = ""
		}
	}

	// Empty command is an alias for `search`.
	if head == "" && tail != "" {
		head = "search"
	}

	log.Printf("head=%#v tail=%#v", head, tail)

	switch head {
	case "h":
		fallthrough
	case "help":
		err = s.cmdHelp(ds, m, tail)

	case "s":
		fallthrough
	case "search":
		err = s.cmdSearch(ds, m, tail)

	case "events":
		err = s.cmdEvents(ds, m)

	case "emo":
		fallthrough
	case "emoji":
		err = s.cmdEmoji(ds, m, tail)
	}

	if err != nil {
		_, _ = ds.ChannelMessageSend(m.ChannelID, err.Error())
	}
}

// Open starts a connection with Discord.
func (s *Skasaha) Open() error {
	if s.session != nil {
		return errors.New("session already opened")
	}

	var (
		err error
		ds  *discordgo.Session

		log = s.logger()
	)

	ds, err = discordgo.New("Bot " + s.Token)
	if err != nil {
		return err
	}

	ds.AddHandler(s.OnMessageCreate)

	log.Print("opening session")

	err = ds.Open()
	if err != nil {
		return err
	}

	log.Print("session opened successfully")

	s.session = ds

	return nil
}

// Close terminates the connection with Discord.
func (s *Skasaha) Close() error {
	if s.session == nil {
		return errors.New("no session")
	}

	log := s.logger()
	log.Print("closing session")

	err := s.session.Close()

	return err
}
