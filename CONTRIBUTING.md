# Contributing

Creating a new driver for Gocialite is a lot simple, thank also to @vibbix.  

## Create a file
I suggest you to duplicate `drivers/bitbucket.go` since it's a complete code and name it with your provider name, ex: *myprovider.go*.

## Set variables

Change [line 12](https://github.com/danilopolani/gocialite/blob/master/drivers/bitbucket.go#L12) with:

```go
const myProviderDriverName = "myprovider"
```

It will be used in the `Driver()` function, example: `gocial.Driver("myprovider")`.  
Now on [line 19](https://github.com/danilopolani/gocialite/blob/master/drivers/bitbucket.go#L19) you have to create the mapping from API to populate the User struct.  
The relation is `"json_field_name": "StructFieldName"`, so if in our JSON there's a field called "first_name", it will be `"first_name": "FirstName"`.  

If there's some nested/complex field, please see the next chapter **User callback hook**.

Finally, on [lines 26-30](https://github.com/danilopolani/gocialite/blob/master/drivers/bitbucket.go#L26-L30) you have to fill the fields for the endpoint baseurl and the path of the user endpoint.  
In the case of Bitbucket, the email address is retrievable only from another endpoint, so we put in it also `emailEndpoint`, but usually you will need only `userEndpoint`.

If your provider has the user endpoint located to `https://api.myprovider.com/me`, the struct will be this:

```go
var MyProviderAPIMap = map[string]string{
	"endpoint":      "https://api.myprovider.com",
	"userEndpoint":  "/me",
}
```

Of course remember to **rename all the variables**. The ones that start with a capital letter are exported, so remember to write the first letter capitalized.

## User callback hook

If you have some complex field or you need to call some other endpoint in order to retrieve a field, you can do that in this section.  
In the case of Bitbucket, we use this hook to populate two fields: *avatar* from a nested array/map and *email* from another endpoint.  

The `client` variable is an `oAuth` client so it's already set up for oAuth details like `access_token`.

## Testing
Use the [example page](https://github.com/danilopolani/gocialite/wiki/Example) as starting point. Set up the credentials of your app in the `providerSecrets` variable, like:

```go
providerSecrets := map[string]map[string]string{
		...
		"myprovider": {
			"clientID":     "xxxxxxxxxxxxxx",
			"clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			"redirectURL":  "http://localhost:9090/auth/myprovider/callback",
		},
	}
```

And add the scopes for your app (or an empty slice) in the `providerScopes` variable:

```go
providerScopes := map[string][]string{
		...
		"myprovider": []string{}, // Or []string{"my_scope", "my_other_scope"}
	}
```

Now run `go run example.go` (or the name of your file), navigate to http://localhost:9090/auth/myprovider (or the name of your provider) and you will be redirected to the oAuth login.  
If everything works correctly, when you will be redirected to http://localhost:9090/auth/myprovider/callback, in your terminal you will see the content of the `User` struct populated (line `fmt.Printf("%#v", gocial.User)`).

## PR

Now that everything works, you can open a Pull Request, it will be tested and if it works, it will be merged and added to the README.
