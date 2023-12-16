package utils

import "context"

var attemptKey = struct{}{}

func ContextWithAttempt(ctx context.Context, attempt int) context.Context {
	return context.WithValue(ctx, attemptKey, attempt)
}

func GetAttempt(ctx context.Context) int {
	if attempt, ok := ctx.Value(attemptKey).(int); ok {
		return attempt
	}
	return 0
}
