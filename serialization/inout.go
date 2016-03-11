package serialization

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	stats "github.com/milanaleksic/flowdock_stats/cmdcolors"
	"io/ioutil"
)

const cacheFile = "users.dat"

/*
GetKnownUsers gets from local cache file all known users.
In case cache file doesn't exist, it will return empty catalog
*/
func GetKnownUsers() (catalog *Catalog) {
	catalog = &Catalog{Users: map[string]*Catalog_User{}}
	bytes, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		stats.Warn("No users file found, will generate from beginning")
	} else {
		err = proto.Unmarshal(bytes, catalog)
		if err != nil {
			panic("Users file users.dat found, but it's nor parsable! Try deleting the file.")
		}
	}
	return
}

/*
SaveUsers is able to save the catalog to the default users cache file - users.dat
*/
func SaveUsers(catalog proto.Message) {
	bytes, err := proto.Marshal(catalog)
	if err != nil {
		stats.Warn(fmt.Sprint("Unable to save known users ", err))
	}
	err = ioutil.WriteFile(cacheFile, bytes, 0644)
	if err != nil {
		stats.Warn(fmt.Sprint("Unable to save known users ", err))
	}
}
