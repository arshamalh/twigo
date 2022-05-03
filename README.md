# twigo (Golang Twitter Library)

<p align="center">
  <img src="./twigo.png" alt="twigo logo" width="300">
</p>

### Twigo is a fast and easy to use twitter API library help you write best twitter bots.

Currently we only support twitter api v2, but version 1.1 will be added soon.

Caution! we are still in Beta phase.

## installation

```bash
go get github.com/arshamalh/twigo
```

# How to use
Easily make a new client!
```go
package main

import "github.com/arshamalh/twigo"

twigo.NewClient(
    "ConsumerKey",
    "ConsumerSecret",
    "AccessToken",
    "AccessTokenSecret",
    // You can use bearer token or four other keys (ConsumerKey, ...), both is not mandatory.
    // We'll find bearer token automatically if it's not specified.
    // Also, you can use twigo.BearerFinder() function.
    "BearerToken",
)
```

And use any function you need, for example, get a tweet like this:

```go
response, _ := client.GetTweet(tweet_id, false, nil)
fmt.Printf("%#v\n", response)
tweet := response.Data
```

Passing some extra fields and params:

```go
fields := map[string]interface{}{"tweet.fields": []string{"author_id", "created_at", "public_metrics"}}
response, _ := client.GetTweet(tweet_id, false, fields)
fmt.Printf("%#v\n", response)
tweet := response.Data
```

If the tweet doesn't exist or it's from a suspended account, 
we will return an empty struct instead, like this:

```go
&twigo.TweetResponse{
  Data: twigo.Tweet{
    ID:"", 
    Text:"", 
    PublicMetrics:map[string]int(nil)
    // Other fields
  }, 
  Meta: twigo.MetaEntity{
    ResultCount:0, 
    NewestID:"",
    NextToken:""
    // Other fields
  }, 
  RateLimits: twigo.RateLimits{
    Limit:300, 
    Remaining:294, 
    ResetTimestamp:1651585436
  }
}
```

**More examples:**

```go
response, err := client.GetLikingUsers(
  "1431751228145426438", 
  false, // should we use oauth_1? Can be true, depend on your preferences, but maybe we will change it if needed.
  map[string]interface{}{
    "max_results": 5,
  })

if err != nil {
  fmt.Println(err)
}

fmt.Printf("%+v\n", response)
```

And result will be a Go struct like this:

```Go
{
  Data:[
    {
      ID:1506623384850948096 
      Name:Rhea Baker 
      Username:RheaBak16183941
      Verified:false
    } {
      ID:1506621280098869252 
      Name:Janice Lane 
      Username:JaniceL44359093
      Verified:false
    }
    // And more...
  ] 
  Includes: // Some structs, read more on docs...
  Errors:[] 
  Meta:{
    ResultCount:5 
    NewestID: 
    OldestID: 
    PreviousToken: 
    NextToken:7140dibdnow9c7btw480y5xgmlpwtbsh4fyqnqmwz9k4w
  }
  RateLimits:{
    Limit:75
    Remaining:74
    ResetTimestamp:1650553033
  }
}
```

More examples:

```go
response, err := client.GetUsersByUsernames(
  []string{"arshamalh", "elonmusk", "someone_else"}, 
  true, // Also we suggest you to use false as default.
  nil, // There is no param in this example.
)
```

Retweeting and Liking a tweet:

```go
client.Retweet("1516784368601153548")
client.Like("1431751228145426438")
```

Or Maybe deleting your like and retweet:

```go
client.UnRetweet("1516784368601153548")
client.Unlike("1516784368601153548")
```

And finally! Tweeting! (creating tweet)

```go
client.CreateTweet("This is a test tweet", nil)
```

**Simple, right?**

### How to paginate over results?
if your method is paginatable, you can paginate using NextPage method attached to the response, like this:

```go
response, _ := client.GetUserTweets("1216345203453452289", false, nil)
for {
  fmt.Printf("%#v\n", response)
  response, err = response.NextPage()  // Here the magic happens! NextPage method attached to response

  if err != nil {  // Break the loop whenever there is no more pages.
    fmt.Println(err)
    break
  }
}
```

A little bit more complex example:

```go
response, err := client.GetUserTweets("1216345203453452289", false, nil)

if err != nil {  // Is there any response at all?
  fmt.Println(err)
} else {
  page_number := 1  // To keep track of page number.
  for {
    fmt.Printf("Page %d: %#v\n", page_number, response)
    response, err = response.NextPage()

    if err != nil {
      fmt.Println(err)
      break
    }

    // If rate limits exceeds, wait until it charges again, 
    // it's not a good approach, because all of your app will sleep, 
    // it's highly recommended to do it in a goroutine instead.
    if response.RateLimits.Remaining == 0 {
      fmt.Println("Rate limit reached")
      sleeping_duration := time.Until(time.Unix(response.RateLimits.ResetTimestamp, 0))
      time.Sleep(sleeping_duration)
		}
    page_number++
  }
}
```

### Rate limits
How many actions can we do?

You can simpy read RateLimits attribute on the Response!
```go
  RateLimits:{
    Limit:75 // An static number depending on the endpoint that you are calling or your authentication method.
    Remaining:74 // An dynamic method that decreases after each call, and will reset every once in a while.
    ResetTimestamp:1650553033 // Reset (charge up) remaining calls in this timestamp.
  }
```

## Contribution
Feel free to open an issue, contribute and contact us!
