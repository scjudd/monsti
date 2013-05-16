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

package node

import (
	"fmt"
	"github.com/monsti/service"
	"github.com/monsti/service/login"
	"net/url"
	"strings"
)

// Service represents the RPC connection to the Nodes service.
type Service struct {
	*service.Service
}

// NewService returns a new Service.
func NewService() *Service {
	var service_ Service
	service_.Service = new(service.Service)
	return &service_
}

type Node struct {
	Path string "path,omitempty"
	// Content type of the node.
	Type        string
	Title       string
	ShortTitle  string
	Description string
	// Order of the node compared to its siblings.
	Order int
	// Don't show the node in navigations if Hide is true.
	Hide bool
}

// PathToID returns an ID for the given node based on it's path.
//
// The ID is simply the path of the node with all slashes replaced by two
// underscores and the result prefixed with "node-".
//
// PathToID will panic if the path is not set.
//
// For example, a node with path "/foo/bar" will get the ID "node-__foo__bar".
func (n Node) PathToID() string {
	if len(n.Path) == 0 {
		panic("Can't calculate ID of node with unset path.")
	}
	return "node-" + strings.Replace(n.Path, "/", "__", -1)
}

// A request to be processed by a nodes service.
type Request struct {
	// The requested node.
	Node Node
	// The query values of the request URL.
	Query url.Values
	// Method of the request (GET,POST,...).
	Method string
	// User session
	Session login.Session
	// Action to perform (e.g. "edit").
	Action string
}

// Response to a node request.
type Response struct {
	// The html content to be embedded in the root template.
	Body []byte
	// Raw must be set to true if Body should not be embedded in the root
	// template. The content type will be automatically detected.
	Raw bool
	// If set, redirect to this target using error 303 'see other'.
	Redirect string
	// The node as received by GetRequest, possibly with some fields
	// updated (e.g. modified title).
	//
	// If nil, the original node data is used.
	Node *Node
}

// Write appends the given bytes to the body of the response.
func (r *Response) Write(p []byte) (n int, err error) {
	r.Body = append(r.Body, p...)
	return len(p), nil
}

// Request performs the given request.
func (s *Service) Request(req *Request) (*Response, error) {
	var res Response
	err := s.Client.Call("Node.Request", req, &res)
	if err != nil {
		return nil, fmt.Errorf("node: RPC error for Request:", err)
	}
	return &res, nil
}

// GetNodeType returns all supported node types.
func (s *Service) GetNodeTypes() ([]string, error) {
	var res []string
	err := s.Client.Call("Node.GetNodeTypes", 0, &res)
	if err != nil {
		return nil, fmt.Errorf("node: RPC error for Request:", err)
	}
	return res, nil
}

/*

func (m *NodeRPC) GetRequest(arg int, reply *client.Request) error {
	if m.Worker.Ticket != nil {
		return errors.New("monsti: Still waiting for response to last request.")
	}
	ticket := <-m.Worker.Tickets
	m.Worker.Ticket = &ticket
	request := client.Request{
		Method:  m.Worker.Ticket.Request.Method,
		Node:    m.Worker.Ticket.Node,
		Query:   m.Worker.Ticket.Request.URL.Query(),
		Session: m.Worker.Ticket.Session,
		Action:  m.Worker.Ticket.Action}
	*reply = request
	return nil
}

func (m *NodeRPC) SendResponse(res client.Response, reply *int) error {
	m.Worker.Ticket.ResponseChan <- res
	m.Worker.Ticket = nil
	return nil
}
*/
