package logger

import (
	"fmt"
	"net"
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BuilderOption struct {
	UDPIP       string
	UDPPort     int
	PrettyPrint bool
}

func getLogLevel(logLevel string) zap.AtomicLevel {
	switch logLevel {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	}
}

func NewLogger(serviceName string, level string, options ...func(*BuilderOption)) *zap.Logger {
	// build config
	cfg := &BuilderOption{}
	for _, option := range options {
		option(cfg)
	}

	// create multiple sync target if UDP logging is enabled
	syncer := zapcore.AddSync(os.Stdout)
	if cfg.UDPIP != "" && cfg.UDPPort > 0 {
		syncer = zapcore.NewMultiWriteSyncer(os.Stdout, newUDPSyncer(cfg.UDPIP, cfg.UDPPort))
	}

	// create new core with log duplication
	var core zapcore.Core
	if cfg.PrettyPrint {
		// create new console formatter config with colored level
		config := zap.NewDevelopmentEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(config), syncer, getLogLevel(level))
	} else {
		// create new ECS formatter config
		core = ecszap.NewCore(ecszap.NewDefaultEncoderConfig(), syncer, getLogLevel(level))
	}

	// create new log instance
	log := zap.New(core, zap.AddCaller())
	log = log.With(zap.String("service_name", serviceName))

	return log
}

func WithID(log *zap.Logger, contextName string, scopeName string) *zap.Logger {
	return log.With(zap.String("context", contextName), zap.String("scope", scopeName))
}

type UDPSyncer struct {
	conn *net.UDPConn
}

func newUDPSyncer(bindIP string, bindPort int) *UDPSyncer {
	// ResolveUDPAddr returns an address of UDP end point.
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", bindIP, bindPort))
	if err != nil {
		fmt.Println("Failed to resolve address", err)
	}

	// DialUDP connects to the remote address raddr on the network net
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Failed to dial address", err)
	}

	return &UDPSyncer{conn: conn}
}

func (s *UDPSyncer) Write(p []byte) (n int, err error) {
	return s.conn.Write(p)
}

func (s *UDPSyncer) Sync() error {
	return nil
}
