package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRootRedirectsToDownloadPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &ginService{server: gin.New()}
	svc.SetupRoutes()

	rec := httptest.NewRecorder()
	svc.server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusFound {
		t.Fatalf("状态码 = %d，期望 %d", rec.Code, http.StatusFound)
	}
	if got := rec.Header().Get("Location"); got != "/download/page/" {
		t.Fatalf("Location = %q，期望 /download/page/", got)
	}
}
