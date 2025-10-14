package cfprefs

import "github.com/jheddings/go-cfprefs/internal"

func SetStr(appID, key, value string) error {
	return internal.Set(appID, key, value)
}
