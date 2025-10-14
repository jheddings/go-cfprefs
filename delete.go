package cfprefs

import "github.com/jheddings/go-cfprefs/internal"

func Delete(appID, key string) error {
	return internal.Delete(appID, key)
}
