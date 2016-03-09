package serialization

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	stats "github.com/milanaleksic/flowdock_stats/cmdcolors"
	"io/ioutil"
)

func GetKnownUsers() (catalog *Catalog) {
	catalog = &Catalog{Users: map[string]*Catalog_User{}}
	bytes, err := ioutil.ReadFile("users.dat")
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

func SaveUsers(catalog *Catalog) {
	bytes, err := proto.Marshal(catalog)
	if err != nil {
		stats.Warn(fmt.Sprint("Unable to save known users ", err))
	}
	err = ioutil.WriteFile("users.dat", bytes, 0644)
	if err != nil {
		stats.Warn(fmt.Sprint("Unable to save known users ", err))
	}
}
