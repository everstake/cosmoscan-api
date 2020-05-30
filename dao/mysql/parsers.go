package mysql

import (
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (m DB) GetParsers() (parsers []dmodels.Parser, err error) {
	q := squirrel.Select("*").From(dmodels.ParsersTable)
	err = m.find(&parsers, q)
	if err != nil {
		return nil, err
	}
	return parsers, nil
}

func (m DB) GetParser(title string) (parser dmodels.Parser, err error) {
	q := squirrel.Select("*").From(dmodels.ParsersTable).
		Where(squirrel.Eq{"par_title": title})
	err = m.first(&parser, q)
	return parser, err
}

func (m DB) UpdateParser(parser dmodels.Parser) error {
	q := squirrel.Update(dmodels.ParsersTable).
		Where(squirrel.Eq{"par_id": parser.ID}).
		SetMap(map[string]interface{}{
			"par_height": parser.Height,
		})
	return m.update(q)
}
