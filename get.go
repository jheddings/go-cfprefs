package cfprefs

import "github.com/jheddings/go-cfprefs/internal"

func GetStr(appID, key string) (string, error) {
	return internal.Get(appID, key)
}

/* TODO
func GetInt(appID, key, user, host string) (int, error) {
	return internal.GetInt(appID, key, user, host)
}

func GetUInt(appID, key, user, host string) (uint, error) {
	return internal.GetInt(appID, key, user, host)
}

func GetBool(appID, key, user, host string) (bool, error) {
	return internal.GetBool(appID, key, user, host)
}

func GetFloat(appID, key, user, host string) (float64, error) {
	return internal.GetFloat(appID, key, user, host)
}
*/
