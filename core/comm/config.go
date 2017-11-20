/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	// Is the configuration cached?
	configurationCached = false
	// Is TLS enabled
	tlsEnabled bool
	// Max send and receive bytes for grpc clients and servers
	maxRecvMsgSize = 100 * 1024 * 1024
	maxSendMsgSize = 100 * 1024 * 1024
	// Default peer keepalive options
	keepaliveOptions = KeepaliveOptions{
		ClientKeepaliveTime:    60,   // 1 min
		ClientKeepaliveTimeout: 20,   // 20 sec - gRPC default
		ServerKeepaliveTime:    7200, // 2 hours - gRPC default
		ServerKeepaliveTimeout: 20,   // 20 sec - gRPC default
	}
	// chaincode keepalive options separate from peer keepalive
	// options above (for flexibility)
	chaincodeKeepaliveOptions = KeepaliveOptions{
		ClientKeepaliveTime:    60, // 1 min
		ClientKeepaliveTimeout: 20, // 20 sec - gRPC default
		ServerKeepaliveTime:    60, // 1 min
		ServerKeepaliveTimeout: 20, // 20 sec - gRPC default
	}
)

// ServerConfig defines the parameters for configuring a GRPCServer instance
type ServerConfig struct {
	// SecOpts defines the security parameters
	SecOpts *SecureOptions
	// KaOpts defines the keepalive parameters
	KaOpts *KeepaliveOptions
}

// SecureOptions defines the security parameters (e.g. TLS) for a
// GRPCServer instance
type SecureOptions struct {
	//PEM-encoded X509 public key to be used by the server for TLS communication
	ServerCertificate []byte
	//PEM-encoded private key to be used by the server for TLS communication
	ServerKey []byte
	//Set of PEM-encoded X509 certificate authorities to optionally send
	//as part of the server handshake
	ServerRootCAs [][]byte
	//Set of PEM-encoded X509 certificate authorities to use when verifying
	//client certificates
	ClientRootCAs [][]byte
	//Whether or not to use TLS for communication
	UseTLS bool
	//Whether or not TLS client must present certificates for authentication
	RequireClientCert bool
}

// KeepAliveOptions is used to set the gRPC keepalive settings for both
// clients and servers
type KeepaliveOptions struct {
	// ClientKeepaliveTime is the duration in seconds after which if the client
	// does not see any activity from the server it pings the server to see
	// if it is alive
	ClientKeepaliveTime int
	// ClientKeepaliveTimeout is the duration the client waits for a response
	// from the server after sending a ping before closing the connection
	ClientKeepaliveTimeout int
	// ServerKeepaliveTime is the duration in seconds after which if the server
	// does not see any activity from the client it pings the client to see
	// if it is alive
	ServerKeepaliveTime int
	// ServerKeepaliveTimeout is the duration the server waits for a response
	// from the client after sending a ping before closing the connection
	ServerKeepaliveTimeout int
}

// cacheConfiguration caches common package scoped variables
func cacheConfiguration() {
	if !configurationCached {
		tlsEnabled = viper.GetBool("peer.tls.enabled")
		configurationCached = true
	}
}

// TLSEnabled return cached value for "peer.tls.enabled" configuration value
func TLSEnabled() bool {
	if !configurationCached {
		cacheConfiguration()
	}
	return tlsEnabled
}

// MaxRecvMsgSize returns the maximum message size in bytes that gRPC clients
// and servers can receive
func MaxRecvMsgSize() int {
	return maxRecvMsgSize
}

// SetMaxRecvMsgSize sets the maximum message size in bytes that gRPC clients
// and servers can receive
func SetMaxRecvMsgSize(size int) {
	maxRecvMsgSize = size
}

// MaxSendMsgSize returns the maximum message size in bytes that gRPC clients
// and servers can send
func MaxSendMsgSize() int {
	return maxSendMsgSize
}

// SetMaxSendMsgSize sets the maximum message size in bytes that gRPC clients
// and servers can send
func SetMaxSendMsgSize(size int) {
	maxSendMsgSize = size
}

// SetKeepaliveOptions sets the gRPC keepalive options for both clients and
// servers
func SetKeepaliveOptions(ka KeepaliveOptions) {
	keepaliveOptions = ka
}

// ServerKeepaliveOptions returns the gRPC keepalive options for servers
func ServerKeepaliveOptions() []grpc.ServerOption {
	return serverKeepaliveOptionsWithKa(&keepaliveOptions)
}

func serverKeepaliveOptionsWithKa(ka *KeepaliveOptions) []grpc.ServerOption {
	var serverOpts []grpc.ServerOption
	kap := keepalive.ServerParameters{
		Time:    time.Duration(ka.ServerKeepaliveTime) * time.Second,
		Timeout: time.Duration(ka.ServerKeepaliveTimeout) * time.Second,
	}
	serverOpts = append(serverOpts, grpc.KeepaliveParams(kap))
	kep := keepalive.EnforcementPolicy{
		// needs to match clientKeepalive
		MinTime: time.Duration(ka.ClientKeepaliveTime) * time.Second,
		// allow keepalive w/o rpc
		PermitWithoutStream: true,
	}
	serverOpts = append(serverOpts, grpc.KeepaliveEnforcementPolicy(kep))
	return serverOpts
}

// ClientKeepaliveOptions returns the gRPC keepalive options for clients
func ClientKeepaliveOptions() []grpc.DialOption {
	return clientKeepaliveOptionsWithKa(&keepaliveOptions)
}

func clientKeepaliveOptionsWithKa(ka *KeepaliveOptions) []grpc.DialOption {
	var dialOpts []grpc.DialOption
	kap := keepalive.ClientParameters{
		Time:                time.Duration(ka.ClientKeepaliveTime) * time.Second,
		Timeout:             time.Duration(ka.ClientKeepaliveTimeout) * time.Second,
		PermitWithoutStream: true,
	}
	dialOpts = append(dialOpts, grpc.WithKeepaliveParams(kap))
	return dialOpts
}
