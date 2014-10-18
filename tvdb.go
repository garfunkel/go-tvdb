// Package tvdb provides a simple, sexy and easy golang module for TheTVDB.
package tvdb

import (
	"fmt"
	"errors"
	"regexp"
	"net/url"
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/xml"
)

const (
	// APIKey is the TheTVDB API key.
	APIKey = "DECE3B6B5464C552"

	// GetSeriesURL is used to get basic series information by name.
	GetSeriesURL = "http://thetvdb.com/api/GetSeries.php?seriesname=%v"

	// GetSeriesByIDURL is used to get basic series information by ID.
	GetSeriesByIDURL = "http://thetvdb.com/api/%v/series/%v/en.xml"

	// GetSeriesByIMDBIDURL is used to get basic series information by IMDb ID.
	GetSeriesByIMDBIDURL = "http://thetvdb.com/api/GetSeriesByRemoteID.php?imdbid=%v"

	// GetDetailURL is used to get detailed series/episode information by ID.
	GetDetailURL = "http://thetvdb.com/api/%v/series/%v/all/en.xml"

	// SearchSeriesURL is used for series web searches.
	SearchSeriesURL = "http://thetvdb.com/?string=%v&searchseriesid=&tab=listseries&function=Search"

	// SearchSeriesRegexPattern is used for series web search matching.
	SearchSeriesRegexPattern = `(?P<before><a href="/\?tab=series&amp;id=)(?P<seriesId>\d+)(?P<after>\&amp;lid=\d*">)`
)

// SearchSeriesRegex is used for series web search matching.
var SearchSeriesRegex = regexp.MustCompile(SearchSeriesRegexPattern)

// PipeList type representing pipe-separated string values.
type PipeList []string

// UnmarshalXML unmarshals an XML element with string value into a pip-separated list of strings.
func (pipeList *PipeList) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) (err error) {
	content := ""

	if err = decoder.DecodeElement(&content, &start); err != nil {
		return err
	}

	*pipeList = strings.Split(strings.Trim(content, "|"), "|")

	return
}

// Episode represents a TV show episode on TheTVDB.
type Episode struct {
	ID uint64 `xml:"id"`
	CombinedEpisodeNumber string `xml:"Combined_episodenumber"`
	CombinedSeason uint64 `xml:"Combined_season"`
	DvdChapter string `xml:"DVD_chapter"`
	DvdDiscID string `xml:"DVD_discid"`
	DvdEpisodeNumber string `xml:"DVD_episodenumber"`
	DvdSeason string `xml:"DVD_season"`
	Director PipeList `xml:"Director"`
	EpImgFlag string `xml:"EpImgFlag"`
	EpisodeName string `xml:"EpisodeName"`
	EpisodeNumber uint64 `xml:"EpisodeNumber"`
	FirstAired string `xml:"FirstAired"`
	GuestStars string `xml:"GuestStars"`
	ImdbID string `xml:"IMDB_ID"`
	Language string `xml:"Language"`
	Overview string `xml:"Overview"`
	ProductionCode string `xml:"ProductionCode"`
	Rating string `xml:"Rating"`
	RatingCount string `xml:"RatingCount"`
	SeasonNumber uint64 `xml:"SeasonNumber"`
	Writer PipeList `xml:"Writer"`
	AbsoluteNumber string `xml:"absolute_number"`
	Filename string `xml:"filename"`
	LastUpdated string `xml:"lastupdated"`
	SeasonID uint64 `xml:"seasonid"`
	SeriesID uint64 `xml:"seriesid"`
	ThumbAdded string `xml:"thumb_added"`
	ThumbHeight string `xml:"thumb_height"`
	ThumbWidth string `xml:"thumb_width"`
}

// Series represents TV show on TheTVDB.
type Series struct {
	ID uint64 `xml:"id"`
	Actors PipeList `xml:"Actors"`
	AirsDayOfWeek string `xml:"Airs_DayOfWeek"`
	AirsTime string `xml:"Airs_Time"`
	ContentRating string `xml:"ContentRating"`
	FirstAired string `xml:"FirstAired"`
	Genre PipeList `xml:"Genre"`
	ImdbID string `xml:"IMDB_ID"`
	Language string `xml:"Language"`
	Network string `xml:"Network"`
	NetworkID string `xml:"NetworkID"`
	Overview string `xml:"Overview"`
	Rating string `xml:"Rating"`
	RatingCount string `xml:"RatingCount"`
	Runtime string `xml:"Runtime"`
	SeriesID string `xml:"SeriesID"`
	SeriesName string `xml:"SeriesName"`
	Status string `xml:"Status"`
	Added string `xml:"added"`
	AddedBy string `xml:"addedBy"`
	Banner string `xml:"banner"`
	Fanart string `xml:"fanart"`
	LastUpdated string `xml:"lastupdated"`
	Poster string `xml:"poster"`
	Zap2ItID string `xml:"zap2it_id"`
	Seasons map[uint64][]*Episode
}

// SeriesList represents a list of TV shows, often used for returning search results.
type SeriesList struct {
	Series []*Series `xml:"Series"`
}

// EpisodeList represents a list of TV show episodes.
type EpisodeList struct {
	Episodes []*Episode `xml:"Episode"`
}

// GetDetail gets more detail for all TV shows in a list.
func (seriesList *SeriesList) GetDetail() (err error) {
	for seriesIndex := range seriesList.Series {
		if err = seriesList.Series[seriesIndex].GetDetail(); err != nil {
			return
		}
	}

	return
}

// GetDetail gets more detail for a TV show, including information on it's episodes.
func (series *Series) GetDetail() (err error) {
	response, err := http.Get(fmt.Sprintf(GetDetailURL, APIKey, strconv.FormatUint(series.ID, 10)))

	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	if err = xml.Unmarshal(data, series); err != nil {
		return
	}

	episodeList := EpisodeList{}

	if err = xml.Unmarshal(data, &episodeList); err != nil {
		return
	}

	if series.Seasons == nil {
		series.Seasons = make(map[uint64][]*Episode)
	}

	for _, episode := range episodeList.Episodes {
		series.Seasons[episode.SeasonNumber] = append(series.Seasons[episode.SeasonNumber], episode)
	}

	return
}

// GetSeries gets a list of TV series by name, by performing a simple search.
func GetSeries(name string) (seriesList SeriesList, err error) {
	response, err := http.Get(fmt.Sprintf(GetSeriesURL, url.QueryEscape(name)))

	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	err = xml.Unmarshal(data, &seriesList)

	return
}

// GetSeriesByID gets a TV series by ID.
func GetSeriesByID(id uint64) (series *Series, err error) {
	response, err := http.Get(fmt.Sprintf(GetSeriesByIDURL, APIKey, id))

	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	seriesList := SeriesList{}

	if err = xml.Unmarshal(data, &seriesList); err != nil {
		return
	}

	if len(seriesList.Series) != 1 {
		err = errors.New("incorrect number of series")

		return
	}

	series = seriesList.Series[0]

	return
}

// GetSeriesByIMDBID gets series from IMDb's ID.
func GetSeriesByIMDBID(id string) (series *Series, err error) {
	response, err := http.Get(fmt.Sprintf(GetSeriesByIMDBIDURL, id))

	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	seriesList := SeriesList{}

	if err = xml.Unmarshal(data, &seriesList); err != nil {
		return
	}

	if len(seriesList.Series) != 1 {
		err = errors.New("incorrect number of series")

		return
	}

	series = seriesList.Series[0]

	return
}

// SearchSeries searches for TV shows by name, using the more sophisticated
// search on TheTVDB's homepage. This is the recommended search method.
func SearchSeries(name string, maxResults int) (seriesList SeriesList, err error) {
	response, err := http.Get(fmt.Sprintf(SearchSeriesURL, url.QueryEscape(name)))

	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	groups := SearchSeriesRegex.FindAllSubmatch(buf, -1)
	doneSeriesIDs := make(map[uint64]struct{})

	for _, group := range groups {
		seriesID := uint64(0)
		var series *Series
		seriesID, err = strconv.ParseUint(string(group[2]), 10, 64)

		if _, ok := doneSeriesIDs[seriesID]; ok {
			continue
		}

		if err != nil {
			return
		}

		series, err = GetSeriesByID(seriesID)

		if err != nil {
			continue
		}

		seriesList.Series = append(seriesList.Series, series)
		doneSeriesIDs[seriesID] = struct{}{}

		if len(seriesList.Series) == maxResults {
			break
		}
	}

	return
}
