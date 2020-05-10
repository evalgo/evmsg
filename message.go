/*
Copyright (c) 2020 EvAlgo Developers All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
this list of conditions and the following disclaimer.

THIS SOFTWARE IS PROVIDED BY EvAlgo Developers "AS IS" AND ANY
EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
IN NO EVENT SHALL EvAlgo Developers BE LIABLE FOR ANY DIRECT,
INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY
OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package evmsg

import (
	"encoding/json"
	"net/http"
)

type Debug struct {
	Error   string `json:"error"`
	Warning string `json:"warning"`
	Info    string `json:"info"`
}

type Message struct {
	Version string      `json:"version"`
	Scope   string      `json:"scope"`
	ID      string      `json:"id"`
	State   string      `json:"state"`
	Secret  string      `json:"secret"`
	Command string      `json:"command"`
	Token   string      `json:"token"`
	Debug   *Debug      `json:"debug"`
	Data    interface{} `json:"data"`
}

func MessageJSON(resp *Message) ([]byte, error) {
	return json.Marshal(resp)
}

func MessageHttpJSON(w http.ResponseWriter, j []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func NewMessage() *Message {
	return &Message{
		State: "Response",
		Debug: &Debug{},
		Data:  nil,
	}
}

func (msg *Message) Error() string {
	return msg.Debug.Error
}
