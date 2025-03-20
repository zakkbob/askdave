package orm

import (
	"github.com/ZakkBob/AskDave/gocommon/tasks"

	"context"
	"fmt"
)

func SaveResults(r *tasks.Results) error {
	robotsQuery := `UPDATE robots 
		SET (allowed_patterns, disallowed_patterns, last_crawl) = ($1, $2, CURRENT_DATE) 
		FROM site WHERE site.url = $3 
		AND site.id = robots.site_id;`

	for urlS, robotsResult := range r.Robots {
		if !robotsResult.Changed {
			continue
		}

		if !robotsResult.Success {
			continue // temporary - should probably do something
		}

		_, err := dbpool.Exec(context.Background(), robotsQuery, robotsResult.Validator.AllowedStrings(), robotsResult.Validator.DisallowedStrings(), urlS)
		if err != nil {
			return fmt.Errorf("unable to save robots result: %v", err)
		}
	}

	for urlS, pageResult := range r.Pages {
		// Add crawl to db
		if pageResult.Success {
			if pageResult.Changed {
				// add crawl
				query := `INSERT INTO crawl (page_id, datetime, success, content_changed, title, og_title, og_description, hash)
					SELECT page.id, CURRENT_TIMESTAMP, TRUE, TRUE, $1, $2, $3, $4 
					FROM page 
					WHERE url = $5;`

				_, err := dbpool.Exec(context.Background(), query, pageResult.Page.Title, pageResult.Page.OgTitle, pageResult.Page.OgDescription, pageResult.Page.Hash, urlS)
				if err != nil {
					return fmt.Errorf("unable to save page result: %v", err)
				}

				// delete existing links
				query = `DELETE FROM link 
					USING page 
					WHERE link.src = page.id 
					AND page.url = $1;`
				_, err = dbpool.Exec(context.Background(), query, urlS)
				if err != nil {
					return fmt.Errorf("unable to delete links for url '%s': %v", urlS, err)
				}

				//add new links
				query = `INSERT INTO link (src, dst, count) 
					SELECT page.id, $2, $3 
					FROM page 
					WHERE page.url=$1;` //src, dst, count
				for _, dst := range pageResult.Page.Links {
					_, err := dbpool.Exec(context.Background(), query, urlS, dst.String(), 1)
					if err != nil {
						return fmt.Errorf("unable to add link from '%s' to '%s': %v", urlS, dst.String(), err)
					}
				}
			} else {
				query := `INSERT INTO crawl (page_id, datetime, success, content_changed)
					SELECT page.id, CURRENT_TIMESTAMP, TRUE, FALSE 
					FROM page 
					WHERE url = $1;`

				_, err := dbpool.Exec(context.Background(), query, urlS)
				if err != nil {
					return fmt.Errorf("unable to save page result: %v", err)
				}
			}
		} else {
			query := `INSERT INTO crawl (page_id, datetime, success)
				SELECT page.id, CURRENT_TIMESTAMP, FALSE 
				FROM page 
				WHERE url = $1;`

			_, err := dbpool.Exec(context.Background(), query, urlS)

			if err != nil {
				return fmt.Errorf("unable to save page result: %v", err)
			}
		}

		// Set next crawl date
		const maxCrawlInterval = 30
		const minCrawlInterval = 1

		var query string

		if pageResult.Changed {
			query = `UPDATE page SET 
				interval_delta = LEAST(interval_delta - 1, -1), 
				crawl_interval = LEAST(GREATEST(crawl_interval + interval_delta, $1),$2), 
				next_crawl = next_crawl + interval '1' day * crawl_interval 
				WHERE url = $3;`
		} else {
			query = `UPDATE page SET 
				interval_delta = GREATEST(interval_delta + 1, 1), 
				crawl_interval = LEAST(GREATEST(crawl_interval + interval_delta, $1),$2), 
				next_crawl = next_crawl + interval '1' day * crawl_interval 
				WHERE url = $3;`
		}

		_, err := dbpool.Exec(context.Background(), query, maxCrawlInterval, minCrawlInterval, urlS)
		if err != nil {
			return fmt.Errorf("unable to update page crawl scheduling: %v", err)
		}
	}

	return nil
}
