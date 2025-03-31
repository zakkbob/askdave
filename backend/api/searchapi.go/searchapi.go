package searchapi

import (
	"context"
	"log"
	"net/http"

	"github.com/ZakkBob/AskDave/backend/orm"
	"github.com/gin-gonic/gin"
)

type Search struct {
	Query string `form:"q"`
}

type Result struct {
	SiteName    string `json:"site_name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func Init(router *gin.Engine, prefix string) {
	router.GET(prefix+"/", GetSearchResults)
}

func GetSearchResults(c *gin.Context) {
	var search Search
	err := c.ShouldBind(&search)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, err := orm.DbPool().Query(context.Background(), `SELECT url, sitename, title, description
												FROM
												(
													SELECT CONCAT(site.url, page.path), page.og_sitename, page.title, page.og_description, 
													(LOWER(CONCAT(site.url, page.path)) LIKE LOWER($1))::int + 
													(LOWER(page.og_sitename) LIKE LOWER($1))::int +
													(LOWER(page.title) LIKE LOWER($1))::int +
													(LOWER(page.og_description) LIKE LOWER($1))::int +
													(LOWER(page.og_title) LIKE LOWER($1))::int
													FROM page 
													LEFT JOIN site 
													ON page.site = site.id
												) results(url, sitename, title, description, rank)
												WHERE rank > 0
												ORDER BY rank DESC
												LIMIT 20;`, "%"+search.Query+"%")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result Result
	var results []Result

	for rows.Next() {
		err := rows.Scan(&result.URL, &result.SiteName, &result.Title, &result.Description)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		results = append(results, result)
	}

	log.Println(search.Query)
	c.IndentedJSON(http.StatusOK, results)
}
