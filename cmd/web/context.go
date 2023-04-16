package main
type contextKey string
const isAuthenticatedContextKey = contextKey("isAuthenticated")
const isFacultyContextKey = contextKey("isFaculty")
const isAuthorizedContextKey = contextKey("isAuthorized")
