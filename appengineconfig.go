package appengineconfig

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

var ConfigName = "AppEngineConfig" // default key-prefix to utilize

type Value struct {
	Val string
}

var configCache map[string]string
var configLock sync.Mutex

func init() {
	configCache = make(map[string]string)
}

func Get(c context.Context, key, def string, values ...interface{}) string {
	configLock.Lock()
	defer configLock.Unlock()
	val, ok := configCache[key]
	if ok {
		return val
	}
	value := Value{}
	dskey := datastore.NewKey(c, ConfigName, key, 0, nil)
	err := datastore.Get(c, dskey, &value)
	if err == datastore.ErrNoSuchEntity {
		value.Val = def
		log.Infof(c, "Creating default config for key - %s - default value is - %s", key, value.Val)
		_, err = datastore.Put(c, dskey, &value)
		if err != nil {
			log.Errorf(c, "Error creating default config for key - %s - error is - %v", key, err)
		}
		return fmt.Sprintf(def, values...) // return default, totally new config setting
	}
	if err != nil {
		log.Errorf(c, "Error fetching config for key - %s - error is - %v", key, err)
		return def // error, return the default
	}
	configCache[key] = value.Val
	return fmt.Sprintf(value.Val, values...)
}
