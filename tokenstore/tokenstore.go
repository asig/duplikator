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
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type Store struct {
    filename string
    Token *Token
}

var (
	accessTokenFlag    = flag.String("access_token", "", "Access token to use")
	tokenStoreFileFlag = flag.String("token_store", "", "File to store tokens in")
)

func Init() (*Store, error) {
	store := Store{filename: tokenStoreFileName()}

	b, err := ioutil.ReadFile(store.filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	token, err := TokenFromString(string(b))
	if err == nil {
		store.Token = token
	}

	if *accessTokenFlag != "" {
		token, err = TokenFromString(*accessTokenFlag)
		if err != nil {
			return nil, err
		}
		store.Token = token
	}
	return &store, nil
}

func (store *Store) Save() error {
	s := ""
	if store.Token != nil {
		s, _ = store.Token.String()
	}
	err := os.MkdirAll(path.Dir(store.filename), 0700)
	if err != nil {
		return err;
	}
	return ioutil.WriteFile(store.filename, []byte(s), 0600)
}

func tokenStoreFileName() string {
    name := *tokenStoreFileFlag
    if name != "" {
        return name
    }
    return buildTokenStoreName()
}

func buildTokenStoreName() string {
	name := *tokenStoreFileFlag
	if name == "" {
		name = filepath.Join(getDataDir(), "token_store")
	}
	return name
}

func getDataDir() string {
	var dir string
	switch runtime.GOOS {
	case "darwin":
		dir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	case "windows":
		dir = os.Getenv("APPDATA")
	default:
		dir = os.Getenv("XDG_DATA_HOME");
		if dir == "" {
			dir = filepath.Join(os.Getenv("HOME"), ".local", "share")
		}
	}
	return  filepath.Join(dir, "duplikator")
}
