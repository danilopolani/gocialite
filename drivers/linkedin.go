package drivers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/danilopolani/gocialite/structs"
    "golang.org/x/oauth2/linkedin"
)

const (
    linkedinDriverName = "linkedin"
)

func init() {
    registerDriver(linkedinDriverName, LinkedInDefaultScopes, LinkedInUserFn, linkedin.Endpoint, LinkedInAPIMap, LinkedInUserMap)
}

// LinkedInUserMap is the map to create the User struct
var LinkedInUserMap = map[string]string{
    "id":                 "ID",
    "localizedFirstName": "FirstName",
    "localizedLastName":  "LastName",
    "vanityName":         "_",
    "firstName":          "_",
    "lastName":           "_",
    "formattedName":      "_",
    "emailAddress":       "Email",
    "pictureUrl":         "Avatar",
}

// LinkedInAPIMap is the map for API endpoints
var LinkedInAPIMap = map[string]string{
    "endpoint":      "https://api.linkedin.com",
    "userEndpoint":  "/v2/me",
    "emailEndpoint": "/v2/emailAddress?q=members&projection=(elements*(handle~))",
}

// LinkedInUserFn is a callback to parse additional fields for User
var LinkedInUserFn = func(client *http.Client, u *structs.User) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Fprintln(os.Stderr, r)
        }
    }()
    u.FullName = u.FirstName + " " + u.LastName
    // Retrieve email
    req, err := client.Get(LinkedInAPIMap["endpoint"] + LinkedInAPIMap["emailEndpoint"])
    if err != nil {
        return
    }
    defer req.Body.Close()
    /*
       {
         "elements": [
           {
             "handle": "urn:li:emailAddress:3775708763",
             "handle~": {
               "emailAddress": "hsimpson@linkedin.com"
             }
           }
         ]
       }
    */
    email := make(map[string]interface{})
    err = json.NewDecoder(req.Body).Decode(&email)
    if err != nil {
        return
    }

    u.Email = email["elements"].([]interface{})[0].(map[string]interface{})["handle~"].(map[string]interface{})["emailAddress"].(string)

}

// LinkedInDefaultScopes contains the default scopes
var LinkedInDefaultScopes = []string{"r_emailaddress", "r_liteprofile"}
