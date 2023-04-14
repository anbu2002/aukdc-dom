package main
type contextKey string
const isAuthenticatedContextKey = contextKey("isAuthenticated")
const isAuthorizedContextKey = contextKey("isAuthorized")
