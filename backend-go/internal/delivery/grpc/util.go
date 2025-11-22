package grpc

func fromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
