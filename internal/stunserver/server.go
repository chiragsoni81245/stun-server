package stunserver

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/pion/stun"
	"github.com/rs/zerolog"
)

type Config struct {
	Addrs          []string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	RateLimitPerIP int
}

type Server struct {
	cfg     Config
	logger  zerolog.Logger
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	limiter *IPRateLimiter
	conns   []net.PacketConn
}

func New(cfg Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		cfg:     cfg,
		logger:  NewLogger(),
		ctx:     ctx,
		cancel:  cancel,
		limiter: NewIPRateLimiter(cfg.RateLimitPerIP, time.Second),
	}
}

func (s *Server) Start() error {
	for _, addr := range s.cfg.Addrs {
		conn, err := net.ListenPacket("udp4", addr)
		if err != nil {
			return err
		}

		s.conns = append(s.conns, conn)

		s.logger.Info().Msgf("Listening on %s", addr)

		go s.readLoop(conn)
	}

	return nil
}

func (s *Server) readLoop(conn net.PacketConn) {
	buf := make([]byte, 1500)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		_ = conn.SetReadDeadline(time.Now().Add(s.cfg.ReadTimeout))

		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		ip := addr.(*net.UDPAddr).IP.String()
		if !s.limiter.Allow(ip) {
			continue
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		s.wg.Add(1)
		go s.handle(conn, addr, data)
	}
}

func (s *Server) handle(conn net.PacketConn, addr net.Addr, data []byte) {
	// START timer immediately when request handling begins
	start := time.Now()
	defer func() {
		RequestDuration.Observe(time.Since(start).Seconds())
		s.wg.Done()
	}()

	msg := new(stun.Message)
	msg.Raw = data

	if err := msg.Decode(); err != nil {
		ErrorCount.Inc()
		return
	}

	RequestCount.WithLabelValues("binding").Inc()

	if msg.Type.Method != stun.MethodBinding ||
		msg.Type.Class != stun.ClassRequest {
		return
	}

	udpAddr := addr.(*net.UDPAddr)

	response := stun.MustBuild(
		stun.TransactionID,
		stun.BindingSuccess,
		&stun.XORMappedAddress{
			IP:   udpAddr.IP,
			Port: udpAddr.Port,
		},
	)

	_ = conn.SetWriteDeadline(time.Now().Add(s.cfg.WriteTimeout))
	_, _ = conn.WriteTo(response.Raw, addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("Shutdown initiated")

	// Stop accepting new work
	s.cancel()

	// Close all listeners
	for _, conn := range s.conns {
		_ = conn.Close()
	}

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info().Msg("Shutdown complete")
		return nil
	case <-ctx.Done():
		s.logger.Warn().Msg("Shutdown timed out")
		return ctx.Err()
	}
}
