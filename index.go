package skasaha

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

type Index struct {
	Mapping *mapping.IndexMappingImpl
	Index   bleve.Index
}

func (s *Skasaha) rebuildIndex() error {
	var (
		err error

		m     *mapping.IndexMappingImpl
		index bleve.Index

		log = s.logger()
	)

	log.Print("rebuilding index")

	m = bleve.NewIndexMapping()

	index, err = bleve.NewMemOnly(m)
	if err != nil {
		return err
	}

	for id, embed := range s.embeds {
		if embed.Title != "" {
			log.Printf("indexing %#v", embed.Title)

			index.Index(id.String(), embed.Title)
		}
	}

	s.index = &Index{
		Mapping: m,
		Index:   index,
	}

	log.Print("index updated")

	return nil
}
