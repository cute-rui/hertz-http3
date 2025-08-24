/*
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package quic

import (
	"context"
	"net"

	"github.com/cloudwego/hertz/pkg/network"
	quicgo "github.com/quic-go/quic-go"
)

var _ network.StreamConn = &conn{}

type conn struct {
	rawConn interface{}
	c *quicgo.Conn
}

//Dirty hack
type ctx struct {
	context.Context
	c *quicgo.Conn
}

func (c *ctx) Done() <-chan struct{} {
	return c.c.HandshakeComplete()
}

func (c *conn) GetVersion() uint32 {
	return uint32(c.c.ConnectionState().Version)
}

func (c *conn) GetRawConnection() interface{} {
	return c.rawConn
}

func (c *conn) AcceptStream(ctx context.Context) (network.Stream, error) {
	stream, err := c.c.AcceptStream(ctx)
	return newStream(stream), err
}

func (c *conn) AcceptUniStream(ctx context.Context) (network.ReceiveStream, error) {
	stream, err := c.c.AcceptUniStream(ctx)
	return newReadStream(stream), err
}

func (c *conn) OpenStream() (network.Stream, error) {
	stream, err := c.c.OpenStream()
	return newStream(stream), err
}

func (c *conn) OpenStreamSync(ctx context.Context) (network.Stream, error) {
	stream, err := c.c.OpenStreamSync(ctx)
	return newStream(stream), err
}

func (c *conn) OpenUniStream() (network.SendStream, error) {
	stream, err := c.c.OpenUniStream()
	return newWriteStream(stream), err
}

func (c *conn) OpenUniStreamSync(ctx context.Context) (network.SendStream, error) {
	stream, err := c.c.OpenUniStreamSync(ctx)
	return newWriteStream(stream), err
}

func (c *conn) CloseWithError(err network.ApplicationError, errMsg string) error {
	return c.c.CloseWithError(quicgo.ApplicationErrorCode(err.ErrCode()), errMsg)
}

func (c *conn) Context() context.Context {
	return c.c.Context()
}

func (c *conn) HandshakeComplete() context.Context {
	return &ctx{context.Background(), c.c}
}

func (c *conn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

func (c *conn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

func newStreamConn(qc *quicgo.Conn) *conn {
	return &conn{qc, qc}
}
