package skasaha

import (
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/bwmarrin/discordgo"
)

// cmdEvents handles the events command.
func (s *Skasaha) cmdSearch(ds *discordgo.Session, m *discordgo.MessageCreate, query string) error {
	var (
		err error
		id  int64

		results *bleve.SearchResult

		index = s.index.Index
		log   = s.logger()
	)

	bq := bleve.NewMatchQuery(query)
	search := bleve.NewSearchRequest(bq)

	results, err = index.Search(search)
	if err != nil {
		return err
	}

	if len(results.Hits) == 0 {
		_, err = ds.ChannelMessageSend(m.ChannelID, "No results.")

		return err
	}

	hit := results.Hits[0]
	rawID := hit.ID

	id, err = strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return err
	}

	embed := s.embeds[Snowflake(id)]

	_, err = ds.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Print(err)
	}

	return err
}
