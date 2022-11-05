
```clickTT``` is a command line tool to fetch match dates from the official clickTT website.  
You can filter data and also get the output in different formats:

- tabular
- ICS (iCalender)
- JSON
- YAML

```clickTT``` is a more advanced version of what you can already do with a little shell snippet:

```
curl \
-s \
-X POST \
-d 'searchType=0&searchTimeRange=5&searchTimeRangeFrom=&searchTimeRangeTo=&selectedTeamId=WONoSelectionString&club=1416&searchMeetings=Suchen' \
https://ttbw.click-tt.de/cgi-bin/WebObjects/nuLigaTTDE.woa/wa/clubMeetings \
| html2text -style pretty -width 200 -ascii
```

## Compile

```
CGO_ENABLED=0 go build
```

### Examples

```
Usage: clickTT OPTIONS

OPTIONS

-G, --filter-group <string>   filter this group
-L, --filter-league <string>  filter this league
-a, --alarms                  add ALARMS (only in ICS output)
-c, --club <id>               club ID (default: 1416)
-g, --group <string>          show this group only (default: all)
-h, --help                    print help and exit
-i, --ics                     ICS output
-j, --json                    JSON output
-l, --league <string>         show this league only (default: all)
-t, --table                   TABLE output
-y, --yaml                    YAML output

EXAMPLES:

 Show all matches of the next 6 months in TABULAR format

  clickTT -t

Show all matches of the next 6 months in ICS format

  clickTT -i

Show all matches of the next 6 months in ICS format with alarms

  clickTT -i -a

Show only matches of "H KLA" league of the next 6 months  
in ICS format with alarms

  clickTT -i -a -l "H KLA"

Show all matches of the next 6 months  
in TABULAR format for club 1440

  clickTT -t -c 1440

Show all matches of the next 6 months  
in ICS format with alarms for club 1440

  clickTT -i -a -c 1440

Show all matches of the next 6 months  
in JSON format and pipe it to jq

  clickTT -j | jq .

Show only matches of "H KLA" league of the next 6 months  
in JSON format and pipe it to jq

  clickTT -j -l "H KLA"| jq .

Show all matches of the next 6 months in YAML format

  clickTT -y

Show only matches of "H KLA" league of the next 6 months  
in YAML format

  clickTT -y -l "H KLA"

Show only matches of group "TTG Bischweier" of the next 6 months  
in TABULAR format but not in the "J19 BK" league

  clickTT -t -L "J19 BK" -g "TTG Bischweier"

Show all matches of the next 6 months in TABULAR format  
but filter out the "TTC Muggensturm II" group

  clickTT -t -G "TTC Muggensturm II"

```
