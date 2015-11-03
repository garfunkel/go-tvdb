// Package tvdb provides a simple, sexy and easy golang module for TheTVDB.
package tvdb

import (
	"encoding/json"
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	// APIKey is the TheTVDB API key.
	//APIKey = "DECE3B6B5464C552"
	APIKey = "C53F79E0F7BEBD54"

	// API base URL.
	APIURL = "https://api-beta.thetvdb.com"

	// Login API URL.
	APILoginURL = APIURL + "/login"

	// Refresh token API URL.
	APIRefreshTokenURL = APIURL + "/refresh_token"

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

type Actor struct {
	ID uint64 `json:"id"`
	SeriesID uint64 `json:"seriesId"`
	Name string `json:"name"`
	Role pipeList `json:"role"`
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
	Data []Series `json:"data"`
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

func (tvdb *TheTVDB) RefreshToken() (err error) {
	// Login again if JWT has expired.
	if tvdb.jwt.Expired() {
		err = tvdb.Login()

		return
	}

	request, err := http.NewRequest("GET", APIRefreshTokenURL, nil)

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

	apiResponse := apiLoginResponse{}

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return
	}

	tvdb.jwt, err = DecodeJWT(apiResponse.JWT)

	return
}

func (tvdb *TheTVDB) Languages() (languages []Language, err error) {
	// Login again if JWT has expired.
	if tvdb.jwt.Expired() {
		err = tvdb.Login()
	// Refresh JWT if it is about to expire.
	} else if tvdb.jwt.AboutToExpire() {
		err = tvdb.RefreshToken()
	}

	if err != nil {
		return
	}

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
	// Login again if JWT has expired.
	if tvdb.jwt.Expired() {
		err = tvdb.Login()
	// Refresh JWT if it is about to expire.
	} else if tvdb.jwt.AboutToExpire() {
		err = tvdb.RefreshToken()
	}

	if err != nil {
		return
	}

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
	// Login again if JWT has expired.
	if tvdb.jwt.Expired() {
		err = tvdb.Login()
	// Refresh JWT if it is about to expire.
	} else if tvdb.jwt.AboutToExpire() {
		err = tvdb.RefreshToken()
	}

	if err != nil {
		return
	}

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

func (tvdb *TheTVDB) SearchSeries(params map[string]string) (seriesList []Series, err error) {
	// Login again if JWT has expired.
	if tvdb.jwt.Expired() {
		err = tvdb.Login()
	// Refresh JWT if it is about to expire.
	} else if tvdb.jwt.AboutToExpire() {
		err = tvdb.RefreshToken()
	}

	if err != nil {
		return
	}

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

	for _, series := range seriesList {
		series.tvdb = tvdb
	}

	return
}

// GetSeriesByID gets a TV series by ID.
func (tvdb *TheTVDB) GetSeriesByID(id uint64) (series Series, err error) {
	// Login again if JWT has expired.
	if tvdb.jwt.Expired() {
		err = tvdb.Login()
	// Refresh JWT if it is about to expire.
	} else if tvdb.jwt.AboutToExpire() {
		err = tvdb.RefreshToken()
	}

	if err != nil {
		return
	}

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
	series.tvdb = tvdb

	return
}
