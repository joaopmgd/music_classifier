# Music Classifier

## Description

This project extracts music data from file names or metadata and organizes the files into a directory tree based on Artists and Genres. It allows you to either copy or create links to the new files in the specified directory structure.

## Features

- Extracts music data from file names or metadata.
- Organizes music files into a directory structure by Artists and Genres.
- Option to copy or create symbolic links to the new files.

## Requirements

- Go 1.16 or higher

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/joaopmgd/music_classifier.git
   cd music_classifier
   ```

2. Build the project:

   ```bash
   go build
   ```

## Configuration

### General Configuration

The configuration for the execution is stored in the `config.json` file. Edit this file to customize the behavior of the script.

#### `config.json` Structure

```json
{
    "music_path_directory": "/Users/{user}/Music/Musics/",
    "music_name_origin": "METADATA",
    "save_music_per_artist": true,
    "save_music_per_genre": true,
    "print_artist_popularity": true,
    "print_error_logs": true
}
```

- `music_path_directory`: Path to the directory containing the original music files.
- `music_name_origin`: Specifies the origin of the music data, either `METADATA` or `FILE_NAME`.
- `save_music_per_artist`: If true, saves the music in an Artists tree directory.
- `save_music_per_genre`: If true, saves the music in a Genres tree directory.
- `print_artist_popularity`: If true, logs your personal artist preferences after the script is finished.
- `print_error_logs`: If true, logs the files that had errors.

### Genre Classification

The genre classification of the artist is specified in the `genres.json` file. Edit this file to define or update the genre associations for artists.

#### `genres.json` Structure

```json
[
    {
        "artist": "Artist Name 1",
        "genres": [
            "Genre 1",
            "Genre 2"
        ]
    },
    {
        "artist": "Artist Name 2",
        "genres": [
            "Genre 1",
            "Genre 3"
        ]
    }
]
```

## Usage

1. Ensure that the `config.json` and `genres.json` files are correctly configured.
2. Run the executable:

   ```bash
   ./music_classifier
   ```

## Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

## Acknowledgements

- [go-taglib from wtolson](https://github.com/wtolson/go-taglib)
