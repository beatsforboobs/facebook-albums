package main

import (
	"net/http"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

type Album struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	CoverPhoto  string `json:"cover_photo"`
	Privacy     string `json:"everyone"`
	CreatedTime string `json:"created_time"`
}

func Export(accessToken, bucketName string, finished chan struct{}) error {
	albums := struct {
		Albums []Album `json:"data"`
	}{}

	resp, err := http.Get()
}

func ExportAlbums(accessToken, bucketName string, finished chan struct{}) {
	auth := aws.EnvAuth()
	api := s3.New(auth, aws.USEast)
	bucket := api.Bucket(bucketName)
	//bucket.Put(path, data, "application/json", "public-read")
	// https://graph.facebook.com/v2.2/beatsforboobs/albums?access_token=...
	// https://graph.facebook.com/v2.2/10152035682934317/photos?access_token=...
}
