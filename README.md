
```clickTT``` is a command line tool to fetch match dates from the official clickTT website.  
You can filter data and also get the output in different formats:

- tabular
- ICS (iCalender)
- JSON
- YAML

```clickTT``` is a more advanced version what you can already do with a little shell snippet:

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
