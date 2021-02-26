package main

import (
	"encoding/hex"
	"net/url"
	"path"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func getHexUUID(id uuid.UUID) string {
	uuidbin, err := id.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(uuidbin)
}

func getURL(id uuid.UUID) string {
	hexstring := getHexUUID(id)
	filename := hexstring + ".mp3"

	u, err := url.Parse(Config.BaseURL)

	if err != nil {
		panic("Configured Base URL is invalid")
	}

	u.Path = path.Join("/instapod/files/", filename)
	return u.String()
}
