package main

// ts := oauth2.StaticTokenSource(
// 	&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
// )
// tc := oauth2.NewClient(oauth2.NoContext, ts)
// client := github.NewClient(tc)
// user := os.Args[1]
// log.Println("Gathering data for", user)

// followers, err := followers.Get(user, client)
// if err != nil {
// 	log.Fatalln(err)
// }

// log.Println("You have a total of", len(followers), "followers!")
// for _, follower := range followers {
// 	log.Println(*follower.Login)
// }
