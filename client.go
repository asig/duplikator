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
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/asig/duplikator/tokenstore"
	"io"
	"log"
	"os"
	"strings"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/asig/duplikator/edam"
	"github.com/mrjones/oauth"
)

type environmentType int

const (
	SANDBOX environmentType = iota
	PRODUCTION
	YINXIANG
)

var (
	_k                          = []byte{0xcb, 0x2c, 0x75, 0x1d, 0x3b, 0xdd, 0x2d, 0xd5, 0x3e, 0x0f, 0xfc, 0x62, 0x84, 0x79, 0x93, 0xa0, 0x51, 0xa8, 0x31, 0x1b, 0x63, 0x5e, 0xe7, 0x67, 0xc0, 0x27, 0xad, 0x5e, 0x4c, 0x10, 0x62, 0x75}
	creds                       = []byte{0x13, 0x30, 0xe1, 0x6f, 0x0e, 0x0d, 0xbe, 0x23, 0x5b, 0xcf, 0x24, 0x17, 0xba, 0xca, 0x68, 0xdd, 0x84, 0x09, 0x4f, 0xcd, 0x3b, 0x4c, 0xd8, 0x24, 0x50, 0xc7, 0x82, 0x6e, 0x46, 0xf9, 0x43, 0xda, 0x3a, 0xc1, 0x07, 0x26, 0xbf, 0xad, 0xb5, 0xd9, 0x74, 0xf3, 0x50, 0x69, 0x7b};
	nonce                       = []byte{0xc5, 0x8f, 0xce, 0xbd, 0xbd, 0x26, 0x11, 0xf7, 0x75, 0xfd, 0xb5, 0x14};
	consumerKey string;
	consumerSecret string
)

func init() {
	c, _ := aes.NewCipher(_k)
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}
	m, err := gcm.Open(nil, nonce, creds, nil)
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(string(m), "|")
	consumerKey, consumerSecret = parts[0], parts[1]
}

func obfuscateCreds(key, secret string) {
	c, _ := aes.NewCipher(_k)
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(key+"|"+secret), nil)
	fmt.Printf("creds = []byte{%s};\n", formatBytes(ciphertext))
	fmt.Printf("nonce = []byte{%s};\n", formatBytes(nonce))

	os.Exit(0)
}

func formatBytes(b []byte) string {
	bytes := []string{}
	for _, b := range b {
		bytes = append(bytes, fmt.Sprintf("0x%02x", b))
	}
	return strings.Join(bytes, ", ")
}

func (e environmentType) host() string {
	switch e {
	case SANDBOX:
		return "sandbox.evernote.com"
	case YINXIANG:
		return "app.yinxiang.com"
	default:
		return "www.evernote.com"
	}
}

type evernoteClient struct {
	host        string
	authToken   string
	oauthClient *oauth.Consumer
	userStore   edam.UserStore
}

func newEvernoteClient(envType environmentType) *evernoteClient {
	host := envType.host()
	client := oauth.NewConsumer(
		consumerKey, consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   fmt.Sprintf("https://%s/oauth", host),
			AuthorizeTokenUrl: fmt.Sprintf("https://%s/OAuth.action", host),
			AccessTokenUrl:    fmt.Sprintf("https://%s/oauth", host),
		},
	)
	return &evernoteClient{
		host:        host,
		oauthClient: client,
	}
}

func (c *evernoteClient) authenticate(ts *tokenstore.Store) error {
	if !ts.Token.IsValid() {
		var err error
		ts.Token, err = login(c.oauthClient);
		if err != nil {
			return err
		}
		ts.Save()
	}
	c.authToken = ts.Token.Token
	return nil
}

func (c *evernoteClient) getUserStore() (edam.UserStore, error) {
	if c.userStore != nil {
		return c.userStore, nil
	}
	evernoteUserStoreServerURL := fmt.Sprintf("https://%s/edam/user", c.host)
	thriftTransport, err := thrift.NewTHttpClient(evernoteUserStoreServerURL)
	thriftClient := thrift.NewTStandardClient(thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(thriftTransport), thrift.NewTBinaryProtocolFactory(true, true).GetProtocol(thriftTransport))
	if err != nil {
		return nil, err
	}
	c.userStore = &throttlingUserStore{edam.NewUserStoreClient(thriftClient)}
	return c.userStore, nil
}

func (c *evernoteClient) getNoteStore(ctx context.Context) (edam.NoteStore, error) {
	us, err := c.getUserStore()
	if err != nil {
		return nil, err
	}
	userUrls, err := us.GetUserUrls(ctx, c.authToken)
	if err != nil {
		return nil, err
	}

	thriftTransport, err := thrift.NewTHttpClient(userUrls.GetNoteStoreUrl())
	thriftClient := thrift.NewTStandardClient(thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(thriftTransport), thrift.NewTBinaryProtocolFactory(true, true).GetProtocol(thriftTransport))
	if err != nil {
		return nil, err
	}

	ns := edam.NewNoteStoreClient(thriftClient)
	return &throttlingNoteStore{ns}, nil
}
