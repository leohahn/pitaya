// Copyright (c) TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package service

import (
	"errors"
	"testing"
	"time"
	"math/rand"

	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya/cluster"
	"github.com/topfreegames/pitaya/component"
	"github.com/topfreegames/pitaya/internal/codec"
	"github.com/topfreegames/pitaya/internal/message"
	"github.com/topfreegames/pitaya/serialize/json"
	"github.com/topfreegames/pitaya/session"
)

type MyComp struct {
	component.Base
}

func (m *MyComp) Init()                        {}
func (m *MyComp) Shutdown()                    {}
func (m *MyComp) Handler1(ss *session.Session) {}
func (m *MyComp) Handler2(ss *session.Session, b []byte) ([]byte, error) {
	return nil, nil
}

type NoHandlerRemoteComp struct {
	component.Base
}

func (m *NoHandlerRemoteComp) Init()     {}
func (m *NoHandlerRemoteComp) Shutdown() {}

func TestNewHandlerService(t *testing.T) {
	dieChan := make(chan bool)
	packetDecoder := codec.NewPomeloPacketDecoder()
	packetEncoder := codec.NewPomeloPacketEncoder()
	serializer := json.NewSerializer()
	heartbeatTimeout := 1 * time.Second
	messageEncoder := message.NewEncoder(rand.Int() % 2 == 0)
	sv := &cluster.Server{}
	remoteSvc := &RemoteService{}
	svc := NewHandlerService(
		dieChan,
		packetDecoder,
		packetEncoder,
		serializer,
		heartbeatTimeout,
		10, 9, 8,
		sv,
		remoteSvc,
		messageEncoder,
	)

	assert.NotNil(t, svc)
	assert.Equal(t, dieChan, svc.appDieChan)
	assert.Equal(t, packetDecoder, svc.decoder)
	assert.Equal(t, packetEncoder, svc.encoder)
	assert.Equal(t, serializer, svc.serializer)
	assert.Equal(t, heartbeatTimeout, svc.heartbeatTimeout)
	assert.Equal(t, 10, svc.messagesBufferSize)
	assert.Equal(t, sv, svc.server)
	assert.Equal(t, remoteSvc, svc.remoteService)
	assert.NotNil(t, svc.chLocalProcess)
	assert.NotNil(t, svc.chRemoteProcess)
}

func TestHandlerServiceRegister(t *testing.T) {
	svc := NewHandlerService(nil, nil, nil, nil, 0, 0, 0, 0, nil, nil, nil)
	err := svc.Register(&MyComp{}, []component.Option{})
	assert.NoError(t, err)
	defer func() { handlers = make(map[string]*component.Handler, 0) }()
	assert.Len(t, svc.services, 1)
	val, ok := svc.services["MyComp"]
	assert.True(t, ok)
	assert.NotNil(t, val)
	assert.Len(t, handlers, 2)
	val2, ok := handlers["MyComp.Handler1"]
	assert.True(t, ok)
	assert.NotNil(t, val2)
	val2, ok = handlers["MyComp.Handler2"]
	assert.True(t, ok)
	assert.NotNil(t, val2)
}

func TestHandlerServiceRegisterFailsIfRegisterTwice(t *testing.T) {
	svc := NewHandlerService(nil, nil, nil, nil, 0, 0, 0, 0, nil, nil, nil)
	err := svc.Register(&MyComp{}, []component.Option{})
	assert.NoError(t, err)
	err = svc.Register(&MyComp{}, []component.Option{})
	assert.Contains(t, err.Error(), "handler: service already defined")
}

func TestHandlerServiceRegisterFailsIfNoHandlerMethods(t *testing.T) {
	svc := NewHandlerService(nil, nil, nil, nil, 0, 0, 0, 0, nil, nil, nil)
	err := svc.Register(&NoHandlerRemoteComp{}, []component.Option{})
	assert.Equal(t, errors.New("type NoHandlerRemoteComp has no exported methods of suitable type"), err)
}

func TestHandlerServiceProcessMessage(t *testing.T) {
	// TODO
}

func TestHandlerServiceLocalProcess(t *testing.T) {
	// TODO
}

func TestHandlerServiceProcessPacket(t *testing.T) {
	// TODO
}

func TestHandlerServiceHandle(t *testing.T) {
	// TODO
}

func TestHandlerServiceDispatch(t *testing.T) {
	// TODO
}
