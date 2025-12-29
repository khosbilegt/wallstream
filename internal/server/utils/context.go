package utils

import "context"

func GetStringFromContext(ctx context.Context, key any) (string, bool) {
	v := ctx.Value(key)
	s, ok := v.(string)
	return s, ok && s != ""
}
