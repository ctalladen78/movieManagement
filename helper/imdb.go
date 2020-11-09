package helper

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/gommon/log"
)

// ImdbByTitle . . .
func ImdbByTitle() string {

	url := "https://movie-database-imdb-alternative.p.rapidapi.com/?i=tt4154796&r=json"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-host", "movie-database-imdb-alternative.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "c1a06d7bdemsh67eb386c3da090ep115cc3jsn6049bf986f74")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Error(err)
		return ""
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
	return string(body)

}
