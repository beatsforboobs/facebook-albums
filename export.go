package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

var (
	buckets = []string{"beatsforboobs-staging", "beatsforboobs-production"}
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

func ExportAlbums(accessToken string) ([]Album, []byte, error) {
	albums := struct {
		Albums []Album `json:"data"`
	}{}

	fmt.Println("using facebook graph api to retrieve beatsforboobs albums")
	resp, err := http.Get("https://graph.facebook.com/v2.2/beatsforboobs/albums?access_token=" + url.QueryEscape(accessToken))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(data, &albums)
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("retrieved %d albums\n", len(albums.Albums))
	return albums.Albums, data, nil
}

func ExportPhotos(accessToken, albumId string) ([]byte, error) {
	fmt.Printf("exporting album, %s\n", albumId)

	uri := fmt.Sprintf("https://graph.facebook.com/v2.2/%s/photos?access_token=%s", albumId, url.QueryEscape(accessToken))
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buffer := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buffer, resp.Body)
	return buffer.Bytes(), err
}

func copyToS3(bucketName, path string, data []byte) error {
	auth, err := aws.EnvAuth()
	if err != nil {
		return err
	}

	fmt.Printf("copying to s3://%s%s\n", bucketName, path)

	api := s3.New(auth, aws.USEast)
	bucket := api.Bucket(bucketName)
	err = bucket.Put(path, data, "application/json", "public-read")
	return err
}

func Export(accessToken string) error {
	albums, data, err := ExportAlbums(accessToken)
	if err != nil {
		return err
	}
	for _, bucket := range buckets {
		err = copyToS3(bucket, "/facebook/albums.json", data)
		if err != nil {
			return err
		}
	}

	for _, album := range albums {
		data, err = ExportPhotos(accessToken, album.Id)
		if err != nil {
			return err
		}

		path := fmt.Sprintf("/facebook/album-%s.json", album.Id)
		for _, bucket := range buckets {
			copyToS3(bucket, path, data)
		}
	}

	return nil
}
