/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"evalgo.org/evmsg"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"

	"github.com/spf13/cobra"
)

var (
	VERSION = "v0.0.1"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("address")
		if err != nil {
			return err
		}
		e := echo.New()
		log.Logger().SetOutput(os.Stdout)
		log.Logger().SetLevel(echoLog.INFO)
		log.Logger().SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
		e.Logger = log.Logger()
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Static("/", "webroot")
		e.POST("/"+VERSION+"/upload", func(c echo.Context) error {
			if evmsg.Bearer(c) == nil {
				file, handler, err := c.Request().FormFile("file")
				if err != nil {
					return err
				}
				defer file.Close()
				c.Logger().Info("Uploaded File: %+v\n", handler.Filename)
				c.Logger().Info("File Size: %+v\n", handler.Size)
				c.Logger().Info("MIME Header: %+v\n", handler.Header)
				uploadTgtFolder := "."
				tempFile, err := ioutil.TempFile(uploadTgtFolder, handler.Filename)
				if err != nil {
					return err
				}
				defer tempFile.Close()
				fileBytes, err := ioutil.ReadAll(file)
				if err != nil {
					return err
				}
				tempFile.Write(fileBytes)
				c.Response().WriteHeader(http.StatusOK)
				msg := evmsg.NewMessage()
				msgB, err := json.Marshal(&msg)
				if err != nil {
					return err
				}
				c.Response().Write(msgB)
				return os.Rename(tempFile.Name(), uploadTgtFolder+"/"+handler.Filename)
			}
			c.Response().WriteHeader(http.StatusUnauthorized)
			c.Response().Write([]byte("401 Unauthorized"))
			return errors.New("401 Unauthorized")
		})
		e.GET("/"+VERSION+"/ws", func(c echo.Context) error {
			s := websocket.Server{
				Handler: websocket.Handler(func(ws *websocket.Conn) {
					defer ws.Close()
				WEBSOCKET:
					for {
						var msg evmsg.Message
						err := websocket.JSON.Receive(ws, &msg)
						if err != nil {
							c.Logger().Error(err)
							if err == io.EOF {
								c.Logger().Info("websocket client closed connection!")
								return
							}
							msg.Debug.Error = err.Error()
						}
						err = evmsg.Auth(&msg)
						if err != nil {
							c.Logger().Error(err)
							err = websocket.JSON.Send(ws, &msg)
							if err != nil {
								c.Logger().Error(err)
							}
							continue WEBSOCKET
						}
						switch msg.Scope {
						case "List":
							switch msg.Command {
							case "getAll":
								msg.State = "Response"
								msg.Data = []map[string]string{
									map[string]string{"test0": "value01", "test1": "value10"},
									map[string]string{"test0": "value01", "test1": "value11"},
									map[string]string{"test0": "value02", "test1": "value12"},
								}
							}
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
		return e.Start(address)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().String("address", "localhost:8888", "address where to start the service at")
}
