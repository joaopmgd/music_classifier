package genres

import (
	"music_classifier/system"
)

// ChatGPT prompt
// I will give a JSON to you with a list of objects with the Artist field inside them. Those are electronic music artists and I would like to you to add a new field in each object of this list called Genre, in which you will describe which electronic music sub-genre the artist plays in a JSON list. If you don't know or never heard about the artist before, you can fill it out as UNKNOWN. Please answer in the same JSON format.
type ArtistGenres struct {
	Artist string   `json:"artist"`
	Genres []string `json:"genres"`
}

type ArtistGenresList []ArtistGenres

const (
	genresFile        = "genres.json"
	unknownMusicGenre = "UNKNOWN"
)

var genres ArtistGenresList

// init reads the genres.json file and unmarshal the data to the genres variable
func init() {
	err := system.ReadFileToData(genresFile, &genres)
	if err != nil {
		panic(err)
	}
	if len(genres) == 0 {
		panic("no genres found")
	}
}

// GetGenresFromArtist returns all genres for one artist
func GetGenresFromArtist(artist ...string) []string {
	finalGenres := []string{}
	for _, a := range genres {
		for _, artist := range artist {
			if a.Artist == artist {
				finalGenres = append(finalGenres, a.Genres...)
			}
		}
	}
	if len(finalGenres) == 0 {
		return []string{unknownMusicGenre}
	}
	return finalGenres
}
