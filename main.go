package main

import (
	"encoding/json"
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

// Default Request Handler
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	css := `<style media="screen">
			body { background: #666; }
		
		
			
		ul.hover_block li, ul.hover_block2 li {
			list-style:none;
			float:left;
			background: #fff;
			padding: 10px;
			width:200px; position: relative;
			margin-right: 20px; 
			margin-bottom: 10px;
			}

		ul.hover_block li a, ul.hover_block2 li a {
			display: block;
			position: relative;
			overflow: hidden;
			height: 83px;
			width: 178px;
			padding: 16px;
			color: #000;
			font: 1.6em/1.3 Helvetica, Arial, sans-serif;
		}

		ul.hover_block li a, ul.hover_block2 li a { text-decoration: none; }

		ul.hover_block li img, ul.hover_block2 li img {
			position: absolute;
			top: 0;
			left: 0;
			border: 0;
height: 100%;

vertical-align: middle;
		}
		</style>`
	script := `<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script> 
		<script type="text/javascript">
	$(function() {
			$('ul.hover_block li').hover(function(){
				$(this).find('img').animate({top:'182px'},{queue:false,duration:500});
			}, function(){
				$(this).find('img').animate({top:'0px'},{queue:false,duration:500});
			});
			$('ul.hover_block2 li').hover(function(){
				$(this).find('img').animate({left:'300px'},{queue:false,duration:500});
			}, function(){
				$(this).find('img').animate({left:'0px'},{queue:false,duration:500});
			});
		});
		</script>`
	screen_name := ""
	if r.URL.Path[1:] != "" {
		screen_name = r.URL.Path[1:]
	} else {
		screen_name = "spolsky"
	}
	userInfo, profPic := twitHandler(screen_name)

	html1 := `<ul class="hover_block" style='position:relative'>`
	fmt.Fprintf(w, " %s", css)
	fmt.Fprintf(w, " %s", script)

	//sorting
	sort.Sort(ByFollowersCount{userInfo})
	//displaying 
	imgNormalSplt := strings.Split(profPic, ".")
	imgExt := imgNormalSplt[len(imgNormalSplt)-1]
	imgSrc1 := strings.TrimRight(profPic, "normal."+imgExt)
	imgSrc2 := imgSrc1[0 : len(imgSrc1)-1]
	imgSrc := imgSrc2 + "." + imgExt
	fmt.Fprintf(w, "<table><tr>")
	fmt.Fprintf(w, "<td><img src='"+imgSrc+"' alt='image not found' width='200px' height='200px'/> </td>")
	fmt.Fprintf(w, "<td><h1 style='text-align:right'>Twitter Puzzle</h1>")
	fmt.Fprintf(w, "<h4><i>To solve the Puzzle <br/>you only need to guess these <b>Numbers</b></i></h4></td>")

	fmt.Fprintf(w, "</tr></table>")
	fmt.Fprintf(w, "%s", html1)

	//fmt.Fprintf(w, "<table><tr><th>Screen Name</th><th>Image</th><th>Followers count</th></tr>")
	for i, ele := range userInfo {
		if i >= 10 {
			break
		}
		//strBArr:=strings.Split(ele.Profile_image_url,"/")
		//endStr:=strBArr[len(strBArr)-1]
		imgNormalSplt := strings.Split(ele.Profile_image_url, ".")
		imgExt := imgNormalSplt[len(imgNormalSplt)-1]
		imgSrc1 := strings.TrimRight(ele.Profile_image_url, "normal."+imgExt)
		imgSrc2 := imgSrc1[0 : len(imgSrc1)-1]
		//copy(imgSrc2,imgSrc1[0:len(imgSrc1)-2])
		imgSrc := imgSrc2 + "." + imgExt

		fmt.Println(imgSrc, "imgSrc")
		fmt.Println(ele.Profile_image_url, "Profile_image_url")

		fmt.Fprintf(w, "<li><a href='#'><img src='"+imgSrc+"' alt='img not found'>%d</a></li>", ele.Followers_count)

		if (i+1)%2 == 0 {
			fmt.Fprintf(w, "</ul>")
			fmt.Fprintf(w, "<ul class='hover_block2'>")
		} //fmt.Fprintf(w, "<tr><td>"+ ele.Screen_name+ "</td><td><img src='"+ele.Profile_image_url+"'/></td><td>%d</td></tr>",ele.Followers_count)
	}
	fmt.Fprintf(w, "</ul>")

	//fmt.Fprintf(w, "</table>")
}

//type to use with sorter interface
type UsersInfo []UserInfo

//idiomatic go implementing sorter
func (s UsersInfo) Len() int      { return len(s) }
func (s UsersInfo) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

//to sort  by followers count
type ByFollowersCount struct{ UsersInfo }

//making the highest first
func (s ByFollowersCount) Less(i, j int) bool {
	return s.UsersInfo[i].Followers_count > s.UsersInfo[j].Followers_count
}

func main() {
	//start handling the requests at root with a handler function
	http.HandleFunc("/", defaultHandler)
	//starting the web server at port 8080
	//http.ListenAndServe(":8080", nil)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}

//this includes all twitter interactions and returns the required list of data in []UserInfo
func twitHandler(screen_name string) ([]UserInfo, string) {
	var (
		client *twittergo.Client
	)
	//initializing twitter client configuration
	config := &oauth1a.ClientConfig{
		ConsumerKey:    "XpxtW8CTzEu2zWB1jC1uQ",
		ConsumerSecret: "eowFB6GdzlEca88SATTqpTUcvmd3UiIig4Mi6yaHej8",
	}

	client = twittergo.NewClient(config, nil)
	//set the screen name of the user 

	//ids := getFollowersIds(screen_name, client)
	//fmt.Println(ids)
	//getInfo(ids, client)

	//get the tweets of the user by using screen name
	tweets := getRetweets(screen_name, client)
	//get user ids who retweeted those tweets
	ids := getRetweeters(tweets, client)
	//get user inforamtion like followers count etc usinng the ids
	userInfo := getInfo(ids, client)
	return userInfo, getUserProfPic(screen_name, client)

}

/**
--getRetweeters(tline twittergo.Timeline, client *twittergo.Client) []uint
--Returns user ids by querying retweeters api 1.1
*/
func getRetweeters(tline twittergo.Timeline, client *twittergo.Client) []uint {
	var (
		req  *http.Request
		err  error
		resp *twittergo.APIResponse
		//flw  Followers
	)

	query := url.Values{}

	query.Set("count", "50")
	ids := make([]uint, 0)
	//loop for each retweeted tweet and get the user ids
	for i, tweet := range tline {
		query.Set("id", tweet.IdStr())
		url := fmt.Sprintf("/1.1/statuses/retweeters/ids.json?%v", query.Encode())
		fmt.Println(url)
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Could not parse request: %v\n", err)
			os.Exit(1)
		}
		resp, err = client.SendRequest(req)
		if err != nil {
			fmt.Printf("Could not send request: %v\n", err)
			os.Exit(1)
		}
		//fmt.Println(resp)
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))

		flw := new(Followers)
		json.Unmarshal(b, &flw)
		ids = append(ids, flw.Ids...)
		fmt.Println("retweeted ids", i, ids)
	}

	fmt.Printf("----------------retweeters/ids.json-------------------\n Rate limit: %v\n", resp.RateLimit())
	fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
	fmt.Printf("Rate limit reset: %v\n", resp.RateLimitReset())
	return ids

}

/**
--getRetweets(screen_name string, client *twittergo.Client) []twittergo.Tweet
--Returns filtered timeline (a collection of tweets) by using /1.1/statuses/user_timeline.json api 1.1
*/
func getRetweets(screen_name string, client *twittergo.Client) []twittergo.Tweet {
	var (
		req  *http.Request
		err  error
		resp *twittergo.APIResponse
		//flw  Followers
	)
	//set query params
	query := url.Values{}
	query.Set("screen_name", screen_name)
	query.Set("count", "20")
	query.Set("exclude_replies", "true")

	query.Set("trim_user", "true")
	url := fmt.Sprintf("/1.1/statuses/user_timeline.json?%v", query.Encode())
	//	query.Set("q", "twitterapi")
	//url := fmt.Sprintf("/1.1/search/tweets.json?%v", query.Encode())
	fmt.Println(url)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(resp)
	b, _ := ioutil.ReadAll(resp.Body)

	tline := new(twittergo.Timeline)
	json.Unmarshal(b, &tline)

	for i, ele := range *tline {
		fmt.Println(i, ele.Text())
	}
	fmt.Printf("-------------user_timeline.json--------------------------\nRate limit: %v\n", resp.RateLimit())
	fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
	fmt.Printf("Rate limit reset: %v\n", resp.RateLimitReset())
	return *tline
}

type Followers struct {
	Ids []uint
}
type UserInfo struct {
	Screen_name       string
	Profile_image_url string
	Followers_count   int
}

/**
--getFollowersIds(screen_name string, client *twittergo.Client) []uint
--Returns follower ids by using /1.1/followers/ids.json? api 1.1
*/
func getFollowersIds(screen_name string, client *twittergo.Client) []uint {
	var (
		req  *http.Request
		err  error
		resp *twittergo.APIResponse
	)

	query := url.Values{}
	query.Set("screen_name", screen_name)
	query.Set("count", "10")
	query.Set("cursor", "-1")

	url := fmt.Sprintf("/1.1/followers/ids.json?%v", query.Encode())
	//	query.Set("q", "twitterapi")
	//url := fmt.Sprintf("/1.1/search/tweets.json?%v", query.Encode())
	fmt.Println(url)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(resp)
	b, _ := ioutil.ReadAll(resp.Body)

	flw := new(Followers)
	json.Unmarshal(b, &flw)

	fmt.Println(string(b))
	fmt.Printf("Rate limit: %v\n", resp.RateLimit())
	fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
	fmt.Printf("Rate limit reset: %v\n", resp.RateLimitReset())
	return flw.Ids
}

/**
--getInfo(ids []uint, client *twittergo.Client) []UserInfo 
--Returns a slice of UserInfo by using /1.1/users/lookup.json api 1.1
*/
func getInfo(ids []uint, client *twittergo.Client) []UserInfo {
	var (
		req  *http.Request
		err  error
		resp *twittergo.APIResponse
	)
	//set query string
	query := url.Values{}
	idStr := ""
	for i, e := range ids {
		if i>= 100 {
			break;
		}
		idStr = fmt.Sprintf("%d", e) + "," + idStr
	}
	query.Set("user_id", idStr)
	fmt.Println(idStr)

	url := fmt.Sprintf("/1.1/users/lookup.json?%v", query.Encode())

	fmt.Println(url)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	//fmt.Println(resp)
	b, _ := ioutil.ReadAll(resp.Body)

	uInfo := make([]UserInfo, 0)
	json.Unmarshal(b, &uInfo)
	fmt.Println("uInfo", uInfo)
	//fmt.Println(string(b))
	fmt.Printf("---------users/lookup.------------------Rate limit: %v\n", resp.RateLimit())
	fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
	fmt.Printf("Rate limit reset: %v\n", resp.RateLimitReset())
	return uInfo
}
func getUserProfPic(screen_name string, client *twittergo.Client) string {
	var (
		req  *http.Request
		err  error
		resp *twittergo.APIResponse
	)
	//set query string
	query := url.Values{}

	query.Set("screen_name", screen_name)

	url := fmt.Sprintf("/1.1/users/lookup.json?%v", query.Encode())

	fmt.Println(url)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	//fmt.Println(resp)
	b, _ := ioutil.ReadAll(resp.Body)

	uInfo := make([]UserInfo, 0)
	json.Unmarshal(b, &uInfo)
	fmt.Println("uInfo", uInfo)
	//fmt.Println(string(b))
	fmt.Printf("---------users/lookup.------------------Rate limit: %v\n", resp.RateLimit())
	fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
	fmt.Printf("Rate limit reset: %v\n", resp.RateLimitReset())
	return uInfo[0].Profile_image_url
}
