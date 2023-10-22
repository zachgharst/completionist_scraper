package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func main() {
	f := flag.String("f", "", "file to send output to")
	norm := flag.Bool("n", false, "normalize data by removing commas, %, and other non-numeric characters")

	flag.Usage = func() {
		w := flag.CommandLine.Output()

		fmt.Fprintf(
			w,
			"Scrapes completionist.me for a profile's stats as a json object\nusage: %s <steam id>\n",
			os.Args[0],
		)

		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	nodes, err := scrapeCompletionistNodes(flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}

	s, err := sprintDataFromNodes(nodes, *norm)
	if err != nil {
		log.Fatalln(err)
	}

	if *f != "" {
		file, err := os.Create(*f)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		_, err = file.Write([]byte(s))
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		fmt.Println(s)
	}
}

func sprintDataFromNodes(nodes []*html.Node, normalize bool) (string, error) {
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

	stats := make(map[string]any)
	for i := range nodes {
		val := ""

		// The first value is "Achievements in Owned" which has an extra span within it.
		if i == 0 {
			val = nodes[0].FirstChild.NextSibling.FirstChild.Data
		} else {
			val = nodes[i].LastChild.Data
		}

		val = strings.TrimSpace(val)

		if normalize {
			nval, err := normalizeValue(val)
			if err != nil {
				log.Printf("couldn't normalize value: %s\n", err)
			}
			val = nval
		}

		stats[keys[i]] = val
	}

	str, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func normalizeValue(val string) (string, error) {
	// Three are three cases:
	//   1. The value is a percent, such as 12.34%
	//   2. The value is a large integer with commas, such as 12,345
	//   3. The value is a time, such as 12h 34m
	val = strings.ReplaceAll(val, ",", "")

	if strings.Contains(val, "%") {
		val = strings.ReplaceAll(val, " %", "")

		conv, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return "", errorf("couldn't parse %s as a float: %s", val, err)
		}

		val = fmt.Sprintf("%.4f", conv/100)
	}

	if strings.Contains(val, "h") || strings.Contains(val, "m") {
		hour, minute := 0, 0
		var err error

		if strings.Contains(val, "h") {
			split := strings.Split(val, " ")
			hourstr, _, _ := strings.Cut(split[0], "h")
			hour, err = strconv.Atoi(hourstr)

			if len(split) > 1 {
				val = split[1]
			} else {
				val = ""
			}
		}

		if strings.Contains(val, "m") {
			minutestr, _, _ := strings.Cut(val, "m")
			minute, err = strconv.Atoi(minutestr)
		}

		if err != nil {
			return "", errorf("couldn't parse %s as a timespan: %s", val, err)
		}

		val = fmt.Sprintf("%.2f", float64(hour)+float64(minute)/60)
	}

	return val, nil
}

func scrapeCompletionistNodes(profile string) ([]*html.Node, error) {
	doc, err := htmlquery.LoadURL("https://completionist.me/steam/profile/" + profile)
	if err != nil {
		return nil, errorf("couldn't load completionist profile: %s", err)
	}

	// This xpath is a bit gross, but rather than work on improving it, I'd
	// rather help implement an API on the website that returns the data
	// instead.
	xpath := "/html/body/div[2]/main/div[1]/div/div[2]/div/div[1]/div/div/div/dl/dt/span|"
	xpath += "/html/body/div[2]/main/div[1]/div/div[2]/div/div[1]/div/div/div/dl/dt/a/span"
	nodes, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		return nil, errorf("couldn't query for values: %s\n", err)
	}

	if len(nodes) == 0 {
		return nil, errorf("couldn't find any stats. Check that your user id is valid.")
	}

	return nodes, nil
}

func errorf(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}
