package main

import (
    "fmt"
    "log"
    "net/http"
    "io"
    "encoding/json"
    "os"
    "strings"

    "github.com/bugsnag/bugsnag-go"
)

var username string = os.Getenv("SPOTIFY_USERNAME")

const (
    LAST_FM_URL = "http://last.fm"
    LAST_FM_USER_NOW_API_URL = "http://ajax.last.fm/user/%s/now"
)

type Artist struct {
    Name string `json:"name"`
    URL string `json:"url"`
}

type Track struct {
    Name string `json:"name"`
    Artist Artist `json:"artist"`
    URL string `json:"url"`
}

type SpotifyLastTrackResponse struct {
    Track Track
}

func prefixString(prefix string, str string) string {
    return strings.Join([]string {prefix, str}, "")
}

func GetLastTrack(username string) (Track, error) {
    url := fmt.Sprintf(LAST_FM_USER_NOW_API_URL, username)
    resp, err := http.Get(url)

    if err != nil {
        log.Printf("Error while fetching last track for user %s - %s", username, err)

        return Track{}, err
    }

    defer resp.Body.Close()

    r := new(SpotifyLastTrackResponse)
    err = json.NewDecoder(resp.Body).Decode(r)

    if err != nil {
        log.Printf("Error decoding response to response - %s", err)

        return Track{}, err
    }

    r.Track.URL = prefixString(LAST_FM_URL, r.Track.URL)
    r.Track.Artist.URL = prefixString(LAST_FM_URL, r.Track.Artist.URL)

    return r.Track, nil
}

func GetSpotify(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/last.fm", 301);
}

func GetLastFm(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Access-Control-Allow-Origin", "*")

    track, err := GetLastTrack(username)

    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)

        // Sent the bugsnag error after the response
        // Not sure if this has any effect of page load time, we'll see...
        bugsnag.Notify(err)
        return
    }

    data, _ := json.Marshal(track)

    io.WriteString(w, string(data))
}

func main() {
    bugsnag.Configure(bugsnag.Configuration{
        APIKey: os.Getenv("BUGSNAG_API_KEY"),
        ReleaseStage: os.Getenv("BUGSNAG_RELEASE_STAGE"),
    })

    http.HandleFunc("/spotify", bugsnag.HandlerFunc(GetSpotify))
    http.HandleFunc("/last.fm", bugsnag.HandlerFunc(GetLastFm))

    http.ListenAndServe(prefixString(":", os.Getenv("PORT")), nil)
}