package twitterService

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (twitter *TwitterService) ReadAccountData(resp *http.Response) (string, string, string, error) {

	// Create variables for scope
	var id, name, followers []byte

	// Reading response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error %v with ReadAccountData(ioutil.ReadAll)", err)
		return "", "", "", err
	}

	// Counter for errors
	errorCounter := 0

	// Parse json array
	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _, _, err = jsonparser.Get(value, "id")
		if err != nil {
			log.Printf("Error %v with ReadAccountData(jsonparser.Get(value, \"id\"))", err)
			errorCounter++
		}
		name, _, _, err = jsonparser.Get(value, "name")
		if err != nil {
			log.Printf("Error %v with ReadAccountData(jsonparser.Get(value, \"name\"))", err)
			errorCounter++
		}
		followers, _, _, err = jsonparser.Get(value, "followers_count")
		if err != nil {
			log.Printf("Error %v with ReadAccountData(jsonparser.Get(value, \"followers_count\"))", err)
			errorCounter++
		}
	})

	// Check error with parsing json
	if errorCounter != 0 {
		return "", "", "", err
	}

	return string(id), string(name), string(followers), nil
}

func (twitter *TwitterService) ReadAccountFollowing(resp *http.Response) ([]AccountData, error, string) {

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	var data, data2, data3, legacy []byte
	var id, name, followers []byte
	var following []AccountData
	var myCursor string // Our cursors
	var cursor []byte   // Byte --> string(myCursor)

	data, _, _, err = jsonparser.Get(body, "data", "user", "result", "timeline", "timeline", "instructions")
	if err != nil {
		log.Printf("Error %v with ReadAccountFollowing(jsonparser.Get(value, \"instructions\"))", err)
		return nil, err, ""
	}

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		data2, _, _, err = jsonparser.Get(value, "entries")
	})

	_, err = jsonparser.ArrayEach(data2, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

		// Findind cursor
		entryId, _, _, err := jsonparser.Get(value, "entryId")

		if strings.Count(string(entryId), "cursor-bottom") == 1 {
			cursor, _, _, err = jsonparser.Get(value, "content", "value")

			myCursor = string(cursor)
			fmt.Println("\n", myCursor)
		}

		// Getting followings
		data3, _, _, err = jsonparser.Get(value, "content", "itemContent", "user_results", "result")
		legacy, _, _, err = jsonparser.Get(value, "content", "itemContent", "user_results", "result", "legacy")

		id, _, _, err = jsonparser.Get(data3, "rest_id")
		name, _, _, err = jsonparser.Get(legacy, "name")
		followers, _, _, err = jsonparser.Get(legacy, "followers_count")

		if len(id) != 0 || len(name) != 0 || len(followers) != 0 {

			account := AccountData{
				Id:        string(id),
				Name:      string(name),
				Followers: string(followers),
			}

			following = append(following, account)
		}

	})

	if err != nil {
		log.Printf("Error %v with ReadAccountFollowing(jsonparser.ArrayEach(value, \"legacy\"))", err)
		return nil, err, ""
	}

	return following, nil, myCursor
}
