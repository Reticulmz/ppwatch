package main

import (
	"testing"
	"net/url"
	"github.com/stretchr/testify/assert"
)

func TestAPICheckRecentPlayURLConstructor(t *testing.T) {
	apichecker := NewAPIChecker("username", "apikey")

	constructedurl, err := apichecker.constructRecentPlayURL()
	assert.Nil(t, err)

	apiurl, err := url.Parse("http://osu.ppy.sh/api/get_user_recent")
	assert.Nil(t, err)

	query := apiurl.Query()
	query.Set("k", "apikey")
	query.Set("u", "username")
	query.Set("limit", "1")
	query.Set("type", "string")
	apiurl.RawQuery = query.Encode()

	assert.Equal(t, apiurl, constructedurl)	
}

func TestAPICheckUserGetURLConstructor(t *testing.T) {
	apichecker := NewAPIChecker("username", "apikey")

	constructedurl, err := apichecker.constructUserGetURL()
	assert.Nil(t, err)

	apiurl, err := url.Parse("http://osu.ppy.sh/api/get_user")
	assert.Nil(t, err)

	query := apiurl.Query()
	query.Set("k", "apikey")
	query.Set("u", "username")
	query.Set("type", "string")
	apiurl.RawQuery = query.Encode()

	assert.Equal(t, apiurl, constructedurl)	
}

func TestAPICheckBeatmapGetURLConstructor(t *testing.T) {
	apichecker := NewAPIChecker("username", "apikey")

	constructedurl, err := apichecker.constructBeatmapGetURL(1)
	assert.Nil(t, err)

	apiurl, err := url.Parse("http://osu.ppy.sh/api/get_beatmaps")
	assert.Nil(t, err)

	query := apiurl.Query()
	query.Set("k", "apikey")
	query.Set("b", "1")
	apiurl.RawQuery = query.Encode()

	assert.Equal(t, apiurl, constructedurl)	
}