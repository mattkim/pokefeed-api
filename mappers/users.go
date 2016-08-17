package mappers

import (
	"github.com/pokefeed/pokefeed-api/models"
	"github.com/pokefeed/pokefeed-api/structs"
)

// FacebookUserMapperDBToJSON Map the db struct to the json struct
func FacebookUserMapperDBToJSON(
	user models.UserRow,
	facebookInfo models.FacebookInfoRow,
) structs.UserStruct {
	var username string

	if len(user.Username) > 0 {
		username = user.Username
	} else if len(facebookInfo.FacebookName) > 0 {
		username = facebookInfo.FacebookName
	} else {
		username = ""
	}

	return structs.UserStruct{
		UUID:     user.UUID,
		Username: username,
		Email:    user.Email,
	}
}
