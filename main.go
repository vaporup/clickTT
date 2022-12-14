package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/galdor/go-cmdline"
	"github.com/go-yaml/yaml"
)

type Event struct {
	Day    string
	Date   string
	Time   string
	League string
	Home   string
	Guest  string
}

func main() {

	// CLI

	cl := cmdline.New()

	cl.AddFlag("7", "last-seven-days", "Events of last 7 days")
	cl.AddFlag("t", "table", "TABLE output")
	cl.AddFlag("i", "ics", "ICS output")
	cl.AddFlag("j", "json", "JSON output")
	cl.AddFlag("y", "yaml", "YAML output")
	cl.AddFlag("a", "alarms", "add ALARMS (only in ICS output)")
	cl.AddOption("l", "league", "string", "show this league only")
	cl.SetOptionDefault("l", "all")
	cl.AddOption("L", "filter-league", "string", "filter this league")
	cl.AddOption("g", "group", "string", "show this group only")
	cl.SetOptionDefault("g", "all")
	cl.AddOption("G", "filter-group", "string", "filter this group")
	cl.AddOption("c", "club", "id", "club ID")
	cl.SetOptionDefault("c", "1416")
	cl.Parse(os.Args)

	if len(os.Args) == 1 {

		fmt.Println("")

		cl.PrintUsage(os.Stderr)

		fmt.Fprintf(os.Stderr,
			"\nEXAMPLES:\n")

		fmt.Fprintf(os.Stderr,
			"\n Show all matches of the next 6 months in TABULAR format\n")
		fmt.Fprintf(os.Stderr, "\n  %s -t\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow all matches of the next 6 months in ICS format\n")
		fmt.Fprintf(os.Stderr, "\n  %s -i\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow all matches of the next 6 months in ICS format with alarms\n")
		fmt.Fprintf(os.Stderr, "\n  %s -i -a\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow only matches of \"H KLA\" league of the next 6 months in ICS format with alarms\n")
		fmt.Fprintf(os.Stderr, "\n  %s -i -a -l \"H KLA\"\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow all matches of the next 6 months in TABULAR format for club 1440\n")
		fmt.Fprintf(os.Stderr, "\n  %s -t -c 1440\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow all matches of the next 6 months in ICS format with alarms for club 1440\n")
		fmt.Fprintf(os.Stderr, "\n  %s -i -a -c 1440\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow all matches of the next 6 months in JSON format and pipe it to jq\n")
		fmt.Fprintf(os.Stderr, "\n  %s -j | jq .\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow only matches of \"H KLA\" league of the next 6 months in JSON format and pipe it to jq\n")
		fmt.Fprintf(os.Stderr, "\n  %s -j -l \"H KLA\"| jq .\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow all matches of the next 6 months in YAML format\n")
		fmt.Fprintf(os.Stderr, "\n  %s -y\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow only matches of \"H KLA\" league of the next 6 months in YAML format\n")
		fmt.Fprintf(os.Stderr, "\n  %s -y -l \"H KLA\"\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow only matches of group \"TTG Bischweier\" of the next 6 months in TABULAR format but not in the \"J19 BK\" league\n")
		fmt.Fprintf(os.Stderr, "\n  %s -t -L \"J19 BK\" -g \"TTG Bischweier\"\n", cl.ProgramName)

		fmt.Fprintf(os.Stderr,
			"\nShow all matches of the next 6 months in TABULAR format but filter out the \"TTC Muggensturm II\" group\n")
		fmt.Fprintf(os.Stderr, "\n  %s -t -G \"TTC Muggensturm II\"\n", cl.ProgramName)

		os.Exit(1)
	}

	// List to store all events

	events := []Event{}

	// HTTP params

	params := url.Values{}

	params.Add("searchType", "0")
	params.Add("club", cl.OptionValue("c"))
	//params.Add("club", "1416")
	if cl.IsOptionSet("7") {
		params.Add("searchTimeRange", "-1") // last 7 days
	} else {
		params.Add("searchTimeRange", "5") // next 6 months
	}

	resp, err := soup.PostForm("https://ttbw.click-tt.de/cgi-bin/WebObjects/nuLigaTTDE.woa/wa/clubMeetings", params)

	if err != nil {
		log.Fatal(err)
	}

	doc := soup.HTMLParse(resp)
	trs := doc.Find("table", "class", "result-set").FindAll("tr")

	var day string
	var date string

	var lenType int
	var lenHome int
	var lenGuest int
	var lenths int

	for _, tr := range trs {

		storeEvent := false

		e := Event{}

		tds := tr.FindAll("td")
		ths := tr.FindAll("th")

		if len(ths) > 0 {
			lenths = len(ths)
		}

		if lenths == 9 {
			// normal table length

			for td_num, td := range tds {

				text := td.Text()
				attrs := td.Attrs()

				data := strings.TrimSpace(text)
				data = strings.Split(data, "\n")[0]

				if td_num == 0 {
					if attrs["class"] == "tabelle-rowspan" {
						e.Day = strings.Trim(day, ".")
					} else {
						day = data
						e.Day = strings.Trim(day, ".")
					}
				}

				if td_num == 1 {
					if attrs["class"] == "tabelle-rowspan" {
						e.Date = date
					} else {
						date = data
						e.Date = date
					}
				}

				if td_num == 2 {
					e.Time = data
				}

				if td_num == 3 {
					continue
				}

				if td_num == 4 {
					if len(data) > lenType {
						lenType = len(data)
					}
					e.League = data
				}

				if td_num == 5 {
					if len(data) > lenHome {
						lenHome = len(data)
					}
					e.Home = data
				}

				if td_num == 6 {
					if len(data) > lenGuest {
						lenGuest = len(data)
					}
					e.Guest = data
				}

				if td_num > 6 {
					continue
				}

			}
		}

		if lenths == 10 {
			// e.g. with "POKAL" events
			// the table has one column more (some number)

			for td_num, td := range tds {

				text := td.Text()
				attrs := td.Attrs()

				data := strings.TrimSpace(text)
				data = strings.Split(data, "\n")[0]

				if td_num == 0 {
					if attrs["class"] == "tabelle-rowspan" {
						e.Day = strings.Trim(day, ".")
					} else {
						day = data
						e.Day = strings.Trim(day, ".")
					}
				}

				if td_num == 1 {
					if attrs["class"] == "tabelle-rowspan" {
						e.Date = date
					} else {
						date = data
						e.Date = date
					}
				}

				if td_num == 2 {
					e.Time = data
				}

				if td_num == 3 {
					continue
				}

				if td_num == 4 {
					continue
				}

				if td_num == 5 {
					if len(data) > lenType {
						lenType = len(data)
					}
					e.League = data
				}

				if td_num == 6 {
					if len(data) > lenHome {
						lenHome = len(data)
					}
					e.Home = data
				}

				if td_num == 7 {
					if len(data) > lenGuest {
						lenGuest = len(data)
					}
					e.Guest = data
				}

				if td_num > 7 {
					continue
				}

			}
		}

		// Skip empty events
		if (e != Event{}) {

			if cl.OptionValue("l") == "all" && cl.OptionValue("g") == "all" {
				storeEvent = true
			}

			if cl.OptionValue("l") == e.League && cl.OptionValue("g") == "all" {
				storeEvent = true
			}

			if cl.OptionValue("l") == "all" && (cl.OptionValue("g") == e.Home || cl.OptionValue("g") == e.Guest) {
				storeEvent = true
			}

			if cl.OptionValue("l") == e.League && (cl.OptionValue("g") == e.Home || cl.OptionValue("g") == e.Guest) {
				storeEvent = true
			}

			if cl.OptionValue("L") == e.League || (cl.OptionValue("G") == e.Home || cl.OptionValue("G") == e.Guest) {
				storeEvent = false
			}

			if storeEvent {
				events = append(events, e)

			}

		}
	}

	// TABLE OUTPUT
	if cl.IsOptionSet("t") {

		fmt.Printf("%-4v ", "TAG")
		fmt.Printf("%-12v ", "DATUM")
		fmt.Printf("%-7v ", "ZEIT")
		fmt.Printf("%-*v ", lenType+2, "LIGA")
		fmt.Printf("%-*v ", lenHome+2, "HEIM")
		fmt.Printf("%-*v \n", lenGuest+2, "GAST")

		for _, event := range events {

			fmt.Printf("%-4v ", event.Day)
			fmt.Printf("%-12v ", event.Date)
			fmt.Printf("%-7v ", event.Time)
			fmt.Printf("%-*v ", lenType+2, event.League)
			fmt.Printf("%-*v ", lenHome+2, event.Home)
			fmt.Printf("%-*v \n", lenGuest+2, event.Guest)

		}

		os.Exit(0)
	}

	// ICS OUTPUT
	if cl.IsOptionSet("i") {

		fmt.Println("BEGIN:VCALENDAR")
		fmt.Println("VERSION:2.0")
		fmt.Println("PRODID:-//vaporup//NONSGML clickTTermine//EN")
		fmt.Println("")

		for idx, event := range events {

			fmt.Println("")
			fmt.Println("BEGIN:VEVENT")
			fmt.Printf("UID:uid%d@clickTTermine\n", idx)
			fmt.Println("SUMMARY:", event.Home, "-", event.Guest, "("+event.League+")")
			//fmt.Println("DESCRIPTION:", event.Home, "-", event.Guest, event.Date+" "+event.Time, "("+event.League+")")

			const layout = "02.01.2006 15:04"
			tm, _ := time.Parse(layout, event.Date+" "+event.Time)
			dtstart := tm.Format("20060102T150405")

			fmt.Printf("DTSTART:%v\n", dtstart)

			// WITH ALARMS
			if cl.IsOptionSet("a") {
				fmt.Println("BEGIN:VALARM")
				fmt.Println("TRIGGER:-PT180M")
				fmt.Println("ACTION:DISPLAY")
				fmt.Println("DESCRIPTION:Reminder")
				fmt.Println("END:VALARM")

				fmt.Println("BEGIN:VALARM")
				fmt.Println("TRIGGER:-P2D")
				fmt.Println("ACTION:DISPLAY")
				fmt.Println("DESCRIPTION:Reminder")
				fmt.Println("END:VALARM")
			}

			fmt.Println("END:VEVENT")

		}

		fmt.Println("")
		fmt.Println("END:VCALENDAR")

		os.Exit(0)
	}

	// JSON OUTPUT
	if cl.IsOptionSet("j") {
		j, err := json.Marshal(events)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(j))
		os.Exit(0)
	}

	// YAML OUTPUT
	if cl.IsOptionSet("y") {
		y, err := yaml.Marshal(events)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(y))
		os.Exit(0)
	}

}
