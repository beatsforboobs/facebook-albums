package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

func ExportAlbums(accessToken string) ([]Album, error) {
	albums := struct {
		Albums []Album `json:"data"`
	}{}

	resp, err := http.Get("https://graph.facebook.com/v2.2/beatsforboobs/albums?access_token=" + url.QueryEscape(accessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(albums)
	return albums, err
}

func ExportPhotos(accessToken, albumId string) ([]byte, error) {
	uri := fmt.Sprintf("https://graph.facebook.com/v2.2/%s/photos?access_token%s", albumId, url.QueryEscape(accessToken))
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buffer := bytes.NewBuffer()
	err = io.Copy(buffer, resp.Body)
	return buffer.Bytes(), err
}

func Export(accessToken, bucketName string, finished chan struct{}) {
	albums := ExportAlbums(accessToken)
	for _, albumId := range albums {

	}
	auth := aws.EnvAuth()
	api := s3.New(auth, aws.USEast)
	bucket := api.Bucket(bucketName)
	//bucket.Put(path, data, "application/json", "public-read")
	// https://graph.facebook.com/v2.2/beatsforboobs/albums?access_token=...
	// https://graph.facebook.com/v2.2/10152035682934317/photos?access_token=...
}
