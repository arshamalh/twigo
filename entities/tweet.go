package entities

import "time"

type Tweet struct {
	ID                 string                         `json:"id"`
	Text               string                         `json:"text"`
	CreatedAt          time.Time                      `json:"created_at,omitempty"`
	AuthorID           string                         `json:"author_id,omitempty"`
	ConversationID     string                         `json:"conversation_id,omitempty"`
	InReplyToUserID    string                         `json:"in_reply_to_user_id,omitempty"`
	Attachments        map[string][]string            `json:"attachments,omitempty"`
	ContextAnnotations []map[string]ContextAnnotation `json:"context_annotations,omitempty"`
	Entities           TweetEntities                  `json:"entities,omitempty"`
	Geo                Geo                            `json:"geo,omitempty"`
	Lang               string                         `json:"lang,omitempty"`
	NonPublicMetrics   NonPublicMetrics               `json:"non_public_metrics,omitempty"`
	OrganicMetrics     OrganicMetrics                 `json:"organic_metrics,omitempty"`
	PossiblySensitive  bool                           `json:"possibly_sensitive,omitempty"`
	PromotedMetrics    PromotedMetrics                `json:"promoted_metrics,omitempty"`
	PublicMetrics      TweetPublicMetrics             `json:"public_metrics,omitempty"`
	ReferencedTweets   []ReferencedTweet              `json:"referenced_tweets,omitempty"`
	ReplySettings      string                         `json:"reply_settings,omitempty"`
	Source             string                         `json:"source,omitempty"`
	Withheld           TweetWithheld                  `json:"withheld,omitempty"`
}

type ContextAnnotation struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TweetEntities struct {
	Annotations []Annotation     `json:"annotations"`
	CashTags    []TweetEntityTag `json:"cashtags"`
	HashTags    []TweetEntityTag `json:"hashtags"`
	Mentions    []TweetEntityTag `json:"mentions"`
	URLs        []URL            `json:"urls"`
}

type Annotation struct {
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Probability    float64 `json:"probability"`
	Type           string  `json:"type"`
	NormalizedText string  `json:"normalized_text"`
}

type TweetEntityTag struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Tag   string `json:"tag"`
}

type URL struct {
	Start       int    `json:"start"`
	End         int    `json:"end"`
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
	DisplayURL  string `json:"display_url"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UnwoundURL  string `json:"unwound_url"`
}

type Geo struct {
	Coordinates struct {
		Type        string `json:"type"`
		Coordinates []int  `json:"coordinates"`
	} `json:"coordinates"`
	PlaceID string `json:"place_id"`
}

type NonPublicMetrics struct {
	ImpressionCount   int `json:"impression_count"`
	UrlLinkClicks     int `json:"url_link_clicks"`
	UserProfileClicks int `json:"user_profile_clicks"`
}

type OrganicMetrics struct {
	ImpressionCount   int `json:"impression_count"`
	LikeCount         int `json:"like_count"`
	ReplyCount        int `json:"reply_count"`
	RetweetCount      int `json:"retweet_count"`
	UrlLinkClicks     int `json:"url_link_clicks"`
	UserProfileClicks int `json:"user_profile_clicks"`
}

type PromotedMetrics struct {
	ImpressionCount   int `json:"impression_count"`
	LikeCount         int `json:"like_count"`
	ReplyCount        int `json:"reply_count"`
	RetweetCount      int `json:"retweet_count"`
	UrlLinkClicks     int `json:"url_link_clicks"`
	UserProfileClicks int `json:"user_profile_clicks"`
}

type TweetPublicMetrics struct {
	RetweetCount int `json:"retweet_count"`
	ReplyCount   int `json:"reply_count"`
	LikeCount    int `json:"like_count"`
	QuoteCount   int `json:"quote_count"`
}

type ReferencedTweet struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type TweetWithheld struct {
	Copyright    bool     `json:"copyright"`
	CountryCodes []string `json:"country_codes"`
}
