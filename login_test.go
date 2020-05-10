package evmsg

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	//"evalgo.org/eve"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

var (
	EchoTestServerAddress  = "localhost:8888"
	StartTestServerCounter = 0
)

func StartEchoTestServer() *echo.Echo {
	if StartTestServerCounter < 1 {
		StartTestServerCounter = 1
		e := echo.New()
		log.Logger().SetOutput(os.Stdout)
		log.Logger().SetLevel(echoLog.INFO)
		log.Logger().SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
		e.Logger = log.Logger()
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.GET("/ws", func(c echo.Context) error {
			s := websocket.Server{
				Handler: websocket.Handler(func(ws *websocket.Conn) {
					defer ws.Close()
					for {
						msg := NewMessage()
						err := websocket.JSON.Receive(ws, &msg)
						if err != nil {
							c.Logger().Error(err)
						}
						c.Logger().Info(msg)
						err = Auth(msg)
						if err != nil {
							c.Logger().Error(err)
						}
						err = websocket.JSON.Send(ws, &msg)
						if err != nil {
							c.Logger().Error(err)
						}
					}
				}),
				Handshake: func(*websocket.Config, *http.Request) error {
					return nil
				},
			}
			s.ServeHTTP(c.Response(), c.Request())
			return nil
		})
		// Start server
		go func() {
			if err := e.Start(EchoTestServerAddress); err != nil {
				e.Logger.Info("shutting down the server")
			}
		}()
		return e
	}
	return nil
}

func StopEchoTestServer(e *echo.Echo) {
	if e != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
		StartTestServerCounter = 0
		cancel()
		//time.Sleep(5 * time.Second)
	}
}

func Test_Unit_Auth(t *testing.T) {
	e := StartEchoTestServer()
	origin := "http://" + EchoTestServerAddress + "/"
	url := "ws://" + EchoTestServerAddress + "/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Error(err)
	}
	msg := NewMessage()
	msg.ID = ID
	msg.Secret = Secret
	msg.Scope = "Login"
	msg.Command = "getToken"
	err = websocket.JSON.Send(ws, &msg)
	if err != nil {
		t.Error(err)
	}
	err = websocket.JSON.Receive(ws, &msg)
	if err != nil {
		t.Error(err)
	}
	if msg.Debug.Error != "" {
		t.Error(errors.New(msg.Debug.Error))
	}
	t.Log(msg)
	StopEchoTestServer(e)
}

func Test_Unit_AuthErrorCredentials(t *testing.T) {
	e := StartEchoTestServer()
	origin := "http://" + EchoTestServerAddress + "/"
	url := "ws://" + EchoTestServerAddress + "/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Error(err)
	}
	msg := NewMessage()
	msg.ID = "myID"
	msg.Secret = "test"
	msg.Scope = "Login"
	msg.Command = "getToken"
	err = websocket.JSON.Send(ws, &msg)
	if err != nil {
		t.Error(err)
	}
	err = websocket.JSON.Receive(ws, &msg)
	if err != nil {
		t.Error(err)
	}
	if msg.Debug.Error != "the given ID or Secret is not valid!" {
		t.Error(errors.New(msg.Debug.Error))
	}
	t.Log(msg)
	StopEchoTestServer(e)
}

func Test_Unit_AuthErrorToken(t *testing.T) {
	e := StartEchoTestServer()
	origin := "http://" + EchoTestServerAddress + "/"
	url := "ws://" + EchoTestServerAddress + "/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Error(err)
	}
	msg := NewMessage()
	msg.ID = ID
	msg.Secret = Secret
	msg.Scope = "Login"
	msg.Command = "getToken"
	err = websocket.JSON.Send(ws, &msg)
	if err != nil {
		t.Error(err)
	}
	err = websocket.JSON.Receive(ws, &msg)
	if err != nil {
		t.Error(err)
	}
	if msg.Debug.Error != "" {
		t.Error(errors.New(msg.Debug.Error))
	}
	//data := msg.Data.([]interface{})[0]
	//msg.Token = data.(interface{}).(map[string]interface{})["token"].(string)
	msg.Token = "thisIsAnErrorToken"
	msg.Debug = &Debug{}
	msg.State = "Created"
	msg.Scope = "Test"
	msg.Command = "getTokenError"
	msg.Data = nil
	err = websocket.JSON.Send(ws, &msg)
	if err != nil {
		t.Error(err)
	}
	err = websocket.JSON.Receive(ws, &msg)
	if err != nil {
		t.Error(err)
	}
	if msg.Debug.Error != "the given Token is not valid!" {
		t.Error(errors.New("expected not valid token went wrong"))
	}
	StopEchoTestServer(e)
}

func Test_Unit_CheckRequiredKeys(t *testing.T) {
	keys := []string{"test1", "test2"}
	msg := NewMessage()
	msg.Data = []interface{}{
		map[string]interface{}{"test1": "value1", "test2": "value2"},
	}
	err := CheckRequiredKeys(msg, keys)
	if err != nil {
		t.Error(err)
	}
	if msg.Value("test1") != "value1" {
		t.Error("could not get value1 from test1 key")
	}
	t.Log(msg.Values())
}
