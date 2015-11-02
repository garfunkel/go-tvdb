// Package tvdb provides a simple, sexy and easy golang module for TheTVDB.
package tvdb

import (
	"encoding/json"
	"encoding/base64"
	"errors"
	"fmt"
	"bytes"
	"time"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	// APIKey is the TheTVDB API key.
	//APIKey = "DECE3B6B5464C552"
	APIKey = "C53F79E0F7BEBD54"

	// API base URL.
	APIURL = "https://api-beta.thetvdb.com"

	// Login API URL.
	APILoginURL = APIURL + "/login"

	// Languages API URL.
	APILanguagesURL = APIURL + "/languages"

	// Language by ID API URL.
	APILanguageByIDURL = APILanguagesURL + "/%v"

	// Search series params API URL.
	APISearchSeriesParamsURL = APIURL + "/search/series/params"

	// Search series API URL.
	APISearchSeriesURL = APIURL + "/search/series"

	// Get series by ID API URL.
	APISeriesURL = APIURL + "/series/%v"

	// Get series actors API URL.
	APISeriesActorsURL = APISeriesURL + "/actors"

	// Get series images API URL.
	APISeriesImagesURL = APISeriesURL + "/images"

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
func (pipeList *PipeList) UnmarshalJSON(data []byte) (err error) {
	*pipeList = strings.Split(strings.Trim(string(data), "|"), "|")

	return
}

type jwtTime time.Time

func (t *jwtTime) UnmarshalJSON(data []byte) (err error) {
	unixTime, err := strconv.ParseInt(string(data), 10, 64)

	if err != nil {
		return
	}

	*t = jwtTime(time.Unix(unixTime, 0))

	return
}

type jwt struct {
	Header struct {
		Algorithm string `json:"alg"`
	}
	Claims struct {
		IssuedAt jwtTime `json:"orig_iat"`
		Expires jwtTime `json:"exp"`
		ID string `json:"id"`
	}
	Signature string
	JWT string
}

type TheTVDB struct {
	apiKey string
	jwt jwt
}

type apiLoginResponse struct {
	JWT	string `json:"token"`
}

type apiSearchSeriesParamsResponse struct {
	Data struct {
		Params []string `json:"params"`
	} `json:"data"`
}

// Series represents TV show on TheTVDB.
type Series struct {
	ID            uint64   `json:"id"`
	SeriesName    string   `json:"seriesName"`
	aliases []string `json:"aliases"`
	Banner        string   `json:"banner"`
	SeriesID      string   `json:"seriesID"`
	Status        string   `json:"status"`
	FirstAired    string   `json:"firstAired"`
	Network       string   `json:"network"`
	NetworkID     string   `json:"networkId"`
	Runtime       string   `json:"runtime"`
	Genres []string `json:"genre"`
	Overview      string   `json:"overview"`
	LastUpdated   jwtTime   `json:"lastUpdated"`
	AirsDayOfWeek string   `json:"airsDayOfWeek"`
	AirsTime      string   `json:"airsTime"`
	Rating        string   `json:"rating"`
	IMDbID        string   `json:"imdbId"`
	Zap2ItID      string   `json:"zap2itId"`
	Added         string   `json:"added"`
	AddedBy       string   `json:"addedBy"`
	siteRating int `json:"siteRating"`

	/*Actors        PipeList `xml:"Actors"`
	ContentRating string   `xml:"ContentRating"`
	Language      string   `xml:"Language"`
	RatingCount   string   `xml:"RatingCount"`
	Fanart        string   `xml:"fanart"`
	Poster        string   `xml:"poster"`
	Seasons       map[uint64][]*Episode*/
}

// SeriesList represents a list of TV shows, often used for returning search results.
type SeriesList []Series

func (series *Series) Actors() (actors []Actor, err error) {
	// Check JWT expiry.

	// Login if JWT expired.

	// Refresh JWT if it is about to expire.

	request, err := http.NewRequest("GET", fmt.Sprintf(APISeriesActorsURL, series.ID), nil)

	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tvdb.jwt.JWT))

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	apiResponse := apiSeriesActorsResponse{}

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return
	}

	actors = apiResponse.Data

	return
}

func (series *Series) Images() (images []Image, err error) {
	// Check JWT expiry.

	// Login if JWT expired.

	// Refresh JWT if it is about to expire.

	request, err := http.NewRequest("GET", fmt.Sprintf(APISeriesImagesURL, series.ID), nil)

	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tvdb.jwt.JWT))

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	apiResponse := apiSeriesImagesResponse{}

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return
	}

	images = []Image{apiResponse.Data}

	return
}

type Actor struct {
	ID uint64 `json:"id"`
	SeriesID uint64 `json:"seriesId"`
	Name string `json:"name"`
	Role PipeList `json:"role"`
	SortOrder uint64 `json:"sortOrder"`
	Image string `json:"image"`
	ImageAuthor uint64 `json:"imageAuthor"`
	ImageAdded string `json:"imageAdded"`
	LastUpdated string `json:"lastUpdated"`
}

type Image struct {
	FanArt uint64 `json:"fanart"`
	Poster uint64 `json:"poster"`
	Season uint64 `json:"season"`
	SeasonWide uint64 `json:"seasonwide"`
	Series uint64 `json:"series"`
}

type Language struct {
	ID uint64 `json:"id"`
	Abbreviation string `json:"abbreviation"`
	Name string `json:"name"`
	EnglishName string `json:"englishName"`
}

type apiSearchSeriesResponse struct {
	Data SeriesList `json:"data"`
}

type apiSeriesResponse struct {
	Data Series `json:"data"`
}

type apiSeriesActorsResponse struct {
	Data []Actor `json:"data"`
}

type apiSeriesImagesResponse struct {
	Data Image `json:"data"`
}

type apiLanguagesResponse struct {
	Data []Language `json:"data"`
}

type apiLanguageByIDResponse struct {
	Data Language `json:"data"`
}

func DecodeJWT(jwtStr string) (j jwt, err error) {
	fields := strings.Split(jwtStr, ".")

	if len(fields) != 3 {
		err = errors.New("Invalid JWT string")

		return
	}

	header, err := base64.StdEncoding.DecodeString(fields[0])

	if err != nil {
		return
	}

	claims, err := base64.StdEncoding.DecodeString(fields[1])

	if err != nil {
		return
	}

	signature, err := base64.RawURLEncoding.DecodeString(fields[2])

	if err != nil {
		return
	}

	j = jwt{
		JWT: jwtStr,
		Signature: string(signature),
	}

	err = json.Unmarshal(header, &j.Header)

	if err != nil {
		return
	}

	err = json.Unmarshal(claims, &j.Claims)

	return
}

func New(apiKey string) (tvdb *TheTVDB) {
	tvdb = &TheTVDB{
		apiKey: apiKey,
	}

	return
}

func (tvdb *TheTVDB) Login() (err error) {
	data := fmt.Sprintf(`{"apikey": "%s"}`, tvdb.apiKey)

	request, err := http.NewRequest("POST", APILoginURL, bytes.NewBufferString(data))

	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	apiResponse := apiLoginResponse{}

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return
	}

	tvdb.jwt, err = DecodeJWT(apiResponse.JWT)

	return
}

// Episode represents a TV show episode on TheTVDB.
type Episode struct {
	ID                    uint64   `xml:"id"`
	CombinedEpisodeNumber string   `xml:"Combined_episodenumber"`
	CombinedSeason        uint64   `xml:"Combined_season"`
	DvdChapter            string   `xml:"DVD_chapter"`
	DvdDiscID             string   `xml:"DVD_discid"`
	DvdEpisodeNumber      string   `xml:"DVD_episodenumber"`
	DvdSeason             string   `xml:"DVD_season"`
	Director              PipeList `xml:"Director"`
	EpImgFlag             string   `xml:"EpImgFlag"`
	EpisodeName           string   `xml:"EpisodeName"`
	EpisodeNumber         uint64   `xml:"EpisodeNumber"`
	FirstAired            string   `xml:"FirstAired"`
	GuestStars            string   `xml:"GuestStars"`
	ImdbID                string   `xml:"IMDB_ID"`
	Language              string   `xml:"Language"`
	Overview              string   `xml:"Overview"`
	ProductionCode        string   `xml:"ProductionCode"`
	Rating                string   `xml:"Rating"`
	RatingCount           string   `xml:"RatingCount"`
	SeasonNumber          uint64   `xml:"SeasonNumber"`
	Writer                PipeList `xml:"Writer"`
	AbsoluteNumber        string   `xml:"absolute_number"`
	Filename              string   `xml:"filename"`
	LastUpdated           string   `xml:"lastupdated"`
	SeasonID              uint64   `xml:"seasonid"`
	SeriesID              uint64   `xml:"seriesid"`
	ThumbAdded            string   `xml:"thumb_added"`
	ThumbHeight           string   `xml:"thumb_height"`
	ThumbWidth            string   `xml:"thumb_width"`
}

// EpisodeList represents a list of TV show episodes.
type EpisodeList struct {
	Episodes []*Episode `xml:"Episode"`
}

/*// GetDetail gets more detail for all TV shows in a list.
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

	// Unmarshal data into SeriesList object
	seriesList := SeriesList{}
	if err = xml.Unmarshal(data, &seriesList); err != nil {
		return
	}

	// Copy the first result into the series object
	*series = *seriesList.Series[0]

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
}*/

func (tvdb *TheTVDB) Languages() (languages []Language, err error) {
	// Check JWT expiry.

	// Login if JWT expired.

	// Refresh JWT if it is about to expire.

	request, err := http.NewRequest("GET", APILanguagesURL, nil)

	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tvdb.jwt.JWT))

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	apiResponse := apiLanguagesResponse{}

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return
	}

	languages = apiResponse.Data

	return
}

func (tvdb *TheTVDB) LanguageByID(id uint64) (language Language, err error) {
	// Check JWT expiry.

	// Login if JWT expired.

	// Refresh JWT if it is about to expire.

	request, err := http.NewRequest("GET", fmt.Sprintf(APILanguageByIDURL, id), nil)

	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tvdb.jwt.JWT))

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	apiResponse := apiLanguageByIDResponse{}

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return
	}

	language = apiResponse.Data

	return
}

func (tvdb *TheTVDB) SearchSeriesParams() (params []string, err error) {
	// Check JWT expiry.

	// Login if JWT expired.

	// Refresh JWT if it is about to expire.

	request, err := http.NewRequest("GET", APISearchSeriesParamsURL, nil)

	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tvdb.jwt.JWT))

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	apiResponse := apiSearchSeriesParamsResponse{}

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return
	}

	params = apiResponse.Data.Params

	return
}

func (tvdb *TheTVDB) SearchSeries(params map[string]string) (seriesList SeriesList, err error) {
	// Check JWT expiry.

	// Login if JWT expired.

	// Refresh JWT if it is about to expire.

	request, err := http.NewRequest("GET", APISearchSeriesURL, nil)

	if err != nil {
		return
	}

	query := request.URL.Query()

	for key, value := range params {
		query.Add(key, value)
	}
	
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tvdb.jwt.JWT))

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	apiResponse := apiSearchSeriesResponse{}

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return
	}

	seriesList = apiResponse.Data

	return
}

// GetSeriesByID gets a TV series by ID.
func (tvdb *TheTVDB) GetSeriesByID(id uint64) (series Series, err error) {
	// Check JWT expiry.

	// Login if JWT expired.

	// Refresh JWT if it is about to expire.

	request, err := http.NewRequest("GET", fmt.Sprintf(APISeriesURL, id), nil)

	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tvdb.jwt.JWT))

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	apiResponse := apiSeriesResponse{}

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return
	}

	series = apiResponse.Data

	return
}

/*// GetSeries gets a list of TV series by name, by performing a simple search.
func (tvdb *TheTVDB) GetSeries(name string) (seriesList SeriesList, err error) {
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
			// Some series can't be found, so we will ignore these.
			if _, ok := err.(*xml.SyntaxError); ok {
				err = nil

				continue
			} else {
				return
			}
		}

		seriesList.Series = append(seriesList.Series, series)
		doneSeriesIDs[seriesID] = struct{}{}

		if len(seriesList.Series) == maxResults {
			break
		}
	}

	return
}*/
