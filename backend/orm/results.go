package orm

import (
	"github.com/ZakkBob/AskDave/gocommon/tasks"

	"fmt"
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

		if err != nil { // should probably log this
			continue
		}

		p.Validator = *robotsResult.Validator

		err = p.Save()
		if err != nil {
			return fmt.Errorf("failed to save robots result: %w", err)
		}
	}

	for _, pageResult := range r.Pages {
		p, err := PageByUrl(pageResult.Url)
		if err != nil {
			return fmt.Errorf("failed to get page by url: %w", err)
		}

		// will be nil if failed
		if pageResult.Page != nil {
			p.Page = *pageResult.Page
		}

		p.ScheduleNextCrawl(pageResult.Changed)
		p.Assigned = false

		err = p.Save()
		if err != nil {
			return fmt.Errorf("failed to save page result: %w", err)
		}

		// err = p.SaveCrawl(time.Now(), pageResult.Success, pageResult.FailureReason, pageResult.Changed, pageResult.Page.Hash)
		// if err != nil {
		// 	return fmt.Errorf("unable to save page result: %w", err)
		// }
	}

	return nil
}
