package twitterService

import (
	"fmt"
	"log"
)

const proxy = ""

func (twitter *TwitterService) TestGetAccountData() {
	id, name, followers, err := twitter.GetAccountData("LogisticNode", proxy)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Id: %v Name: %v Followers: %v\n", id, name, followers)
}

func (twitter *TwitterService) TestGetAccountFollowing() {

	userId := ""

	cookie := ""
	bearer := ""
	csrf := ""

	accounts, err, _ := twitter.GetAccountFollowings(userId, proxy, cookie, bearer, csrf)
	if err != nil {
		fmt.Printf("Error %v with test(TestGetAccountFollowing)", err)
	}

	for i := 0; i < len(accounts)-1; i++ {
		fmt.Printf("%v) Id: %v, Name: %v Followers: %v\n", i+1, accounts[i].Id, accounts[i].Name, accounts[i].Followers)
	}

}

func (twitter *TwitterService) TestGetAccountFollowingContinuatuion(cursor string) {
	userId := ""

	cookie := ""
	bearer := ""
	csrf := ""

	accounts, err, _ := twitter.GetAccountFollowingsContinuation(userId, proxy, cookie, bearer, csrf, cursor)
	if err != nil {
		log.Printf("Error %v with test(TestGetAccountFollowingsContinuation)", err)
	}

	fmt.Println(accounts)
}

func (twitter *TwitterService) TestGetAllAccountFollowings() {

	userId := ""

	cookie := ""
	bearer := ""
	csrf := ""

	accounts, err := twitter.GetAllAccountFollowings(userId, proxy, cookie, bearer, csrf)
	if err != nil {
		fmt.Printf("Error %v with test(TestGetAllAccountFollowings)", err)
	}

	for i := 0; i < len(accounts)-1; i++ {
		fmt.Printf("%v) Id: %v, Name: %v Followers: %v\n", i+1, accounts[i].Id, accounts[i].Name, accounts[i].Followers)
	}
}
