package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type APIChecker struct {
	UserName string
	APIKey   string

	LastTime time.Time
	LastPP   map[string]float32

	APIBase string
}

func NewAPIChecker(username, apikey string) *APIChecker {
	return &APIChecker{
		UserName: username,
		APIKey:   apikey,
		LastTime: time.Now(),
		LastPP:   make(map[string]float32),
		APIBase:  "http://osu.ppy.sh/api",
	}
}

func (this *APIChecker) constructRecentPlayURL(modenum string) (*url.URL, error) {
	apiurl, err := url.Parse(fmt.Sprintf("%s/get_user_recent", this.APIBase))
	if err != nil {
		return nil, err
	}

	query := apiurl.Query()
	query.Set("k", this.APIKey)
	query.Set("u", this.UserName)
	query.Set("limit", "1")
	query.Set("m", modenum)
	query.Set("type", "string")

	apiurl.RawQuery = query.Encode()
	return apiurl, nil
}

func (this *APIChecker) constructBeatmapGetURL(beatmapid int) (*url.URL, error) {
	apiurl, err := url.Parse(fmt.Sprintf("%s/get_beatmaps", this.APIBase))
	if err != nil {
		return nil, err
	}

	query := apiurl.Query()
	query.Set("k", this.APIKey)
	query.Set("b", fmt.Sprintf("%d", beatmapid))

	apiurl.RawQuery = query.Encode()
	return apiurl, nil
}

func (this *APIChecker) constructUserGetURL(modenum string) (*url.URL, error) {
	apiurl, err := url.Parse(fmt.Sprintf("%s/get_user", this.APIBase))
	if err != nil {
		return nil, err
	}

	query := apiurl.Query()
	query.Set("k", this.APIKey)
	query.Set("u", this.UserName)
	query.Set("m", modenum)
	query.Set("type", "string")

	apiurl.RawQuery = query.Encode()
	return apiurl, nil
}

func (this *APIChecker) CheckForPlay(gamemodes []string) (bool, *PlayInfo, error) {
	var recentPlayData []map[string]interface{}
	var beatmapInfo []map[string]interface{}
	var userInfo []map[string]interface{}

	var gamemode string
	var modenum string
	var beatmapID int
	var date time.Time
	var err error
	var body []byte
	var resp *http.Response

	// Get the most recent play
	for _, mode := range gamemodes {
		gamemode = mode // to persist outside of the loop
		switch mode {
		case "osu":
			modenum = "0"
		case "taiko":
			modenum = "1"
		case "ctb":
			modenum = "2"
		case "mania":
			modenum = "3"
		default:
			return false, &PlayInfo{}, fmt.Errorf("invalid gamemode given: %s", mode)
		}

		recenturl, err := this.constructRecentPlayURL(modenum)
		if err != nil {
			return false, &PlayInfo{}, fmt.Errorf("failed to construct recent plays url: %s", err)
		}

		resp, err := http.Get(recenturl.String())
		if err != nil {
			return false, &PlayInfo{}, fmt.Errorf("failed to get recent plays: %s", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, &PlayInfo{}, fmt.Errorf("failed to get recent plays: %s", err)
		}

		resp.Body.Close()

		err = json.Unmarshal(body, &recentPlayData)
		if err != nil {
			log.Debugf("failed to unmarshal recent plays, erroring. content: %s", string(body))
			return false, &PlayInfo{}, fmt.Errorf("failed to get recent plays: %s", err)
		}

		if len(recentPlayData) < 1 {
			// no plays for the current gamemode
			continue
		}

		// Check that we've actually got a new score
		date, err = time.Parse("2006-01-02 15:04:05", recentPlayData[0]["date"].(string))
		if err != nil {
			return false, &PlayInfo{}, fmt.Errorf("failed to parse date: %s", err)
		}

		if date.Unix() <= this.LastTime.Unix() {
			log.Debugf("we've seen this play before (%s)", date.Format("2006-01-02 15:04:05"))
			continue
		}

		beatmapID, err = strconv.Atoi(recentPlayData[0]["beatmap_id"].(string))
		if err != nil {
			e := fmt.Errorf("can't convert beatmap_id %s to int: %s", recentPlayData[0]["beatmap_id"].(string), err)
			return false, &PlayInfo{}, e
		}

		log.Debugf("new play: gamemode %s, beatmap %d", mode, beatmapID)
		break
	}

	if beatmapID == 0 {
		// no new plays, just return
		log.Debugf("no new plays, returning")
		return false, &PlayInfo{}, nil
	}

	// Get the beatmap information

	beatmap, err := this.constructBeatmapGetURL(beatmapID)
	if err != nil {
		log.Warnf("failed to construct beatmap url: %s", err)
		return false, &PlayInfo{}, nil
	}

	resp, err = http.Get(beatmap.String())
	if err != nil {
		log.Warnf("failed to get beatmap: %s", err)
		return false, &PlayInfo{}, nil
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("failed to get beatmap: %s", err)
		return false, &PlayInfo{}, nil
	}

	resp.Body.Close()

	err = json.Unmarshal(body, &beatmapInfo)
	if err != nil {
		log.Debugf("failed to unmarshal beatmap, erroring. content: %s", string(body))
		log.Warnf("failed to get beatmap: %s", err)
		return false, &PlayInfo{}, nil
	}

	// Get the user information to check PP

	userurl, err := this.constructUserGetURL(modenum)
	if err != nil {
		log.Warnf("failed to construct user info url: %s", err)
		return false, &PlayInfo{}, nil
	}

	resp, err = http.Get(userurl.String())
	if err != nil {
		log.Warnf("failed to get user info: %s", err)
		return false, &PlayInfo{}, nil
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("failed to get user info: %s", err)
		return false, &PlayInfo{}, nil
	}

	resp.Body.Close()

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		log.Debugf("failed to unmarshal user info, erroring. content: %s", string(body))
		log.Warnf("failed to get user info: %s", err)
		return false, &PlayInfo{}, nil
	}

	// Construct the PlayInfo object

	score, err := strconv.Atoi(recentPlayData[0]["score"].(string))
	if err != nil {
		log.Errorf("can't convert score %s to int: %s", recentPlayData[0]["score"].(string), err)
		return false, &PlayInfo{}, nil
	}

	combo, err := strconv.Atoi(recentPlayData[0]["maxcombo"].(string))
	if err != nil {
		log.Errorf("can't convert maxcombo %s to int: %s", recentPlayData[0]["maxcombo"].(string), err)
		return false, &PlayInfo{}, nil
	}

	pp, err := strconv.ParseFloat(userInfo[0]["pp_raw"].(string), 32)
	if err != nil {
		log.Errorf("can't convert pp %s to float: %s", userInfo[0]["pp_raw"].(string), err)
		return false, &PlayInfo{}, nil
	}

	playinfo := &PlayInfo{
		Time:       date,
		GameMode:   gamemode,
		BeatmapID:  beatmapID,
		Beatmap:    fmt.Sprintf("%s - %s", beatmapInfo[0]["artist"].(string), beatmapInfo[0]["title"].(string)),
		Difficulty: beatmapInfo[0]["version"].(string),
		Rank:       recentPlayData[0]["rank"].(string),
		Score:      score,
		MaxCombo:   combo,
		Perfect:    recentPlayData[0]["perfect"].(string) == "1",
		TotalPP:    this.LastPP[modenum],
		GainedPP:   float32(pp) - this.LastPP[modenum],
	}

	this.LastTime = date
	this.LastPP[modenum] = float32(pp)

	return true, playinfo, nil
}
