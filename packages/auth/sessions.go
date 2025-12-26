package auth

func extractSessionID() string {
	return "session_123456"
}
func GetSessionInfo()string {
	return extractSessionID()
}