package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type IndexItemLocation struct {
	Index []Location `json:"index"`
}
type IndexItemRelation struct {
	Index []Relations `json:"index"`
}
type IndexItemDates struct {
	Index []Dates `json:"index"`
}
type ArtistData struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Relation     string              `json:"relations"`
	Concerts     map[string][]string `json:"datesLocations"`
	Location     map[string]string   `json:"location"`
}

type Dates struct {
	ID   int      `json:"id"`
	Date []string `json:"dates"`
}
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

var (
	ArtistDataAll     []ArtistData
	LocationData      []Location
	DateData          []Dates
	RelationData      []Relations
	IndexDataLocation []IndexItemLocation
	IndexDataRelation []IndexItemRelation
	IndexDataDates    []IndexItemDates
)

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/about", AboutPage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	ArtistDataAll = getArtistData("https://groupietrackers.herokuapp.com/api/artists")
	LocationData = getLocationData("https://groupietrackers.herokuapp.com/api/locations")
	DateData = getDatesData("https://groupietrackers.herokuapp.com/api/dates")
	RelationData = getRelationData("https://groupietrackers.herokuapp.com/api/relation")

	data := struct {
		Artists   []ArtistData
		Locations []Location
		Dates     []Dates
		Relations []Relations
	}{
		Artists:   ArtistDataAll,
		Locations: LocationData,
		Dates:     DateData,
		Relations: RelationData,
	}

	tpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprintf(w, "Template parse error: %v", err)
		return
	}

	search := strings.ToLower(r.URL.Query().Get("Search"))
	if search != "" {
		searchResults := filterArtists(ArtistDataAll, LocationData, DateData, RelationData, search)
		data.Artists = uniqueArtists(searchResults)
	} else {
		data.Artists = ArtistDataAll
	}

	err = tpl.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, "Error executing template: %v", err)
	}
}

func filterArtists(artists []ArtistData, LocationData []Location, DateData []Dates, RelationData []Relations, search string) []ArtistData {
	searchResults := make([]ArtistData, 0)
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), search) {
			searchResults = append(searchResults, artist)
		}
	}
	for _, locations := range LocationData {
		for _, location := range locations.Locations {
			if strings.Contains(strings.ToLower(location), strings.ToLower(search)) {
				for _, artist := range ArtistDataAll {
					if fmt.Sprint(locations.ID) == fmt.Sprint(artist.ID) {
						searchResults = append(searchResults, artist)
					}
				}
			}
		}
	}
	for _, dates := range DateData {
		for _, date := range dates.Date {
			if strings.Contains(strings.ToLower(date), strings.ToLower(search)) {
				for _, artist := range ArtistDataAll {
					if fmt.Sprint(dates.ID) == fmt.Sprint(artist.ID) {
						searchResults = append(searchResults, artist)
					}
				}
			}
		}
	}
	for _, relation := range RelationData {
		for _, datesLocations := range relation.DatesLocations {
			for _, datesLocation := range datesLocations {
				if strings.Contains(strings.ToLower(datesLocation), strings.ToLower(search)) {
					for _, artist := range ArtistDataAll {
						if fmt.Sprint(relation.ID) == fmt.Sprint(artist.ID) {
							searchResults = append(searchResults, artist)
						}
					}
				}
			}
		}
	}
	return searchResults
}

func uniqueArtists(artists []ArtistData) []ArtistData {
	seen := make(map[int]bool)
	unique := []ArtistData{}
	for _, artist := range artists {
		if _, ok := seen[artist.ID]; !ok {
			seen[artist.ID] = true
			unique = append(unique, artist)
		}
	}
	return unique
}

func getLocationData(url string) []Location {
	data1, e1 := http.Get(url)
	if e1 != nil {
		fmt.Println("HTTP GET hatasi:", e1)
		return []Location{}
	}
	defer data1.Body.Close()
	data, e2 := io.ReadAll(data1.Body)
	if e2 != nil {
		fmt.Println("HTTP yaniti okuma hatasi:", e2)
		return []Location{}
	}
	var locationData IndexItemLocation
	e3 := json.Unmarshal(data, &locationData)
	if e3 != nil {
		fmt.Println("JSON parse hatasi:", e3)
		return []Location{}
	}
	return locationData.Index
}

func getRelationData(url string) []Relations {
	data1, e1 := http.Get(url)
	if e1 != nil {
		fmt.Println("HTTP GET hatasi:", e1)
		return []Relations{}
	}
	defer data1.Body.Close()
	data, e2 := io.ReadAll(data1.Body)
	if e2 != nil {
		fmt.Println("HTTP yaniti okuma hatasi:", e2)
		return []Relations{}
	}
	var relationData IndexItemRelation
	e3 := json.Unmarshal(data, &relationData)
	if e3 != nil {
		fmt.Println("JSON parse hatasi:", e3)
		return []Relations{}
	}
	return relationData.Index
}

func getDatesData(url string) []Dates {
	data1, e1 := http.Get(url)
	if e1 != nil {
		fmt.Println("HTTP GET hatasi:", e1)
		return []Dates{}
	}
	defer data1.Body.Close()
	data, e2 := io.ReadAll(data1.Body)
	if e2 != nil {
		fmt.Println("HTTP yaniti okuma hatasi:", e2)
		return []Dates{}
	}
	var dateData IndexItemDates
	e3 := json.Unmarshal(data, &dateData)
	if e3 != nil {
		fmt.Println("JSON parse hatasi:", e3)
		return []Dates{}
	}
	return dateData.Index
}

func getArtistData(url string) []ArtistData {
	data1, e1 := http.Get(url)

	if e1 != nil {
		panic(e1)
	}
	data, _ := io.ReadAll(data1.Body)
	_ = data1.Body.Close()
	e := json.Unmarshal(data, &ArtistDataAll)
	if e != nil {
		panic(e)
	}
	return ArtistDataAll
}

func AboutPage(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("about.html")
	if err != nil {
		fmt.Fprintf(w, "Error parsing template: %v", err)
		return
	}

	artistID := r.URL.Query().Get("id")

	type CombinedData struct {
		Artist   *ArtistData
		Location *Location
		Date     *Dates
		Relation *Relations
	}

	var combined CombinedData

	for _, artist := range ArtistDataAll {
		if fmt.Sprint(artist.ID) == artistID {
			combined.Artist = &artist
			break
		}
	}
	for _, location := range LocationData {
		if fmt.Sprint(location.ID) == artistID {
			combined.Location = &location
			break
		}
	}

	for _, date := range DateData {
		if fmt.Sprint(date.ID) == artistID {
			combined.Date = &date
			break
		}
	}

	for _, relation := range RelationData {
		if fmt.Sprint(relation.ID) == artistID {
			combined.Relation = &relation
			break
		}
	}

	err = tpl.Execute(w, combined)
	if err != nil {
		fmt.Fprintf(w, "Error executing template: %v", err)
	}
}
