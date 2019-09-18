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

package repository

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Entry struct {
	GUID string `json:"guid"`
	UpdateSequenceNum int64 `json:"updated"`
	Title string `json:"title"`
}

type Repo struct {
    filename string
	entries []Entry
};

func New(baseDir string) *Repo {
	name := path.Join(baseDir, "repository.json")
	res := Repo{filename: name, entries: []Entry{} };
	return &res
}

func Load(baseDir string) (*Repo, error) {
	res := New(baseDir);
	log.Printf("Loading repository from %s", res.filename);
	jsonFile, err := os.Open(res.filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Repository %s not found", res.filename);
			return res, nil
		}
		return res, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &res.entries)
	return res, err
}

func (r *Repo) Save() error {
	log.Printf("Writing repository to %s", r.filename);
	file, _ := json.MarshalIndent(r.entries, "", " ")
	return ioutil.WriteFile(r.filename, file, 0644)
}

func (r *Repo) Get(guid string) (entry *Entry, ok bool) {
	for _, e := range r.entries {
		if e.GUID == guid {
			return &e, true
		}
	}
	return nil, false
}

func (r *Repo) GetOrAdd(guid string) *Entry {
	if e, ok := r.Get(guid); ok {
		return e
	}
	r.entries = append(r.entries, Entry{GUID: guid})
	return &r.entries[len(r.entries) - 1]
}

func (r *Repo) Add(entry *Entry) *Entry {
    e := r.GetOrAdd(entry.GUID)
    e.GUID = entry.GUID
    e.UpdateSequenceNum = entry.UpdateSequenceNum
    e.Title = entry.Title
    return e
}

func (r *Repo) GUIDs() []string {
    res := []string{}
	for _, e := range r.entries {
	    res = append(res, e.GUID)
	}
	return res
}
