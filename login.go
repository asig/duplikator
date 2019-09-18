/*
 * Copyright (c) 2019 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of Duplikator.
 *
 * Duplikator is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Duplikator is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Duplikator.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"errors"
	"fmt"
	"github.com/asig/duplikator/tokenstore"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/mrjones/oauth"
)


type authHandler struct {
	verifierChannel chan string
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	verifier, _ := r.URL.Query()["oauth_verifier"]
	v := ""
	if verifier != nil {
		v = verifier[0]
	}
	h.verifierChannel <- v

	w.Write([]byte("You can close this window now."))
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func login(oauthClient *oauth.Consumer) (*tokenstore.Token, error) {
	// start web server for callback.
	verifierChannel := make(chan string)
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	log.Printf("Listeing on %s", listener.Addr())
	server := http.Server{
		Addr:    listener.Addr().String(),
		Handler: &authHandler{verifierChannel},
	}
	go server.Serve(listener)

	// 1. Request token
	callbackURL := "http://" + listener.Addr().String() + "/"
	requestToken, url, err := oauthClient.GetRequestTokenAndUrl(callbackURL)
	if err != nil {
		return nil, err
	}

	// 2. Open browser for authorization
	err = openBrowser(url)
	if err != nil {
		return nil, err
	}

	// 3. Retrieve verifier token
	fmt.Printf("Waiting for verifier\n")
	verifier := <-verifierChannel
	if verifier == "" {
		return nil, errors.New("User didn't authorize")
	}

	// 4. Get Access token
	token, err := oauthClient.AuthorizeToken(requestToken, verifier)
	if err != nil {
		return nil, err
	}

	t := &tokenstore.Token{*token}
	return t, nil
}


