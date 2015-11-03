package tvdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Series represents TV show on TheTVDB.
type Series struct {
	ID            uint64   `json:"id"`
	SeriesName    string   `json:"seriesName"`
	aliases       []string `json:"aliases"`
	Banner        string   `json:"banner"`
	SeriesID      string   `json:"seriesID"`
	Status        string   `json:"status"`
	FirstAired    string   `json:"firstAired"`
	Network       string   `json:"network"`
	NetworkID     string   `json:"networkId"`
	Runtime       string   `json:"runtime"`
	Genres        []string `json:"genre"`
	Overview      string   `json:"overview"`
	LastUpdated   unixTime `json:"lastUpdated"`
	AirsDayOfWeek string   `json:"airsDayOfWeek"`
	AirsTime      string   `json:"airsTime"`
	Rating        string   `json:"rating"`
	IMDbID        string   `json:"imdbId"`
	Zap2ItID      string   `json:"zap2itId"`
	Added         string   `json:"added"`
	AddedBy       string   `json:"addedBy"`
	siteRating    int      `json:"siteRating"`
	tvdb          *TheTVDB
}

func (series *Series) Images() (images []Image, err error) {
	// Login again if JWT has expired.
	if series.tvdb.jwt.Expired() {
		err = series.tvdb.Login()
		// Refresh JWT if it is about to expire.
	} else if series.tvdb.jwt.AboutToExpire() {
		err = series.tvdb.RefreshToken()
	}

	if err != nil {
		return
	}

	request, err := http.NewRequest("GET", fmt.Sprintf(APISeriesImagesURL, series.ID), nil)

	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", series.tvdb.jwt.JWT))

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

func (series *Series) Actors() (actors []Actor, err error) {
	// Login again if JWT has expired.
	if series.tvdb.jwt.Expired() {
		err = series.tvdb.Login()
		// Refresh JWT if it is about to expire.
	} else if series.tvdb.jwt.AboutToExpire() {
		err = series.tvdb.RefreshToken()
	}

	if err != nil {
		return
	}

	request, err := http.NewRequest("GET", fmt.Sprintf(APISeriesActorsURL, series.ID), nil)

	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", series.tvdb.jwt.JWT))

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
