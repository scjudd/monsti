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

package mail

import (
	"github.com/chrneumann/mimemail"
	"github.com/monsti/rpc/types"
	"io"
	"log"
	"net/rpc"
	"net/url"
	"os"
	"strings"
)

// Send given Mail.
//
// An empty From or To field will be filled with the site owner's name and
// address.
func (s *Service) SendMail(m mimemail.Mail) {
	var reply int
	err := s.Call("NodeRPC.SendMail", m, &reply)
	if err != nil {
		s.Logger.Fatal("master: RPC SendMail error:", err) // 
	}
