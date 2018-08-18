package drivers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/danilopolani/gocialite/drivers/option"
	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/vk"
)

// VK is a Russian online social media and social networking service
const vkDriverName = "vk"

func init() {
	err := RegisterDriver(
		option.Driver(vkDriverName),
		option.DefaultScopes(vkDefaultScopes),
		option.Callback(vkUserFn),
		option.Endpoint(vk.Endpoint),
		option.APIMap(vkAPIMap),
		option.UserMap(vkUserMap),
	)
	if err != nil {
		panic(err)
	}
}

// vkUserMap is the map to create the User struct
var vkUserMap = map[string]string{}

// vkAPIMap is the map for API endpoints
var vkAPIMap = map[string]string{
	"endpoint": "https://api.vk.com",
	"userEndpoint": "/method/users.get" +
		"?fields=photo_max_orig,first_name,last_name,uid,screen_name" +
		"&v=5.21&access_token=%ACCESS_TOKEN",
}

// vkUserFn is a callback to parse additional fields for User
var vkUserFn = func(client *http.Client, u *structs.User) {
	resp := u.Raw["response"].([]interface{})
	firstResp := resp[0].(map[string]interface{})
	u.ID = strconv.FormatFloat(firstResp["id"].(float64), 'E', -1, 64)
	u.FirstName = firstResp["first_name"].(string)
	u.LastName = firstResp["last_name"].(string)
	u.Avatar = firstResp["photo_max_orig"].(string)
	u.Username = firstResp["screen_name"].(string)
	u.FullName = fmt.Sprint(u.FirstName, " ", u.LastName)
}

// vkDefaultScopes contains the default scopes
var vkDefaultScopes = []string{""}
