package drivers

import (
    "encoding/json"
    "fmt"
    "net/http"

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
    "id": "ID",
}

// LinkedInAPIMap is the map for API endpoints
var LinkedInAPIMap = map[string]string{
    "endpoint":      "https://api.linkedin.com",
    "userEndpoint":  "/v2/me",
    "emailEndpoint": "/v2/emailAddress?q=members&projection=(elements*(handle~))",
}

// LinkedInUserFn is a callback to parse additional fields for User
var LinkedInUserFn = func(client *http.Client, u *structs.User) {
    /*
       {
          "id":"REDACTED",
          "firstName":{
             "localized":{
                "en_US":"Tina"
             },
             "preferredLocale":{
                "country":"US",
                "language":"en"
             }
          },
          "lastName":{
             "localized":{
                "en_US":"Belcher"
             },
             "preferredLocale":{
                "country":"US",
                "language":"en"
             }
          },
           "profilePicture": {
               "displayImage": "urn:li:digitalmediaAsset:B54328XZFfe2134zTyq"
           }
       }
    */
    raw := u.Raw
    if raw != nil {
        fn := raw["firstName"].(map[string]interface{})
        localized := fn["localized"].(map[string]interface{})
        preferredLocale := fn["preferredLocale"].(map[string]interface{})
        u.FirstName = localized[fmt.Sprintf("%s_%s", preferredLocale["language"],
            preferredLocale["country"])].(string)
        ln := raw["lastName"].(map[string]interface{})
        localized = ln["localized"].(map[string]interface{})
        preferredLocale = ln["preferredLocale"].(map[string]interface{})
        u.LastName = localized[fmt.Sprintf("%s_%s", preferredLocale["language"],
            preferredLocale["country"])].(string)
    }

    // Retrieve email
    req, err := client.Get(LinkedInAPIMap["endpoint"] + LinkedInAPIMap["emailEndpoint"])
    if err != nil {
        return
    }
    defer req.Body.Close()
    /*
       {
           "handle": "urn:li:emailAddress:3775708763",
           "handle~": {
               "emailAddress": "hsimpson@linkedin.com"
           }
       }
    */
    email := make(map[string]interface{})
    err = json.NewDecoder(req.Body).Decode(&email)
    if err != nil {
        return
    }

    fmt.Printf("%+v", email)
    // u.Email = email.Handle.EmailAddress

}

// LinkedInDefaultScopes contains the default scopes
var LinkedInDefaultScopes = []string{"r_emailaddress", "r_liteprofile", "w_member_social"}
