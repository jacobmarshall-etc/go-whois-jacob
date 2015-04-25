package main

import (
    "fmt"
    "log"
    "net/http"
    "io"
    "encoding/json"
    "os"
    "strings"
)

var username string = os.Getenv("SPOTIFY_USERNAME")

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

func GetLastTrack(username string) (Track, error) {
    url := fmt.Sprintf("http://ajax.last.fm/user/%s/now", username)
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

    return r.Track, nil
}

func GetSpotify(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Access-Control-Allow-Origin", "*")

    track, err := GetLastTrack(username)

    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    data, _ := json.Marshal(track)

    io.WriteString(w, string(data))
}

func main() {
    http.HandleFunc("/spotify", GetSpotify)

    http.ListenAndServe(strings.Join(
        []string {":", os.Getenv("PORT")}, ""), nil)
}