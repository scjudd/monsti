// This file is part of Monsti, a web content management system.
// Copyright 2012-2013 Christian Neumann
// 
// Monsti is free software: you can redistribute it and/or modify it under the
// terms of the GNU Affero General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// Monsti is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
// A PARTICULAR PURPOSE.  See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Monsti.  If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// Service represents the rpc connection to a service.
type Service struct {
	Client *rpc.Client
}

// Connect establishes a new RPC connection to the given service.
//
// path is the unix domain socket path to the service.
func (s *Service) Connect(path string) error {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return err
	}
	s.Client = rpc.NewClient(conn)
	return nil
}

type Provider struct {
	Logger *log.Logger
}

// Serve listens on the given unix domain socket path for incoming rpc
// connections.
//
// service is the service name
func (p *Provider) Serve(path string, service string, rcvr interface{}) error {
	listener, err := net.Listen("unix", path)
	if err != nil {
		return fmt.Errorf("Could not listen on unix domain socket %q: %v",
			path, err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("Could not accept connection on %q: %v",
				path, err)
		}
		server := rpc.NewServer()
		if err = server.RegisterName(service, rcvr); err != nil {
			return fmt.Errorf("Could not register RPC methods: %v",
				err.Error())
		}
		go server.ServeConn(conn)
	}
	return nil
}
