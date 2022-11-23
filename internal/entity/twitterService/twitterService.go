package twitterService

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type TwitterService struct {
	client http.Client
}

func NewTwitterService() (*TwitterService, error) {

	client := http.Client{Timeout: 5 * time.Second}

	twitter := &TwitterService{
		client: client,
	}

	return twitter, nil
}

func (twitter *TwitterService) SetProxy(proxy string) error {

	proxyURL, err := url.Parse(proxy)

	if err != nil {
		log.Println(err)
		return err
	}

	// Creating Transport structure for client
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// Add proxy to client
	twitter.client.Transport = transport

	return nil
}

func (twitter *TwitterService) GetAccountData(username string, proxy string) (string, string, string, error) {

	// Setting proxy for client
	err := twitter.SetProxy(proxy)
	if err != nil {
		log.Printf("Error %v with GetAccountData(twitter.SetProxy)", err)
		return "", "", "", nil
	}

	twitterUrl := fmt.Sprintf("https://cdn.syndication.twimg.com/widgets/followbutton/info.json?screen_names=%v", username)

	// Creating request
	req, err := http.NewRequest(http.MethodGet, twitterUrl, nil)

	// Closing request
	req.Close = true
	if err != nil {
		log.Printf("Error %v with GetAccountData(Creating request)", err)
		return "", "", "", err
	}

	// Sending request
	resp, err := twitter.client.Do(req)
	if err != nil {
		log.Printf("Error %v with GetAccountData(Sending request)", err)
		return "", "", "", err
	}

	// Closing connecting
	defer resp.Body.Close()

	// Reading
	id, name, followers, err := twitter.ReadAccountData(resp)
	if err != nil {
		log.Printf("Error %v with GetAccountData(ReadAccountData)", err)
		return "", "", "", err
	}

	return id, name, followers, nil
}

func (twitter *TwitterService) GetAccountFollowings(userId string, proxy string, cookie string, bearer string, csrf string) ([]AccountData, error, string) {

	// Setting proxy for client
	err := twitter.SetProxy(proxy)
	if err != nil {
		log.Printf("Error %v with GetAccountFollowings(twitter.SetProxy)", err)
		return nil, err, ""
	}

	// Creating url without query params
	urlA, err := url.Parse("https://twitter.com/i/api/graphql/ft89HkYoFtD8H3czDEuGhg/Following?variables=")
	if err != nil {
		log.Printf("Error %v with GetAccountFollowings(url.Parse)", err)
	}

	// Creating query params string
	variables := fmt.Sprintf("{\"userId\":%v,\"count\":115,\"includePromotedContent\":false,\"withSuperFollowsUserFields\":true,\"withDownvotePerspective\":false,\"withReactionsMetadata\":false,\"withReactionsPerspective\":false,\"withSuperFollowsTweetFields\":true}", userId)
	features := "{\"dont_mention_me_view_api_enabled\":true,\"interactive_text_enabled\":true,\"responsive_web_uc_gql_enabled\":false,\"vibe_api_enabled\":false,\"responsive_web_edit_tweet_api_enabled\":false,\"standardized_nudges_misinfo\":true,\"responsive_web_enhance_cards_enabled\":false}"

	// Setting query params
	values := urlA.Query()

	values.Add("variables", "features")

	values.Set("variables", variables)
	values.Set("features", features)

	// Creating url with query params
	urlA.RawQuery = values.Encode()

	twitterUrl := urlA.String()
	fmt.Println(twitterUrl)
	// Creating request
	req, err := http.NewRequest(http.MethodGet, twitterUrl, nil)
	if err != nil {
		log.Printf("Error %v with GetAccountFollowings(http.NewRequest)", err)
		return nil, err, ""
	}

	req.Close = true

	// Adding headers
	req.Header.Add("cookie", cookie)
	req.Header.Add("authorization", bearer)
	req.Header.Add("x-csrf-token", csrf)

	// Sending request
	resp, err := twitter.client.Do(req)

	if err != nil {
		log.Printf("Error %v with GetAccountFollowings(twitter.client.Do)", err)
		return nil, err, ""
	}

	// Closing connection
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	// Reading
	accounts, err, cursor := twitter.ReadAccountFollowing(resp)
	if err != nil {
		log.Printf("Error %v with GetAccountFollowings(twitter.ReadAccountFollowing)", err)
		return nil, err, ""
	}

	return accounts, nil, cursor
}

func (twitter *TwitterService) GetAccountFollowingsContinuation(userId string, proxy string, cookie string, bearer string, csrf string, cursor string) ([]AccountData, error, string) {

	// Setting proxy for client
	err := twitter.SetProxy(proxy)
	if err != nil {
		log.Printf("Error %v with GetAccountFollowingsContinuation(twitter.SetProxy)", err)
		return nil, err, ""
	}

	// Creating url without query params
	urlA, err := url.Parse("https://twitter.com/i/api/graphql/ft89HkYoFtD8H3czDEuGhg/Following?variables=")
	if err != nil {
		log.Printf("Error %v with GetAccountFollowingsContinuation(url.Parse)", err)
		return nil, err, ""
	}

	// Creating query params string
	variables := fmt.Sprintf("{\"userId\":%v,\"count\":115,\"cursor\":\"%v\",\"includePromotedContent\":false,\"withSuperFollowsUserFields\":true,\"withDownvotePerspective\":false,\"withReactionsMetadata\":false,\"withReactionsPerspective\":false,\"withSuperFollowsTweetFields\":true}", userId, cursor)
	features := "{\"unified_cards_follow_card_query_enabled\":false,\"dont_mention_me_view_api_enabled\":true,\"responsive_web_uc_gql_enabled\":true,\"vibe_api_enabled\":true,\"responsive_web_edit_tweet_api_enabled\":true,\"standardized_nudges_misinfo\":true,\"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled\":false,\"interactive_text_enabled\":true,\"responsive_web_text_conversations_enabled\":false,\"responsive_web_enhance_cards_enabled\":true}"

	// Setting query params
	values := urlA.Query()

	values.Add("variables", "features")

	values.Set("variables", variables)
	values.Set("features", features)

	// Creating url with query params
	urlA.RawQuery = values.Encode()

	twitterUrl := urlA.String()

	fmt.Println(twitterUrl)

	// Creating request
	req, err := http.NewRequest(http.MethodGet, twitterUrl, nil)
	if err != nil {
		log.Printf("Error %v with GetAccountFollowingsContinuation(http.NewRequest)", err)
		return nil, err, ""
	}

	req.Close = true

	// Adding headers
	req.Header.Add("cookie", cookie)
	req.Header.Add("authorization", bearer)
	req.Header.Add("x-csrf-token", csrf)

	// Sending request
	resp, err := twitter.client.Do(req)
	if err != nil {
		log.Printf("Error %v with GetAccountFollowingsContinuation(twitter.client.Do)", err)
		return nil, err, ""
	}
	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()

	// Reading
	accounts, err, cursor := twitter.ReadAccountFollowing(resp)
	if err != nil {
		log.Printf("Error %v with GetAccountFollowings(twitter.ReadAccountFollowing)", err)
		return nil, err, ""
	}

	return accounts, nil, cursor
}

func (twitter *TwitterService) GetAllAccountFollowings(userId string, proxy string, cookie string, bearer string, csrf string) ([]AccountData, error) {

	var allSubscriptions []AccountData

	// First request
	accounts, err, cursor := twitter.GetAccountFollowings(userId, proxy, cookie, bearer, csrf)
	if err != nil {
		log.Printf("Error %v with GetAllAccountFollowings(twitter.GetAccountFollowings)", err)
		return nil, err
	}

	for i := 0; i < len(accounts)-1; i++ {
		allSubscriptions = append(allSubscriptions, accounts[i])
	}

	log.Printf("First sended")
	x := 2

	// Other continuation requests
	for {
		time.Sleep(1*time.Second)
		accounts, err, cursor = twitter.GetAccountFollowingsContinuation(userId, proxy, cookie, bearer, csrf, cursor)
		if err != nil {
			log.Printf("Error %v with GetAllAccountFollowings(twitter.GetAccountFollowingsContinuation)", err)
			return nil, err
		}

		log.Printf("%v sended", x)
		x++

		if len(accounts) == 0 {
			break
		} else {
			for i := 0; i < len(accounts)-1; i++ {
				allSubscriptions = append(allSubscriptions, accounts[i])
			}
		}
	}

	return allSubscriptions, nil
}
