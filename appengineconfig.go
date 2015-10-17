package appengineconfig

import (
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

var ConfigName = "AppEngineConfig" // default key-prefix to utilize

var configCache map[string]string
var configLock sync.Mutex

func init() {
	configCache = make(map[string]string)
}

func Get(c context.Context, key, def string) string {
	configLock.Lock()
	defer configLock.Lock()
	value, ok := configCache[key]
	if ok {
		return value
	}
	dskey := datastore.NewKey(c, ConfigName, key, 0, nil)
	_, err := datastore.Get(c, dskey, &value)
	if err == datastore.ErrNoSuchEntity {
		_, err = datastore.Put(c, dskey, def)
		if err != nil {
			log.Infof(c, "Creating default config for key - %s - default value is - %s", key, def)
		}
		return def // return default, totally new config setting
	}
	if err != nil {
		log.Errorf(c, "Error fetching config for key - %s - error is - %v", key, err)
		return def // error, return the default
	}
	configCache[key] = value
	return value
}
