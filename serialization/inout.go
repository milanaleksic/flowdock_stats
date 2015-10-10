package serialization

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	stats "milanaleksic.net/milanaleksic/flowdock-stats"
	"fmt"
)

func GetKnownUsers() (catalog *Catalog) {
	catalog = &Catalog{Users:map[string]*Catalog_User{}}
	bytes, err := ioutil.ReadFile("users.dat")
	if err != nil {
		stats.Warn("No users file found, will generate from beginning")
	} else {
		err = proto.Unmarshal(bytes, catalog)
		if err != nil {
			panic("No users found, will generate from beginning")
		}
	}
	return
}

func SaveUsers(catalog *Catalog) {
	bytes, err := proto.Marshal(&catalog)
	if err != nil {
		stats.Warn(fmt.Sprint("Unable to save known users ", err))
	}
	err = ioutil.WriteFile("users.dat", bytes, 0644)
	if err != nil {
		stats.Warn(fmt.Sprint("Unable to save known users ", err))
	}
}