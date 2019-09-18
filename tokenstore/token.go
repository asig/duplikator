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

package tokenstore

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"strconv"
	"time"

	"github.com/mrjones/oauth"
)

type Token struct {
	oauth.AccessToken
}

func TokenFromString(s string) (*Token, error) {
	by, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	t := Token{}
	err = d.Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (t *Token) String() (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(*t)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func (t *Token) IsValid() bool {
	if t == nil {
		return false
	}
	now := int(time.Now().UnixNano() / int64(time.Millisecond))
	expiry, _ := strconv.Atoi(t.AdditionalData["edam_expires"])
	return now < expiry
}
