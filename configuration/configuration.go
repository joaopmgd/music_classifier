package configuration

import (
	"fmt"
	"music_classifier/system"
)

type configuration struct {
	// The path to the musics directory
	MusicsPathDirectory string `json:"music_path_directory"`
	// MusicNameOrigin is the origin of the music name
	MusicNameOrigin MusicNameOrigin `json:"music_name_origin"`
	// SaveMusicPerArtist is a flag to save the music per artist
	SaveMusicPerArtist bool `json:"save_music_per_artist"`
	// SaveMusicPerGenre is a flag to save the music per genre
	SaveMusicPerGenre bool `json:"save_music_per_genre"`
	// PrintArtistPopularity is a flag to print the artist popularity
	PrintArtistPopularity bool `json:"print_artist_popularity"`
	// PrintErrorLogs is a flag to print the error logs
	PrintErrorLogs bool `json:"print_error_logs"`

	// FilePathList is a list of all the files in the musics directory
	// Not used in the configuration file, but it is used to store the file list
	FilePathList []string `json:"file_path_list"`
}

const (
	configFile = "config.json"
)

// The type of music name origin
type MusicNameOrigin string

const (
	// MetadataOrigin is the origin of the music name from the metadata
	MetadataOrigin MusicNameOrigin = "METADATA"
	// FileNameOrigin is the origin of the music name from the file name
	FileNameOrigin MusicNameOrigin = "FILE_NAME"
)

var config configuration

// init reads the configuration file and unmarshal the data to the config variable
func init() {
	err := system.ReadFileToData(configFile, &config)
	if err != nil {
		panic(err)
	}
	// Reads the file list recursively
	filePathList, err := system.ReadFileList(config.MusicsPathDirectory)
	if err != nil {
		panic(err)
	}
	config.FilePathList = filePathList
	if err := config.IsValid(); err != nil {
		panic(err)
	}
}

// IsValid checks if the configuration is valid
func (c *configuration) IsValid() error {
	if c.MusicsPathDirectory == "" {
		return fmt.Errorf("musics path directory is empty")
	}
	if len(c.FilePathList) == 0 {
		return fmt.Errorf("file path list is empty")
	}
	if err := c.MusicNameOrigin.IsValid(); err != nil {
		return err
	}
	return nil
}

// ShouldSaveMusicPerArtist returns true if the music should be saved per artist
func ShouldSaveMusicPerArtist() bool {
	return config.SaveMusicPerArtist
}

// ShouldSaveMusicPerGenre returns true if the music should be saved per genre
func ShouldSaveMusicPerGenre() bool {
	return config.SaveMusicPerGenre
}

// IsValid checks if the music name origin is valid
func (m MusicNameOrigin) IsValid() error {
	switch m {
	case MetadataOrigin, FileNameOrigin:
		return nil
	}
	return fmt.Errorf("invalid music name origin, %s", m)
}

// GetMusicsPathDirectory returns the musics path directory
func GetMusicsPathDirectory() string {
	return config.MusicsPathDirectory
}

// GetFilePathList returns the file path list
func GetFilePathList() []string {
	return config.FilePathList
}

// GetMusicNameOrigin returns the music name origin
func GetMusicNameOrigin() MusicNameOrigin {
	return config.MusicNameOrigin
}

// ShouldPrintArtistPopularity returns true if the artist popularity should be printed
func ShouldPrintArtistPopularity() bool {
	return config.PrintArtistPopularity
}

// ShouldPrintErrorLogs returns true if the error logs should be printed
func ShouldPrintErrorLogs() bool {
	return config.PrintErrorLogs
}
