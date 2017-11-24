package utils

import (
	// "github.com/aurieh/ddg-ng/commandclient"
	"encoding/json"
	"net/http"
)

func GetJson(res *http.Response, target interface{}) error {
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)
}

// func SafesearchDisabledChannel(client *commandclient.CommandClient, id string) (bool, error) {
// 	safesearch := true
// 	rows, err := client.DB.Query("SELECT * FROM channels WHERE id=?;", id)
// 	if err != nil {
// 		return false, err
// 	}
// 	for rows.Next() {
// 		err := rows.Scan(&safesearch)
// 		if err != nil {
// 			return false, err
// 		}
// 		if !safesearch {
// 			break
// 		}
// 	}
// 	err = rows.Close()
// 	if err != nil {
// 		return false, err
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		return false, err
// 	}
// 	return safesearch, nil
// }
