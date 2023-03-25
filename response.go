package itunes_search_go

import (
	"encoding/json"
	"errors"
)

var (
	// ErrNoResults returned when itunes api response contains no results.
	ErrNoResults = errors.New("no results")
	// ErrImpreciseResults returned when itunes api response contains more than 1 result for exact app id.
	ErrImpreciseResults = errors.New("imprecise results")
)

func UnmarshalLookupResponse(data []byte) (LookupResponse, error) {
	var res LookupResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return LookupResponse{}, err
	}

	return res, nil
}

func (l *LookupResponse) Marshal() ([]byte, error) {
	return json.Marshal(l)
}

// LookupResponse represents response of itunes api lookup endpoint.
type LookupResponse struct {
	// ResultCount contains info about total found results number.
	ResultCount int64 `json:"resultCount"`
	// Results is an array with all found results.
	Results []Result `json:"results"`
}

// Result takes care of check the results count and shortcut to get the first result from array (lookup always
// contains 1 result in array).
func (l LookupResponse) Result() (Result, error) {
	if len(l.Results) == 0 {
		return Result{}, ErrNoResults
	}

	if len(l.Results) > 1 {
		return l.Results[0], ErrImpreciseResults
	}

	return l.Results[0], nil
}

// Result represents result entity from itunes api lookup.
type Result struct {
	ScreenshotURLs                     []string      `json:"screenshotUrls"`
	IpadScreenshotURLs                 []string      `json:"ipadScreenshotUrls"`
	AppletvScreenshotURLs              []interface{} `json:"appletvScreenshotUrls"`
	ArtworkURL60                       string        `json:"artworkUrl60"`
	ArtworkURL512                      string        `json:"artworkUrl512"`
	ArtworkURL100                      string        `json:"artworkUrl100"`
	ArtistViewURL                      string        `json:"artistViewUrl"`
	SupportedDevices                   []string      `json:"supportedDevices"`
	Advisories                         []string      `json:"advisories"`
	IsGameCenterEnabled                bool          `json:"isGameCenterEnabled"`
	Kind                               string        `json:"kind"`
	Features                           []string      `json:"features"`
	TrackCensoredName                  string        `json:"trackCensoredName"`
	LanguageCodesISO2A                 []string      `json:"languageCodesISO2A"`
	FileSizeBytes                      string        `json:"fileSizeBytes"`
	SellerURL                          string        `json:"sellerUrl"`
	ContentAdvisoryRating              string        `json:"contentAdvisoryRating"`
	AverageUserRatingForCurrentVersion float64       `json:"averageUserRatingForCurrentVersion"`
	UserRatingCountForCurrentVersion   int64         `json:"userRatingCountForCurrentVersion"`
	AverageUserRating                  float64       `json:"averageUserRating"`
	TrackViewURL                       string        `json:"trackViewUrl"`
	TrackContentRating                 string        `json:"trackContentRating"`
	TrackID                            int64         `json:"trackId"`
	TrackName                          string        `json:"trackName"`
	ReleaseDate                        string        `json:"releaseDate"`
	GenreIDS                           []string      `json:"genreIds"`
	FormattedPrice                     string        `json:"formattedPrice"`
	PrimaryGenreName                   string        `json:"primaryGenreName"`
	IsVppDeviceBasedLicensingEnabled   bool          `json:"isVppDeviceBasedLicensingEnabled"`
	MinimumOSVersion                   string        `json:"minimumOsVersion"`
	SellerName                         string        `json:"sellerName"`
	CurrentVersionReleaseDate          string        `json:"currentVersionReleaseDate"`
	ReleaseNotes                       string        `json:"releaseNotes"`
	PrimaryGenreID                     int64         `json:"primaryGenreId"`
	Version                            string        `json:"version"`
	WrapperType                        string        `json:"wrapperType"`
	Currency                           string        `json:"currency"`
	Description                        string        `json:"description"`
	ArtistID                           int64         `json:"artistId"`
	ArtistName                         string        `json:"artistName"`
	Genres                             []string      `json:"genres"`
	Price                              float64       `json:"price"`
	BundleID                           string        `json:"bundleId"`
	UserRatingCount                    int64         `json:"userRatingCount"`
}
