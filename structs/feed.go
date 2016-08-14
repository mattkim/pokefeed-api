package structs

import "time"

type GetFeedResultStruct struct {
	UUID              string           `json:"uuid"`
	Username          string           `json:"username"`
	CreatedByUserUUID string           `json:"created_by_user_uuid"`
	Message           string           `json:"message"`
	FeedTags          []*FeedTagStruct `json:"feed_tags"`
	Lat               float64          `json:"lat"`
	Long              float64          `json:"long"`
	FormattedAddress  string           `json:"formatted_address"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         time.Time        `json:"deleted_at"`
}

type FeedItemStruct struct {
	UUID              string           `json:"uuid"`
	Username          string           `json:"username"`
	CreatedByUserUUID string           `json:"created_by_user_uuid"`
	Message           string           `json:"message"`
	FeedTags          []*FeedTagStruct `json:"feed_tags"`
	Comments          []*CommentStruct `json:"comments"`
	Lat               float64          `json:"lat"`
	Long              float64          `json:"long"`
	FormattedAddress  string           `json:"formatted_address"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         time.Time        `json:"deleted_at"`
}

type GetFeedsStruct struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type FeedTagStruct struct {
	// TODO: consider renaming
	UUID        string    `json:"uuid"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type PostFeedStruct struct {
	UUID              string              `json:"uuid"`
	CreatedByUserUUID string              `json:"created_by_user_uuid"`
	Message           string              `json:"message"`
	FeedTags          []PostFeedTagStruct `json:"feed_tags"`
	Lat               float64             `json:"lat"`
	Long              float64             `json:"long"`
	FormattedAddress  string              `json:"formatted_address"`
}

type PostFeedTagStruct struct {
	UUID        string `json:"uuid"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	ImageURL    string `json:"image_url"`
}

type ResultStruct struct {
	Result string `json:"result"`
}

type CommentStruct struct {
	UUID              string    `json:"uuid"`
	FeedItemUUID      string    `json:"feed_item_uuid"`
	CreatedByUserUUID string    `json:"created_by_user_uuid"`
	Username          string    `json:"username"`
	Message           string    `json:"message"`
	Lat               float64   `json:"lat"`
	Long              float64   `json:"long"`
	FormattedAddress  string    `json:"formatted_address"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         time.Time `json:"deleted_at"`
}
