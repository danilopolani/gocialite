package drivers

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"github.com/danilopolani/gocialite/drivers/option"
	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/odnoklassniki"
)

// Odnoklassniki is a social network service for classmates and old friends.
// It is popular in Russia and former Soviet Republics.
const okDriverName = "ok"

func init() {
	err := RegisterDriver(
		option.Driver(okDriverName),
		option.DefaultScopes(okDefaultScopes),
		option.Callback(okUserFn),
		option.Endpoint(odnoklassniki.Endpoint),
		option.APIMap(okAPIMap),
		option.UserMap(okUserMap),
		option.Sig(okSigFn),
	)
	if err != nil {
		panic(err)
	}
}

// okUserMap is the map to create the User struct
var okUserMap = map[string]string{
	"last_name":  "LastName",
	"first_name": "FirstName",
	"name":       "FullName",
	"pic_full":   "Avatar",
	"uid":        "ID",
}

// okAPIMap is the map for API endpoints
var okAPIMap = map[string]string{
	"endpoint": "https://api.ok.ru",
	"userEndpoint": "/fb.do" +
		"?access_token=%ACCESS_TOKEN" +
		"&application_key=%APPLICATION_KEY" +
		"&fields=HAS_EMAIL,EMAIL,FIRST_NAME,LAST_NAME,NAME,PIC_FULL,UID" +
		"&format=json" +
		"&method=users.getCurrentUser" +
		"&sig=%SIG",
}

// okUserFn is a callback to parse additional fields for User
var okUserFn = func(client *http.Client, u *structs.User) {
	if u.Raw["has_email"].(bool) && u.Raw["email"] != nil {
		u.Email = u.Raw["email"].(string)
	}
}

// okDefaultScopes contains the default scopes
var okDefaultScopes = []string{"GET_EMAIL"}

// okSigFn contains the signature function
var okSigFn option.SigFunc = func(
	conf *oauth2.Config,
	tok *oauth2.Token,
	params map[string]string,
) string {
	h := md5.New()
	h.Write([]byte(tok.AccessToken))
	h.Write([]byte(conf.ClientSecret))
	sessionKey := hex.EncodeToString(h.Sum(nil))
	h = md5.New()
	h.Write([]byte("application_key=" + params["%APPLICATION_KEY"]))
	h.Write([]byte("fields=HAS_EMAIL,EMAIL,FIRST_NAME,LAST_NAME,NAME,PIC_FULL,UID"))
	h.Write([]byte("format=json"))
	h.Write([]byte("method=users.getCurrentUser"))
	h.Write([]byte(sessionKey))
	return hex.EncodeToString(h.Sum(nil))
}
