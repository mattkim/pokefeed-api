package mappers

import (
	"github.com/pokefeed/pokefeed-api/models"
	"github.com/pokefeed/pokefeed-api/structs"
)

// FeedTagMapperArrayDBToJSON Map the array db struct to the array json struct
func FeedTagMapperArrayDBToJSON(feedTags []models.FeedTagRow) []*structs.FeedTagStruct {
	feedTagsResult := []*structs.FeedTagStruct{}

	for _, feedTag := range feedTags {
		feedTagResult := FeedTagMapperDBToJSON(feedTag)
		feedTagsResult = append(feedTagsResult, &(feedTagResult))
	}

	return feedTagsResult
}

// FeedTagMapperDBToJSON Map the db struct to the json struct
func FeedTagMapperDBToJSON(feedTag models.FeedTagRow) structs.FeedTagStruct {
	return structs.FeedTagStruct{
		UUID:        feedTag.UUID,
		Type:        feedTag.Type,
		Name:        feedTag.Name,
		DisplayName: feedTag.DisplayName,
		ImageURL:    feedTag.ImageURL,
		CreatedAt:   feedTag.CreatedAt.Time,
		UpdatedAt:   feedTag.UpdatedAt.Time,
		DeletedAt:   feedTag.DeletedAt.Time,
	}
}

// FeedItemMapperDBToJSON Map the db struct to the json struct
func FeedItemMapperDBToJSON(
	feedItem models.FeedItemRow,
	user *models.UserRow, // TODO: convert to struct before passing in here.
	facebookInfo *models.FacebookInfoRow, // TODO: convert to struct before passing in here.
	feedTags []*structs.FeedTagStruct,
	comments []*structs.CommentStruct,
) structs.FeedItemStruct {
	var username string

	if len(user.Username) > 0 {
		username = user.Username
	} else {
		// TODO: this could be empty too
		username = facebookInfo.FacebookName
	}

	return structs.FeedItemStruct{
		UUID:              feedItem.UUID,
		Message:           feedItem.Message,
		Lat:               feedItem.Lat,
		Long:              feedItem.Long,
		FormattedAddress:  feedItem.FormattedAddress,
		Username:          username,
		CreatedByUserUUID: user.UUID,
		FeedTags:          feedTags,
		Comments:          comments,
		CreatedAt:         feedItem.CreatedAt.Time,
		UpdatedAt:         feedItem.UpdatedAt.Time,
		DeletedAt:         feedItem.DeletedAt.Time,
	}
}

// CommentMapperDBToJSON Map the db struct to the json struct
func CommentMapperDBToJSON(
	comment models.CommentRow,
	user *models.UserRow, // TODO: convert to struct before passing in here.
	facebookInfo *models.FacebookInfoRow, // TODO: convert to struct before passing in here.
) structs.CommentStruct {
	var username string

	if len(user.Username) > 0 {
		username = user.Username
	} else {
		// TODO: this could be empty too
		username = facebookInfo.FacebookName
	}

	return structs.CommentStruct{
		UUID:         comment.UUID,
		FeedItemUUID: comment.FeedItemUUID,
		Message:      comment.Message,
		// TODO: disable these for now because they are not supported.
		// Lat:               comment.Lat,
		// Long:              comment.Long,
		// FormattedAddress:  comment.FormattedAddress,
		Username:          username,
		CreatedByUserUUID: user.UUID,
		CreatedAt:         comment.CreatedAt.Time,
		UpdatedAt:         comment.UpdatedAt.Time,
		DeletedAt:         comment.DeletedAt.Time,
	}
}
