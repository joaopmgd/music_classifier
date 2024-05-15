package artist

// Popularity is the struct that represents the artist popularity
type Popularity struct {
	Artist string
	Count  int
}

// PopularityMap is a map of Popularity
type PopularityMap map[string]*Popularity

// PopularityList is a list of Popularity
type PopularityList []Popularity

func (a PopularityList) Len() int { return len(a) }
func (a PopularityList) Less(i, j int) bool {
	if a[i].Count != a[j].Count {
		return a[i].Count > a[j].Count
	}
	return a[i].Artist < a[j].Artist
}
func (a PopularityList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// GetArtistNameList returns a list of artist names
func (a PopularityList) GetArtistNameList() []string {
	artistNameList := []string{}
	for _, artist := range a {
		artistNameList = append(artistNameList, artist.Artist)
	}
	return artistNameList
}

// ToList returns a list of Popularity
func (a PopularityMap) ToList() PopularityList {
	PopularityList := PopularityList{}
	for _, Popularity := range a {
		PopularityList = append(PopularityList, *Popularity)
	}
	return PopularityList
}
