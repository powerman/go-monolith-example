package grpc_test

import (
	"context"
	"crypto/tls"
	"net"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/smartystreets/goconvey/convey"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/powerman/go-monolith-example/api/proto/powerman/example/auth"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/ms/auth/internal/config"
	"github.com/powerman/go-monolith-example/ms/auth/internal/srv/grpc"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

func TestMain(m *testing.M) {
	reg := prometheus.NewPedanticRegistry()
	def.Init()
	app.InitMetrics(reg)
	grpc.InitMetrics(reg)
	cfg = config.MustGetServeTest()
	check.TestMain(m)
}

var (
	ctx = context.Background()
	cfg *config.ServeConfig
)

func testNew(t *check.C) (api.NoAuthSvcClient, api.AuthSvcClient, api.AuthIntSvcClient, *app.MockAppl) {
	t.Helper()
	ctrl := gomock.NewController(t)

	mockAppl := app.NewMockAppl(ctrl)

	ca, err := netx.LoadCACert(cfg.TLSCACert)
	t.Must(t.Nil(err))

	cert, err := tls.LoadX509KeyPair(cfg.TLSCert, cfg.TLSKey)
	t.Must(t.Nil(err))
	srv := grpc.NewServer(mockAppl, grpc.Config{
		Cert: &cert,
	})
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	t.Must(t.Nil(err, "net.Listen"))
	errc := make(chan error, 1)
	go func() { errc <- srv.Serve(ln) }()
	ctx, cancel := context.WithTimeout(context.Background(), def.TestTimeout)
	defer cancel()
	conn, err := grpcpkg.DialContext(ctx, strings.Replace(ln.Addr().String(), "127.0.0.1:", "localhost:", 1),
		grpcpkg.WithTransportCredentials(credentials.NewClientTLSFromCert(ca, "")),
		grpcpkg.WithBlock(),
	)
	t.Must(t.Nil(err, "grpc.Dial"))

	certInt, err := tls.LoadX509KeyPair(cfg.TLSCertInt, cfg.TLSKeyInt)
	t.Must(t.Nil(err))
	srvInt := grpc.NewServerInt(mockAppl, grpc.Config{
		Cert: &certInt,
	})
	lnInt, err := net.Listen("tcp", "127.0.0.1:0")
	t.Must(t.Nil(err, "net.Listen"))
	errcInt := make(chan error, 1)
	go func() { errcInt <- srvInt.Serve(lnInt) }()
	connInt, err := grpcpkg.DialContext(ctx, lnInt.Addr().String(),
		grpcpkg.WithTransportCredentials(credentials.NewClientTLSFromCert(ca, "")),
		grpcpkg.WithBlock(),
	)
	t.Must(t.Nil(err, "grpc.Dial internal"))

	t.Cleanup(func() {
		t.Helper()
		t.Nil(conn.Close())
		t.Nil(connInt.Close())
		srv.Stop()
		srvInt.Stop()
		t.Nil(<-errc, "srv.Serve")
		t.Nil(<-errcInt, "srv.Serve internal")
	})
	clientNoAuth := api.NewNoAuthSvcClient(conn)
	clientAuth := api.NewAuthSvcClient(conn)
	clientAuthInt := api.NewAuthIntSvcClient(connInt)
	return clientNoAuth, clientAuth, clientAuthInt, mockAppl
}
