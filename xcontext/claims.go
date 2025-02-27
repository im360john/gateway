package xcontext

import "context"

type claimsKeyType string

const claimsKey claimsKeyType = "claims"

func WithClaims(ctx context.Context, info map[string]any) context.Context {
	return context.WithValue(ctx, claimsKey, info)
}

func Claims(ctx context.Context) map[string]any {
	claims, ok := ctx.Value(claimsKey).(map[string]any)
	if !ok {
		return map[string]any{}
	}
	return claims
}
