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
	"testing"
)

func Test_Unit_MessageJSONMarshal(t *testing.T) {
	respJSON := `{"Version":"","Scope":"","ID":"","State":"","Secret":"","Command":"","Debug":{"Error":"","Warning":"","Info":""},"Data":["evalgo-org","gitea","airflow","bo-dyndns","caldav","carddav","mattermost","triplefake-com","gocd","gocd-agent","gocd-golang-agent","gocd-java8-agent","gocd-npm10-agent","gocd-packer-agent","letsencrypt","rs-lindau-online-de","jitsi-meet","mattermost-rs-lindau-online-de","mattermost-bellissy"]}`
	var r Message
	err := json.Unmarshal([]byte(respJSON), &r)
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}
