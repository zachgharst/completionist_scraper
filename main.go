package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
)

func main() {
    doc, err := htmlquery.LoadURL("https://completionist.me/steam/profile/76561197997651901")
    if err != nil {
        log.Fatalf("Couldn't load completionist profile: %s", err)
    }

    values, err := htmlquery.QueryAll(doc, "/html/body/div[2]/main/div[1]/div/div[2]/div/div[1]/div/div/div/dl/dt/span|/html/body/div[2]/main/div[1]/div/div[2]/div/div[1]/div/div/div/dl/dt/a/span")
    if err != nil {
        log.Fatalf("Couldn't find values: %s", err)
    }

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

    fmt.Println("{")
    for i := range values {
        val := ""
        if i == 0 {
            val = values[0].FirstChild.NextSibling.FirstChild.Data
        } else {
            val = values[i].LastChild.Data
        }

        val = strings.TrimSpace(val)

        if i != len(values) - 1 {
            fmt.Printf("  \"%s\": \"%s\",\n", keys[i], val)
        } else {
            fmt.Printf("  \"%s\": \"%s\"\n", keys[i], val)
        }
    }
    fmt.Println("}")
}
