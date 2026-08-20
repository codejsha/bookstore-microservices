package repo

import "time"

type WorkResult struct {
	Id               int64
	Uid              string
	Title            string
	Description      *string
	CoverUids        []string
	FirstPublishDate *string
	OlKey            *string
	AuthorUids       []string
	AuthorNames      []string
	SubjectUids      []string
	SubjectNames     []string
	CreatedAt        time.Time
	UpdatedAt        *time.Time
}

type EditionResult struct {
	Id             int64
	Uid            string
	Title          string
	Isbn10         *string
	Isbn13         *string
	NumberOfPages  *int32
	PublishDate    *string
	CoverUids      []string
	Languages      []string
	PhysicalFormat *string
	Description    *string
	WorkUid        string
	WorkTitle      string
	PublisherUid   *string
	PublisherName  *string
	OlKey          *string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type AuthorResult struct {
	Id             int64
	Uid            string
	Name           string
	Bio            *string
	BirthDate      *string
	DeathDate      *string
	PhotoUids      []string
	AlternateNames []string
	OlKey          *string
}

type PublisherResult struct {
	Id      int64
	Uid     string
	Name    string
	Address *string
	OlKey   *string
}

type SubjectResult struct {
	Id   int64
	Uid  string
	Name string
}

type WorkSearchAuthor struct {
	Uid            string
	Name           string
	Bio            *string
	BirthDate      *string
	DeathDate      *string
	OlKey          *string
	AlternateNames []string
}

type WorkSearchSubject struct {
	Uid  string
	Name string
}

type WorkSearchPublisher struct {
	Uid  string
	Name string
}

type WorkSearchEdition struct {
	Uid            string
	Title          string
	Isbn10         *string
	Isbn13         *string
	PublishDate    *string
	Languages      []string
	PhysicalFormat *string
	Description    *string
	Publisher      *WorkSearchPublisher
}

type WorkSearchHighlight struct {
	Title       []string
	Description []string
}

type WorkSearchResult struct {
	Uid              string
	Title            string
	Description      *string
	CoverUids        []string
	FirstPublishDate *string
	OlKey            *string
	Authors          []WorkSearchAuthor
	Subjects         []WorkSearchSubject
	Editions         []WorkSearchEdition
	Score            float32
	Highlight        *WorkSearchHighlight
	CreatedAt        time.Time
	UpdatedAt        *time.Time
}
