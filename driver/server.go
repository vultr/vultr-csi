/*
Copyright 2018 Vultr Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Defines Non blocking GRPC server interfaces
type NonBlockingGRPCServer interface {
	// Start services at the endpoint
	Start(endpoint string, ids csi.IdentityServer, cs csi.ControllerServer, ns csi.NodeServer)
	// Waits for the service to stop
	Wait()
	// Stops the service gracefully
	Stop()
	// Stops the service forcefully
	ForceStop()
}

// NewNonBlockingGRPCServer provides the non-blocking GRPC server
func NewNonBlockingGRPCServer() NonBlockingGRPCServer {
	return &nonBlockingGRPCServer{}
}

// NonBlocking server
type nonBlockingGRPCServer struct {
	wg     sync.WaitGroup
	server *grpc.Server
}

func (n *nonBlockingGRPCServer) Start(endpoint string, ids csi.IdentityServer, cs csi.ControllerServer, ns csi.NodeServer) {
	n.wg.Add(1)
	go n.serve(endpoint, ids, cs, ns)
}

func (n *nonBlockingGRPCServer) Wait() {
	n.wg.Wait()
}

func (n *nonBlockingGRPCServer) Stop() {
	n.server.GracefulStop()
}

func (n *nonBlockingGRPCServer) ForceStop() {
	n.server.Stop()
}

func (n *nonBlockingGRPCServer) serve(endpoint string, ids csi.IdentityServer, cs csi.ControllerServer, ns csi.NodeServer) {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(GRPCLogger),
	}

	serveURL, err := url.Parse(endpoint)
	if err != nil {
		log.Fatal(err.Error())
	}

	var addr string
	if serveURL.Scheme == "unix" {
		addr = serveURL.Path
		if errRemove := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			log.Fatalf("Failed to remove %s, error: %s", addr, errRemove.Error())
		}
	} else if serveURL.Scheme == "tcp" {
		addr = serveURL.Host
	} else {
		log.Fatalf("%v endpoint scheme not supported", serveURL.Scheme)
	}

	log.Infof("Start listening with scheme %v, addr %v", serveURL.Scheme, addr)
	listener, err := net.Listen(serveURL.Scheme, addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer(opts...)
	n.server = server

	if ids != nil {
		csi.RegisterIdentityServer(server, ids)
	}
	if cs != nil {
		csi.RegisterControllerServer(server, cs)
	}
	if ns != nil {
		csi.RegisterNodeServer(server, ns)
	}

	log.WithFields(log.Fields{
		"proto":   serveURL.Scheme,
		"address": addr,
	}).Infof("Listening for connections on address: %#v", listener.Addr())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	n.wg.Done()
}

// GRPCLogger provides better error handling for gRPC calls
func GRPCLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger := log.WithFields(log.Fields{
		"GRPC.call":    info.FullMethod,
		"GRPC.request": fmt.Sprintf("%+v", req),
	})

	resp, err := handler(ctx, req)
	if err != nil {
		logger.Errorf("GRPC error: %v", err)
	} else {
		logger.Infof("GRPC response: %+v", resp)
	}
	return resp, err
}
