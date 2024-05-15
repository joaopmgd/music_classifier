package main

import (
	"music_classifier/configuration"
	"music_classifier/music"
)

func main() {

	// Loads all music and artist data from the metadata
	musicList, err := music.GetMusicData()
	if err != nil {
		panic(err)
	}

	// Loop over the saved musics and rename them with the artist name by popularity
	for _, music := range musicList {

		// Create and save the music per artist
		// It links the file to the artist directory
		if configuration.ShouldSaveMusicPerArtist() {
			err = music.SaveMusicPerArtist()
			if err != nil {
				panic(err)
			}
		}

		// Create and save the music per genre
		// It copies the file to the genre directory
		if configuration.ShouldSaveMusicPerGenre() {
			err = music.SaveMusicPerGenre()
			if err != nil {
				panic(err)
			}
		}
	}

	// Print the user artist popularity
	if configuration.ShouldPrintArtistPopularity() {
		music.PrintArtistPopularity()
	}
}
