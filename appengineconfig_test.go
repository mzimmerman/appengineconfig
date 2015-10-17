package appengineconfig_test

import "testing"
import "google.golang.org/appengine/aetest"
import "github.com/mzimmerman/appengineconfig"

func TestGet(t *testing.T) {
	c, closer, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()
	key, def, changed := "testme", "yes!", "changed"
	if want, got := def, appengineconfig.Get(c, key, def); want != got {
		t.Errorf("Expected %s, got %s", want, got)
	}
	if want, got := def, appengineconfig.Get(c, key, changed); want != got {
		t.Errorf("Expected %s, got %s", want, got)
	}
}
