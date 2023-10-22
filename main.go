package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <steam id>\n", os.Args[0])
		return
	}

	s := sprintCompletionistData(os.Args[1])

	// Should be a struct and should be marshalled to json. Dump to stdout,
	// because saving to file means that we can't pipe. We could also return
	// as a string and parameterize on the program args as a file/stdout, but
	// let's not overengineer right now.
	fmt.Println(s)
}

func sprintCompletionistData(profile string) string {
	values, err := scrapeCompletionistNodes(profile)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	// Originally, I tried to get the keys from a similar xpath, but it was
	// hard to keep the xpath order consistent with the values. It's also
	// possible to xpath into the "row" if you will (the key:value is a
	// sequence that doesn't share a parent), but then I had a lot of special
	// logic to parse lots of special rows that weren't consistent. Hard coding
	// the keys seems to be the simplest solution to keep the script
	// maintainable.
	keys := []string{
		"Achievements in Owned",
		"Daily Maximum",
		"Daily Average",
		"Average Completion",
		"Completion",
		"Median Completion",
		"Average in Progress",
		"Median in Progress",
		"Started Games Completed",
		"Daily Maximum Perfect Games",
		"Daily Average Perfect Games",
		"Total Playtime",
		"Average Playtime",
		"Median Playtime",
		"Completions",
		"Completions Average",
		"Completions Median",
		"Achievements",
		"Achievements to Perfection",
		"Achievements in Untouched",
		"Perfect Games",
		"Removed Perfect Games",
		"Games",
		"Games with Achievements",
		"Games in Progress",
		"Games Started",
		"Games Played",
		"Games Untouched",
		"Games Removed",
		"Games Expired",
		"Games Restricted",
	}

	stats := make(map[string]string)
	for i := range values {
		val := ""

		// The first value is "Achievements in Owned" which has an extra span within it.
		if i == 0 {
			val = values[0].FirstChild.NextSibling.FirstChild.Data
		} else {
			val = values[i].LastChild.Data
		}

		stats[keys[i]] = strings.TrimSpace(val)
	}

	str, err := json.MarshalIndent(stats, "", "  ")
	return string(str)
}

func scrapeCompletionistNodes(profile string) ([]*html.Node, error) {
	doc, err := htmlquery.LoadURL("https://completionist.me/steam/profile/" + profile)
	if err != nil {
		return nil, errorf("Couldn't load completionist profile: %s", err)
	}

	// This xpath is a bit gross, but rather than work on improving it, I'd
	// rather help implement an API on the website that returns the data
	// instead.
	values, err := htmlquery.QueryAll(doc, "/html/body/div[2]/main/div[1]/div/div[2]/div/div[1]/div/div/div/dl/dt/span|/html/body/div[2]/main/div[1]/div/div[2]/div/div[1]/div/div/div/dl/dt/a/span")
	if err != nil {
		return nil, errorf("Couldn't query for values: %s\n", err)
	}

	if len(values) == 0 {
		return nil, errorf("Couldn't find any stats. Check that your user id is valid.")
	}

	return values, nil
}

func errorf(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}
