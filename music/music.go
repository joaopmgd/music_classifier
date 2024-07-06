package music

import (
	"fmt"
	"music_classifier/artist"
	"music_classifier/configuration"
	"music_classifier/genres"
	"music_classifier/system"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wtolson/go-taglib"
)

type MusicList []*Music

// Music is the struct that represents a music data
type Music struct {
	Title              string
	Artists            []string
	Genres             []string
	BaseDirectory      string
	Path               string
	OrderedArtistNames string
}

const (
	// Directories to save the music
	artistsDirectory = "/Artists/"
	genreDirectory   = "/Genre/"
)

var (
	artistPopularityMap = artist.PopularityMap{}
)

// SaveMusicPerArtist saves the music per artist
func (m Music) SaveMusicPerArtist() error {
	// Link the file to the artist directory
	for _, artist := range m.Artists {
		artistDirectory := fmt.Sprintf("%s/Artists/%s", configuration.GetMusicsPathDirectory(), artist)
		err := os.MkdirAll(artistDirectory, os.ModePerm)
		if err != nil {
			return err
		}
		artistsFilePath := fmt.Sprintf("%s/%s - %s.mp3", artistDirectory, m.OrderedArtistNames, m.Title)
		err = system.LinkFile(m.Path, artistsFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveMusicPerGenre saves the music per genre
func (m Music) SaveMusicPerGenre() error {
	// Copy the file to the Genre directory
	for _, genre := range m.Genres {
		genreDirectory := fmt.Sprintf("%s/Genre/%s", configuration.GetMusicsPathDirectory(), genre)
		err := os.MkdirAll(genreDirectory, os.ModePerm)
		if err != nil {
			return err
		}
		genreFilePath := fmt.Sprintf("%s/%s - %s.mp3", genreDirectory, m.OrderedArtistNames, m.Title)
		err = system.LinkFile(m.Path, genreFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// setOrderedArtistNames sets a string with the artist names ordered by popularity
func (m *Music) setOrderedArtistNames() {
	// Get the artist with the highest popularity
	artistPopularityList := artist.PopularityList{}
	for _, artist := range m.Artists {
		artistPopularityList = append(artistPopularityList, *artistPopularityMap[artist])
	}
	// Sort the artist popularity list
	sort.Sort(artistPopularityList)
	m.OrderedArtistNames = strings.Join(artistPopularityList.GetArtistNameList(), ", ")
}

// setFilePath sets the file path
func (m *Music) setFilePath() error {
	// Rename the file with the music title - artist name
	savedFilePath := fmt.Sprintf("%s/%s - %s.mp3", m.BaseDirectory, m.OrderedArtistNames, m.Title)
	err := os.Rename(m.Path, savedFilePath)
	if err != nil {
		return err
	}
	m.Path = savedFilePath
	return nil
}

// GetMusicData returns a list of music data
func GetMusicData() (MusicList, error) {
	var m MusicList
	var err error
	switch configuration.GetMusicNameOrigin() {
	case configuration.MetadataOrigin:
		m, err = getMusicDataFromMetadata()
	case configuration.FileNameOrigin:
		m, err = getMusicDataFromFileName()
	default:
		return nil, fmt.Errorf("invalid music name origin, %s", configuration.GetMusicNameOrigin())
	}
	if err != nil {
		return nil, err
	}
	// Validate if the music data is not empty
	if len(m) == 0 {
		return nil, fmt.Errorf("no music data found")
	}
	// Set the ordered artist names and the renamed file path
	for _, music := range m {
		music.setOrderedArtistNames()
		err = music.setFilePath()
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// Split is a function to split the string by the rune
func Split(r rune) bool {
	return r == '/' || r == ','
}

// getMusicDataFromFileName returns a list of music data from the file name
func getMusicDataFromFileName() (MusicList, error) {
	musicList := MusicList{}
	for _, fileStringPath := range configuration.GetFilePathList() {
		// It will not classify the music if it is in the artist directory
		if strings.Contains(filepath.Dir(fileStringPath), artistsDirectory) ||
			strings.Contains(filepath.Dir(fileStringPath), genreDirectory) {
			continue
		}

		// Create the music data from the file name
		artistsString, titleString, valid := strings.Cut(filepath.Base(fileStringPath), "-")
		if !valid {
			if configuration.ShouldPrintErrorLogs() {
				fmt.Printf("Invalid file name: %s\n", fileStringPath)
			}
			continue
		}
		// Sanitize the artist names
		artists := strings.FieldsFunc(artistsString, Split)
		for i, artistData := range artists {
			// Remove trailing spaces
			artists[i] = strings.TrimSpace(artistData)
		}
		// Sanitize the music title
		title := strings.TrimSuffix(titleString, filepath.Ext(titleString))
		musicData := &Music{
			Title:         title,
			Artists:       artists,
			Genres:        genres.GetGenresFromArtist(artists...),
			BaseDirectory: filepath.Dir(fileStringPath),
			Path:          fileStringPath,
		}

		// Save the music in the map of artists
		for _, artistData := range musicData.Artists {
			_, ok := artistPopularityMap[artistData]
			if !ok {
				artistPopularityMap[artistData] = &artist.Popularity{
					Artist: artistData,
					Count:  0,
				}
			}
			if strings.Contains(musicData.Title, artistData) {
				artistPopularityMap[artistData].Count++
			}
			artistPopularityMap[artistData].Count++
		}

		// Save the music metadata
		musicMetadata, err := taglib.Read(fileStringPath)
		if err != nil {
			return nil, err
		}
		musicMetadata.SetTitle(musicData.Title)
		musicMetadata.SetArtist(strings.Join(artists, "/"))
		musicMetadata.SetGenre(strings.Join(genres.GetGenresFromArtist(artists...), ", "))
		err = musicMetadata.Save()
		if err != nil {
			return nil, err
		}

		// Save the music in the music list
		musicList = append(musicList, musicData)
	}
	return musicList, nil
}

// getMusicDataFromMetadata returns a list of music data from the metadata
func getMusicDataFromMetadata() (MusicList, error) {
	musicList := MusicList{}
	for _, fileStringPath := range configuration.GetFilePathList() {

		// It will not classify the music if it is in the artist directory
		if strings.Contains(filepath.Dir(fileStringPath), artistsDirectory) ||
			strings.Contains(filepath.Dir(fileStringPath), genreDirectory) {
			continue
		}

		// Load music
		musicData, err := taglib.Read(fileStringPath)
		if err != nil {
			return nil, err
		}

		// Ignore music without artist or title
		if musicData.Artist() == "" || musicData.Title() == "" {
			if configuration.ShouldPrintErrorLogs() {
				fmt.Printf("Invalid metadata: %s\n", fileStringPath)
			}
			continue
		}

		// Save the music in the map of artists
		for _, artistData := range strings.Split(musicData.Artist(), "/") {
			_, ok := artistPopularityMap[artistData]
			if !ok {
				artistPopularityMap[artistData] = &artist.Popularity{
					Artist: artistData,
					Count:  0,
				}
			}
			if strings.Contains(musicData.Title(), artistData) {
				artistPopularityMap[artistData].Count++
			}
			artistPopularityMap[artistData].Count++
		}

		// Save the music genre
		musicData.SetGenre(strings.Join(genres.GetGenresFromArtist(strings.Split(musicData.Artist(), "/")...), ", "))
		err = musicData.Save()
		if err != nil {
			return nil, err
		}

		musicTitle := musicData.Title()
		// Updates musics with / in the title
		if strings.Contains(musicTitle, "/") {
			// remove everything after the / or - in the title
			musicTitle = strings.Split(musicTitle, "/")[0]
			musicTitle = strings.Split(musicTitle, "-")[0]
		}

		// Save the music in the music list
		musicList = append(musicList, &Music{
			Title:         musicTitle,
			Artists:       strings.Split(musicData.Artist(), "/"),
			Genres:        genres.GetGenresFromArtist(strings.Split(musicData.Artist(), "/")...),
			BaseDirectory: filepath.Dir(fileStringPath),
			Path:          fileStringPath,
		})
		musicData.Close()
	}
	return musicList, nil
}

// PrintArtistPopularity prints the artist popularity
func PrintArtistPopularity() {
	artistPopularityList := artistPopularityMap.ToList()
	sort.Sort(artistPopularityList)
	for _, artist := range artistPopularityList {
		fmt.Printf("%s: %d\n", artist.Artist, artist.Count)
	}
}
