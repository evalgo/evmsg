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
	"errors"
	echo "github.com/labstack/echo/v4"
)

var (
	ID     = "evlxc"
	Secret = "secret"
	Token  = "jEMjppkdwc1+8UkOplSJzhsUkRVqonQ1zmh+xC9Gj5w="
)

func Auth(msg *Message) error {
	if msg.Scope != "Login" {
		if msg.Token != Token {
			msg.Debug.Error = "the given Token is not valid!"
			msg.Data = []map[string]string{}
			return errors.New("the given Token is not valid!")
		}
	}
	switch msg.Scope {
	case "Login":
		switch msg.Command {
		case "getToken":
			msg.Debug.Info = "Login::getToken"
			if msg.ID == ID && msg.Secret == Secret {
				msg.Data = []map[string]string{map[string]string{"token": Token}}
			} else {
				msg.Debug.Error = "the given ID or Secret is not valid!"
				return errors.New("the given ID or Secret is not valid!")
			}
		}
	}
	return nil
}

func Bearer(c echo.Context) error {
	if c.Request().FormValue("token") == "" {
		switch c.Request().Header.Get("Authorization") {
		case "Bearer " + Token:
			return nil
		default:
			return errors.New("401 Unauthorized")
		}
	}
	if c.Request().FormValue("token") == Token {
		return nil
	}
	return errors.New("401 Unauthorized")
}

func CheckRequiredKeys(msg *Message, keys []string) error {
	// if no data was given and no keys
	if msg.Data == nil && len(keys) == 0 {
		// return ok
		return nil
	}
	data := msg.Data.([]interface{})
	// if no msg data and no keys are given
	if len(data) == 0 && len(keys) == 0 {
		// return ok
		return nil
	}
	// if no data is given but keys
	if len(data) == 0 && len(keys) > 0 {
		return errors.New("none of the required keys are given the data is empty!")
	}
	values := data[0].(map[string]interface{})
	for _, key := range keys {
		if _, ok := values[key]; !ok {
			return errors.New("the required key <" + key + "> was not found!")
		}
	}
	return nil
}
