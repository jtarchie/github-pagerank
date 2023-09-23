package rank

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"github.com/jtarchie/pagerank"
)

type Cmd struct {
	DBFilename string `required:"" default:":memory:" help:"the name of the file to save results to"`
}

func (c *Cmd) Run() error {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", c.DBFilename)
	if err != nil {
		return fmt.Errorf("could not open the database: %w", err)
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, `
	SELECT from_id, to_id FROM followers
	`)
	if err != nil {
		return fmt.Errorf("could not query followers: %w", err)
	}
	defer rows.Close()

	graph := pagerank.NewGraph[uint32]()

	var fromID, toID uint32
	for rows.Next() {
		err := rows.Scan(&fromID, &toID)
		if err != nil {
			return fmt.Errorf("could not load row: %w", err)
		}

		graph.Link(fromID, toID, 1)
	}

	type result struct {
		id   uint32
		rank float64
	}

	results := []result{}

	graph.Rank(0.85, 0.000001, func(id uint32, rank float64) {
		results = append(results, result{
			id:   id,
			rank: rank,
		})
	})

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].rank > results[j].rank
	})

	for index, result := range results {
		fmt.Printf("%d: id = %d, rank = %f\n", index, result.id, result.rank)
	}

	return nil
}
