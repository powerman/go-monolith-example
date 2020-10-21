package mono

import "net/http"

func (s *Service) serveHealthCheck(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}
