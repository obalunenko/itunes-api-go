package itunes_search_go

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obalunenko/itunes-api-go/option"
)

// Custom type that allows setting the func that our Mock Do func will run instead
type dofunc func(req *http.Request) (*http.Response, error)

// MockClient is the mock client
type mockHTTPClient struct {
	MockDo dofunc
}

// Overriding what the Do function should "do" in our MockClient
func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func equalITunesLookup(tb testing.TB, want, got LookupResponse) {
	tb.Helper()

	assert.Equal(tb, want.ResultCount, got.ResultCount)
	assert.Equal(tb, got.ResultCount, int64(1))

	wantres := want.Results[0]
	gotres := got.Results[0]

	assert.ElementsMatch(tb, wantres.ScreenshotURLs, gotres.ScreenshotURLs, "ScreenshotURLs")
	assert.ElementsMatch(tb, wantres.IpadScreenshotURLs, gotres.IpadScreenshotURLs, "IpadScreenshotURLs")
	assert.ElementsMatch(tb, wantres.AppletvScreenshotURLs, gotres.AppletvScreenshotURLs, "AppletvScreenshotURLs")
	assert.Equal(tb, wantres.ArtworkURL512, gotres.ArtworkURL512, "ArtworkURL512")
	assert.Equal(tb, wantres.ArtistViewURL, gotres.ArtistViewURL, "ArtistViewURL")
	assert.Equal(tb, wantres.ArtworkURL60, gotres.ArtworkURL60, "ArtworkURL60")
	assert.Equal(tb, wantres.ArtworkURL100, gotres.ArtworkURL100, "ArtworkURL100")
	assert.ElementsMatch(tb, wantres.SupportedDevices, gotres.SupportedDevices, "SupportedDevices")
	assert.ElementsMatch(tb, wantres.Advisories, gotres.Advisories, "Advisories")
	assert.Equal(tb, wantres.IsGameCenterEnabled, gotres.IsGameCenterEnabled, "IsGameCenterEnabled")
	assert.Equal(tb, wantres.Kind, gotres.Kind, "Kind")
	assert.ElementsMatch(tb, wantres.Features, gotres.Features, "Features")
	assert.ElementsMatch(tb, wantres.LanguageCodesISO2A, gotres.LanguageCodesISO2A, "LanguageCodesISO2A")
	assert.Equal(tb, wantres.FileSizeBytes, gotres.FileSizeBytes, "FileSizeBytes")
	assert.Equal(tb, wantres.SellerURL, gotres.SellerURL, "SellerURL")
	assert.Equal(tb, wantres.AverageUserRatingForCurrentVersion, gotres.AverageUserRatingForCurrentVersion,
		"AverageUserRatingForCurrentVersion")
	assert.Equal(tb, wantres.UserRatingCountForCurrentVersion, gotres.UserRatingCountForCurrentVersion,
		"UserRatingCountForCurrentVersion")
	assert.Equal(tb, wantres.TrackContentRating, gotres.TrackContentRating, "TrackContentRating")
	assert.Equal(tb, wantres.TrackCensoredName, gotres.TrackCensoredName, "TrackCensoredName")
	assert.Equal(tb, wantres.TrackViewURL, gotres.TrackViewURL, "TrackViewURL")
	assert.Equal(tb, wantres.ContentAdvisoryRating, gotres.ContentAdvisoryRating, "ContentAdvisoryRating")
	assert.Equal(tb, wantres.AverageUserRating, gotres.AverageUserRating, "AverageUserRating")
	assert.Equal(tb, wantres.TrackID, gotres.TrackID, "TrackID")
	assert.Equal(tb, wantres.TrackName, gotres.TrackName, "TrackName")
	assert.Equal(tb, wantres.ReleaseDate, gotres.ReleaseDate, "ReleaseDate")
	assert.ElementsMatch(tb, wantres.GenreIDS, gotres.GenreIDS, "GenreIDS")
	assert.Equal(tb, wantres.FormattedPrice, gotres.FormattedPrice, "FormattedPrice")
	assert.Equal(tb, wantres.PrimaryGenreName, gotres.PrimaryGenreName, "PrimaryGenreName")
	assert.Equal(tb, wantres.MinimumOSVersion, gotres.MinimumOSVersion, "MinimumOSVersion")
	assert.Equal(tb, wantres.IsVppDeviceBasedLicensingEnabled, gotres.IsVppDeviceBasedLicensingEnabled,
		"IsVppDeviceBasedLicensingEnabled")
	assert.Equal(tb, wantres.SellerName, gotres.SellerName, "SellerName")
	assert.Equal(tb, wantres.CurrentVersionReleaseDate, gotres.CurrentVersionReleaseDate,
		"CurrentVersionReleaseDate")
	assert.Equal(tb, wantres.ReleaseNotes, gotres.ReleaseNotes, "ReleaseNotes")
	assert.Equal(tb, wantres.PrimaryGenreID, gotres.PrimaryGenreID, "PrimaryGenreID")
	assert.Equal(tb, wantres.Currency, gotres.Currency, "Currency")
	assert.Equal(tb, wantres.Version, gotres.Version, "Version")
	assert.Equal(tb, wantres.WrapperType, gotres.WrapperType, "WrapperType")
	assert.Equal(tb, wantres.Description, gotres.Description, "Description")
	assert.Equal(tb, wantres.ArtistID, gotres.ArtistID, "ArtistID")
	assert.Equal(tb, wantres.ArtistName, gotres.ArtistName, "ArtistName")
	assert.ElementsMatch(tb, wantres.Genres, gotres.Genres, "Genres")
	assert.Equal(tb, wantres.Price, gotres.Price, "Price")
	assert.Equal(tb, wantres.BundleID, gotres.BundleID, "BundleID")
	assert.Equal(tb, wantres.UserRatingCount, gotres.UserRatingCount, "UserRatingCount")
}

func TestFetchAppDetails(t *testing.T) {
	ctx := context.Background()

	type args struct {
		id   int
		opts []option.LookupOption
	}

	tests := []struct {
		name       string
		args       args
		resultpath string
		wantErr    bool
	}{
		{
			name: "strava by details",
			args: args{
				id: 1068204657,
				opts: []option.LookupOption{
					option.WithCountry("by"),
				},
			},
			resultpath: filepath.Join("testdata", "id426826309.json"),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := ioutil.ReadFile(tt.resultpath)
			require.NoError(t, err)

			r := ioutil.NopCloser(bytes.NewReader(file))

			cl := client{
				httpclient: &mockHTTPClient{
					MockDo: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:           http.StatusText(http.StatusOK),
							StatusCode:       http.StatusOK,
							Proto:            "HTTP/1.0",
							ProtoMajor:       1,
							ProtoMinor:       0,
							Header:           nil,
							Body:             r,
							ContentLength:    0,
							TransferEncoding: nil,
							Close:            false,
							Uncompressed:     false,
							Trailer:          nil,
							Request:          nil,
							TLS:              nil,
						}, nil
					},
				},
			}

			expected, err := UnmarshalLookupResponse(file)
			require.NoError(t, err)

			got, err := cl.Lookup(ctx, tt.args.id, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			equalITunesLookup(t, expected, got)
		})
	}
}
