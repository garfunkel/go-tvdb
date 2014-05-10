package main

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
	API_KEY = "DECE3B6B5464C552"
	GET_SERIES_URL = "http://thetvdb.com/api/GetSeries.php?seriesname=%v"
	GET_SERIES_BY_ID_URL = "http://thetvdb.com/api/%v/series/%v/en.xml"
	GET_DETAIL_URL = "http://thetvdb.com/api/%v/series/%v/all/en.xml"
	SEARCH_SERIES_URL = "http://thetvdb.com/?string=%v&searchseriesid=&tab=listseries&function=Search"
	SEARCH_SERIES_REGEX = `(?P<before><a href="/\?tab=series&amp;id=)(?P<seriesId>\d+)(?P<after>\&amp;lid=\d*">)`
)

var SearchSeriesRegex = regexp.MustCompile(SEARCH_SERIES_REGEX)

type PipeList []string

func (pipeList *PipeList) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) (err error) {
	content := ""

	if err = decoder.DecodeElement(&content, &start); err != nil {
		return err
	}

	*pipeList = strings.Split(strings.Trim(content, "|"), "|")

	return
}

type Episode struct {
	Id uint64 `xml:"id"`
	CombinedEpisodeNumber uint64 `xml:"Combined_episodenumber"`
	CombinedSeason uint64 `xml:"Combined_season"`
	DvdChapter string `xml:"DVD_chapter"`
	DvdDiscId string `xml:"DVD_discid"`
	DvdEpisodeNumber string `xml:"DVD_episodenumber"`
	DvdSeason string `xml:"DVD_season"`
	Director PipeList `xml:"Director"`
	EpImgFlag string `xml:"EpImgFlag"`
	EpisodeName string `xml:"EpisodeName"`
	EpisodeNumber int `xml:"EpisodeNumber"`
	FirstAired string `xml:"FirstAired"`
	GuestStars string `xml:"GuestStars"`
	ImdbId string `xml:"IMDB_ID"`
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
	SeasonId uint64 `xml:"seasonid"`
	SeriesId uint64 `xml:"seriesid"`
	ThumbAdded string `xml:"thumb_added"`
	ThumbHeight string `xml:"thumb_height"`
	ThumbWidth string `xml:"thumb_width"`
}

type Series struct {
	Id uint64 `xml:"id"`
	Actors PipeList `xml:"Actors"`
	AirsDayOfWeek string `xml:"Airs_DayOfWeek"`
	AirsTime string `xml:"Airs_Time"`
	ContentRating string `xml:"ContentRating"`
	FirstAired string `xml:"FirstAired"`
	Genre PipeList `xml:"Genre"`
	ImdbId string `xml:"IMDB_ID"`
	Language string `xml:"Language"`
	Network string `xml:"Network"`
	NetworkId string `xml:"NetworkID"`
	Overview string `xml:"Overview"`
	Rating string `xml:"Rating"`
	RatingCount string `xml:"RatingCount"`
	Runtime string `xml:"Runtime"`
	SeriesId string `xml:"SeriesID"`
	SeriesName string `xml:"SeriesName"`
	Status string `xml:"Status"`
	Added string `xml:"added"`
	AddedBy string `xml:"addedBy"`
	Banner string `xml:"banner"`
	Fanart string `xml:"fanart"`
	LastUpdated string `xml:"lastupdated"`
	Poster string `xml:"poster"`
	Zap2ItId string `xml:"zap2it_id"`
	Seasons map[uint64][]Episode
}

type SeriesList struct {
	Series []Series `xml:"Series"`
}

type EpisodeList struct {
	Episodes []Episode `xml:"Episode"`
}

func NewSeries(data []byte) (*Series, error) {
	series := Series{}

	series.Seasons = make(map[uint64][]Episode)

	if err := xml.Unmarshal(data, &series); err != nil {
		return nil, err
	}

	return &series, nil
}

func (seriesList *SeriesList) GetDetail() (err error) {
	for seriesIndex := range seriesList.Series {
		if err = seriesList.Series[seriesIndex].GetDetail(); err != nil {
			return
		}
	}

	return
}

func (series *Series) GetDetail() (err error) {
	response, err := http.Get(fmt.Sprintf(GET_DETAIL_URL, API_KEY, strconv.FormatUint(series.Id, 10)))

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
		series.Seasons = make(map[uint64][]Episode)
	}

	for _, episode := range episodeList.Episodes {
		series.Seasons[episode.SeasonNumber] = append(series.Seasons[episode.SeasonNumber], episode)
	}

	return
}

func GetSeries(name string) (seriesList SeriesList, err error) {
	response, err := http.Get(fmt.Sprintf(GET_SERIES_URL, url.QueryEscape(name)))

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

func GetSeriesById(id uint64) (series Series, err error) {
	response, err := http.Get(fmt.Sprintf(GET_SERIES_BY_ID_URL, API_KEY, id))

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

func SearchSeries(name string) (seriesList SeriesList, err error) {
	response, err := http.Get(fmt.Sprintf(SEARCH_SERIES_URL, url.QueryEscape(name)))

	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	groups := SearchSeriesRegex.FindAllSubmatch(buf, -1)

	for _, group := range groups {
		fmt.Println(group)
		seriesId := uint64(0)
		series := Series{}
		seriesId, err = strconv.ParseUint(string(group[2]), 10, 64)

		if err != nil {
			return
		}

		series, err = GetSeriesById(seriesId)

		if err != nil {
			return
		}

		seriesList.Series = append(seriesList.Series, series)
	}

	return
}
