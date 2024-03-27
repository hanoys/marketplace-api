package domain

import "context"

// TODO: interface for token provider

type TokenProvider interface {
	CreateSession(ctx context.Context, tokenString string)
	CloseSession(ctx context.Context, tokenString string)
	RefreshSession(ctx context.Context, refreshTokenString string)
	VerifyToken(ctx context.Context, tokenString string)
}
