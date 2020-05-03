package drivers

import (
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
    "id":            "ID",
    "firstName":     "FirstName",
    "lastName":      "LastName",
    "profilePicture":    "Avatar",
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
    u.FirstName = u.FirstName["localized"][fmt.Sprintf("%s_%s", u.FirstName["preferredLocale"]["language"],
        u.FirstName["preferredLocale"]["country"])]
    u.LastName = u.LastName["localized"][fmt.Sprintf("%s_%s", u.LastName["preferredLocale"]["language"],
        u.LastName["preferredLocale"]["country"])]

    u.Avatar = u.Avatar["profilePicture"]["displayImage"]

    // Retrieve email
    req, err := client.Get(LinkedinAPIMap["endpoint"] + LinkedinAPIMap["emailEndpoint"])
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
    type EmailRes struct {
        Handle struct {
            EmailAddress string `json:"emailAddress"`
        } `json:"handle~"`
    }
    email := new(EmailRes)
    err = json.NewDecoder(req.Body).Decode(&email)
    if err != nil {
        return
    }

    u.Email = email.Handle.EmailAddress

}

// LinkedInDefaultScopes contains the default scopes
var LinkedInDefaultScopes = []string{"r_emailaddress", "r_liteprofile", "w_member_social"}
