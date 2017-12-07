package drivers

import (
	"io/ioutil"
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/bitbucket"
)

const bitbucketDriverName = "bitbucket"

func init() {
	registerDriver(bitbucketDriverName, BitbucketDefaultScopes, BitbucketUserFn, bitbucket.Endpoint, BitbucketAPIMap, BitbucketUserMap)
}

// BitbucketUserMap is the map to create the User struct
var BitbucketUserMap = map[string]string{
	"account_id":   "ID",
	"username":     "Username",
	"display_name": "FullName",
}

// BitbucketAPIMap is the map for API endpoints
var BitbucketAPIMap = map[string]string{
	"endpoint":      "https://api.bitbucket.org",
	"userEndpoint":  "/2.0/user",
	"emailEndpoint": "/2.0/user/emails",
}

// BitbucketUserFn is a callback to parse additional fields for User
var BitbucketUserFn = func(client *http.Client, u *structs.User) {
	// Set avatar
	u.Avatar = u.Raw["links"].(map[string]interface{})["avatar"].(map[string]interface{})["href"].(string)

	// Retrieve email
	req, err := client.Get(BitbucketAPIMap["endpoint"] + BitbucketAPIMap["emailEndpoint"])
	if err != nil {
		return
	}

	defer req.Body.Close()
	res, _ := ioutil.ReadAll(req.Body)
	data, err := jsonDecode(res)
	if err != nil {
		return
	}

	u.Email = data["values"].([]interface{})[0].(map[string]interface{})["email"].(string)
}

// BitbucketDefaultScopes contains the default scopes
var BitbucketDefaultScopes = []string{"account", "email"}
