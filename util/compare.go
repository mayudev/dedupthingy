package util

import "strings"

type Comparator struct {
	ByArtist  bool
	ByTitle   bool
	ByAlbum   bool
	Sensitive bool
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func NewComparator(matches []string, sensitive bool) *Comparator {
	return &Comparator{
		Sensitive: sensitive,
		ByArtist:  contains(matches, "artist"),
		ByTitle:   contains(matches, "title"),
		ByAlbum:   contains(matches, "album"),
	}

}

func (c *Comparator) ParseString(s string) string {
	if !c.Sensitive {
		return strings.ToLower(s)
	}
	return s
}

func (c *Comparator) CreateComparator(v Metadata) Metadata {
	res := Metadata{}

	if c.ByArtist {
		res.Artist = c.ParseString(v.Artist)
	}

	if c.ByAlbum {
		res.Album = c.ParseString(v.Album)
	}

	if c.ByTitle {
		res.Title = c.ParseString(v.Title)
	}

	return res
}
