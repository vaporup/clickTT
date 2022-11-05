
## Shell Example

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
