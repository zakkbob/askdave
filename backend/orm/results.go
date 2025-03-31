package orm

import (
	"log"
	"time"

	"github.com/ZakkBob/AskDave/gocommon/tasks"
)

func SaveResults(r *tasks.Results) error {
	for _, robotsResult := range r.Robots {
		if !robotsResult.Changed {
			continue
		}

		if !robotsResult.Success {
			continue // temporary - should probably do something
		}

		p, err := SiteByUrl(robotsResult.Url)

		if err != nil {
			log.Printf("failed to get site with url '%s' : %v\n", robotsResult.Url.String(), err)
			continue
		}

		p.Validator = *robotsResult.Validator
		p.LastRobotsCrawl = time.Now()

		err = p.Save()
		if err != nil {
			log.Printf("failed to save robots result: %v\n", err)
			continue
		}
	}

	for _, pageResult := range r.Pages {
		p, err := PageByUrl(pageResult.Url)
		if err != nil {
			log.Printf("failed to get page by url: %v\n", err)
			continue
		}

		// will be nil if failed
		if pageResult.Page != nil {
			p.Page = *pageResult.Page
		}

		p.ScheduleNextCrawl(pageResult.Changed)
		p.Assigned = false

		err = p.Save()
		if err != nil {
			log.Printf("failed to save page result: %v\n", err)
			continue
		}

		// err = p.SaveCrawl(time.Now(), pageResult.Success, pageResult.FailureReason, pageResult.Changed, pageResult.Page.Hash)
		// if err != nil {
		// 	return fmt.Errorf("unable to save page result: %w", err)
		// }
	}

	return nil
}
