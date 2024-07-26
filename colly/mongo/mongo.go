package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/url"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
)

// Storage implements a MongoDB storage backend for colly
type Storage struct {
	Database string
	URI      string
	client   *mongo.Client
	db       *mongo.Database
	visited  *mongo.Collection
	cookies  *mongo.Collection
}

// Init initializes the MongoDB storage
func (s *Storage) Init() error {

	var err error

	clientOptions := options.Client().ApplyURI(s.URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	s.client = client

	s.db = s.client.Database(s.Database)

	s.visited = s.db.Collection("colly_visited")

	s.cookies = s.db.Collection("colly_cookies")

	return nil

}

// VisitedDocument represents the document we want to insert into MongoDB when we have visited a page.
type VisitedDocument struct {
	RequestID string `json:"requestID"`
	Visited   bool   `json:"visited"`
}

// Visited implements colly/storage.VisitedDocument()
func (s *Storage) Visited(requestID uint64) error {

	visited := VisitedDocument{
		RequestID: strconv.FormatUint(requestID, 10),
		Visited:   true,
	}
	_, err := s.visited.InsertOne(context.Background(), visited)

	return err

}

// IsVisited implements colly/storage.IsVisited()
func (s *Storage) IsVisited(requestID uint64) (bool, error) {
	result := &VisitedDocument{}
	visitedQuery := bson.D{{"requestID", strconv.FormatUint(requestID, 10)}}

	err := s.visited.FindOne(nil, visitedQuery).Decode(&result)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {

			return false, nil

		}

		log.Println(err)

		return false, err

	}

	return true, nil

}

type CookiesDocument struct {
	Host    string `json:"host"`
	Cookies string `json:"cookies"`
}

// Cookies implements colly/storage.Cookies()
func (s *Storage) Cookies(u *url.URL) string {

	cookiesQuery := bson.D{{"host", u.Host}}
	result := &CookiesDocument{}

	if err := s.cookies.FindOne(nil, cookiesQuery).Decode(&result); err != nil {

		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
		}

		return ""
	}

	return result.Cookies
}

// SetCookies implements colly/storage.SetCookies()
func (s *Storage) SetCookies(u *url.URL, cookies string) {
	d := CookiesDocument{
		Host:    u.Host,
		Cookies: cookies,
	}

	if _, err := s.cookies.InsertOne(nil, d); err != nil {
		log.Println(err)
	}

}
