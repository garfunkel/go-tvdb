package tvdb

// Episode represents a TV show episode on TheTVDB.
type Episode struct {
	ID                 uint64   `json:"id"`
	AiredSeason        uint64   `json:"airedSeason"`
	AiredEpisodeNumber uint64   `json:"airedEpisodeNumber"`
	EpisodeName        string   `json:"episodeName"`
	FirstAired         string   `json:"firstAired"`
	GuestStars         []string `json:"guestStars"`
	Director           pipeList `json:"director"`
	Writer             []string `json:"writers"`
	Overview           string   `json:"overview"`
	Language           struct {
		EpisodeName string `json:"episodeName"`
		Overview    string `json:"overview"`
	} `json:"language"`
	ProductionCode    string   `json:"productionCode"`
	ShowURL           string   `json:"showUrl"`
	LastUpdated       unixTime `json:"lastUpdated"`
	DVDDiscID         string   `json:"dvdDiscid"`
	DVDSeason         uint64   `json:"dvdSeason"`
	DVDEpisodeNumber  uint64   `json:"dvdEpisodeNumber"`
	DVDChapter        string   `json:"dvdChapter"`
	AbsoluteNumber    uint64   `json:"absoluteNumber"`
	Filename          string   `json:"filename"`
	SeriesID          uint64   `json:"seriesId"`
	LastUpdatedBy     uint64   `json:"json:"lastUpdatedBy"`
	AirsAfterSeason   string   `json:"airsAfterSeason"`
	AirsBeforeSeason  string   `json:"airsBeforeSeason"`
	AirsBeforeEpisode string   `json:"airsBeforeEpisode"`
	ThumbAuthor       uint64   `json:"thumbAuthor"`
	ThumbAdded        string   `json:"thumbAdded"`
	ThumbWidth        string   `json:"thumbWidth"`
	ThumbHeight       string   `json:"thumbHeight"`
	IMDbID            string   `json:"imdbId"`
	SiteRating        float64  `json:"siteRating"`
	tvdb              *TheTVDB
}
