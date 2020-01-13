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
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/asig/duplikator/edam"
	"github.com/asig/duplikator/repository"
	"github.com/asig/duplikator/tokenstore"

	"golang.org/x/net/html"
)

var (
	ns     edam.NoteStore
	client *evernoteClient

	destDirFlag = flag.String("dest_dir", "/tmp/evernote-backup", "Destination directory");
	sandboxFlag = flag.Bool("sandbox", false, "Use sandbox server if true")

	obfuscateFlag = flag.Bool("obfuscate", false, "")
)

type noteWithResources struct {
	note      *edam.Note
	resources map[string]*edam.Resource
}

type command func() error;

func (note noteWithResources) dump() {
	log.Printf("Note: Title = %s", *note.note.Title)
	log.Printf("      ContentLength = %d", *note.note.ContentLength)
	for key, r := range note.resources {
		log.Printf("      Resource[%s].GUID = %s", key, *r.GUID)
		log.Printf("      Resource[%s].Mime = %s", key, *r.Mime)
		if r.Attributes == nil {
			log.Printf("      Resource[%s].Attributes IS NULL!")
		} else {
			filename := "<nil>"
			if r.Attributes.FileName != nil {
				filename = *r.Attributes.FileName
			}
			log.Printf("      Resource[%s].Attributes.FileName = %s", key, filename)
		}
		if r.Data == nil {
			log.Printf("      Resource[%s].Data IS NULL!")
		} else {
			log.Printf("      Resource[%s].Data.Size = %d", key, *r.Data.Size)
			log.Printf("      Resource[%s].Data.BodyHash = %x", key, r.Data.BodyHash)
		}
	}
}

func getCommand() (command, error) {
	args := flag.Args()
	if len(args) == 0 {
		return sync, nil;
	}
	switch (args[0]) {
	case "list":
		if len(args) > 1 {
			guids := args[1:]
			return func() error {
				return list(guids)
			}, nil
		} else {
			return listAll, nil
		}
	case "sync":
		if len(args) > 1 {
			return nil, errors.New("'sync' does not accept parameters")
		}
		return sync, nil
	case "duplicate":
		if len(args) > 1 {
			guids := args[1:]
			return func() error {
				return duplicate(guids)
			}, nil
		} else {
			return duplicateAll, nil
		}
	}
	return nil, fmt.Errorf("%q is not a valid command.", strings.Join(args, " "))
}

func main() {
	flag.Parse()

	if *obfuscateFlag {
		obfuscateCreds(flag.Arg(0), flag.Arg(1))
	}

	tokenStore, err := tokenstore.Init()
	if err != nil {
		log.Fatal(err)
	}

	env := PRODUCTION
	if *sandboxFlag {
		env = SANDBOX
	}
	client = newEvernoteClient(env);
	if err := client.authenticate(tokenStore); err != nil {
		log.Fatal(err);
	}
	ns, err = client.getNoteStore(context.Background())
	if err != nil {
		log.Fatal(err);
	}

	command, err := getCommand()
	if err != nil {
		log.Fatal(err)
	}
	err = command()
	if err != nil {
		log.Fatal(err)
	}
}

func getAllNoteMetadata() ([]*edam.NoteMetadata, error) {
	res := []*edam.NoteMetadata{}

	start := int32(0);
	filter := edam.NewNoteFilter()
	order := int32(edam.NoteSortOrder_CREATED)
	filter.Order = &order;
	resultSpec := &edam.NotesMetadataResultSpec{
		IncludeTitle: boolVal(true),
		IncludeUpdateSequenceNum: boolVal(true),
	}
	for {
		list, err := ns.FindNotesMetadata(context.Background(), client.authToken, filter, start, 100, resultSpec)
		if err != nil {
			return res, err
		}
		if len(list.Notes) == 0 {
			break
		}
		for _, n := range list.Notes {
			res = append(res, n)
		}
		start += int32(len(list.Notes))
	}
	return res, nil
}

func getAllGUIDs() ([]string, error) {
	res := []string{}
	notes, err := getAllNoteMetadata()
	if err != nil {
		return res, err
	}
	for _, n := range notes {
		res = append(res, string(n.GUID))
	}
	return res, nil
}

func sync() error {
	repo, err := repository.Load(*destDirFlag)
	if err != nil {
		return err
	}
	syncedRepo := repository.New(*destDirFlag)
	metadatas, err := getAllNoteMetadata()
	if err != nil {
		return err
	}
	for _, md := range metadatas {
		var e *repository.Entry
		var ok bool
		guid := string(md.GUID)
		if e, ok = repo.Get(guid); ok {
			if e.UpdateSequenceNum <= int64(*md.UpdateSequenceNum) {
				log.Printf("Note %q (%s) is up to date", *md.Title, guid)
				syncedRepo.Add(e)
				continue
			}
			// Existing Note, but needs downloading
			e = syncedRepo.Add(e)
		} else {
	        // New Note, needs downloading
			e = syncedRepo.GetOrAdd(guid)
		}
		log.Printf("Downloading Note %q (%s)", *md.Title, string(md.GUID))
		n, err := fetchNote(guid, context.Background());
		if err != nil {
			return err
		}
		handle(n)
		e.UpdateSequenceNum = int64(*md.UpdateSequenceNum)
		e.Title = *md.Title
	}

	// Delete old files that are not in the new repo
	for _, guid := range repo.GUIDs() {
	    if _, ok := syncedRepo.Get(guid); !ok {
	    	e, _ := repo.Get(guid);
			log.Printf("Deleting %q (%s)", e.Title, e.GUID);
			os.RemoveAll(baseName(e.Title, e.GUID));
        }
	}
	err = syncedRepo.Save()
	return err
}

func listAll() error {
	guids, err := getAllGUIDs()
	if err != nil {
		return err
	}
	return list(guids);
}

func list(guids []string) error {
	nrs := &edam.NoteResultSpec{
		IncludeContent:                boolVal(false),
		IncludeResourcesData:          boolVal(false),
		IncludeResourcesAlternateData: boolVal(false),
	}

	for _, guid := range guids {
		note, err := ns.GetNoteWithResultSpec(context.Background(), client.authToken, edam.GUID(guid), nrs)
		if err != nil {
			log.Printf("Can't get note %s: %s", guid, err);
			continue
		}
		fmt.Printf("%s: %s\n", guid, *note.Title)
	}
	return nil
}

func duplicateAll() error {

	guids, err := getAllGUIDs()
	if err != nil {
		return err
	}
	return duplicate(guids);
}

func duplicate(guids []string) error {
	for _, guid := range guids {
		log.Printf("Downloading Note %s", guid)
		if note, err := fetchNote(guid, context.Background()); err == nil {
			err = handle(note)
			if err != nil {
				log.Printf("Error while handling %s: %s", *note.note.Title, err)
			}
		} else {
			log.Printf("Can't download note %s: %s", guid, err)
		}
	}
	return nil
}

func handle(note noteWithResources) error {
	return note.save()
}

func baseName(title, guid string) string {
	filename := makeFilename(title)
	return filepath.Join(*destDirFlag, filename+"-"+guid)
}

func (note noteWithResources) baseName() string {
	return baseName(*note.note.Title, string(*note.note.GUID))
}

func (note noteWithResources) attachmentFileName(hash string, relative bool) string {
	var basename string
	if relative {
		basename = "files"
	} else {
		basename = note.baseName() + "/files"
	}

	attachmentName := ""
	r := note.resources[hash]
	if r.Attributes.FileName != nil {
		attachmentName = *r.Attributes.FileName
	} else {
		// No filename given, lets create one
		attachmentName = string(*r.GUID)
		exts, _ := mime.ExtensionsByType(*r.Mime)
		if len(exts) == 0 {
			log.Printf("Can't find file suffix for mime type %s", *r.Mime)
		} else {
			attachmentName = attachmentName + exts[0]
		}
	}
	return filepath.Join(basename, attachmentName)
}

func (note noteWithResources) save() error {
	os.RemoveAll(note.baseName());

	// Save html

	filename := note.baseName() + "/" + makeFilename(note.note.GetTitle()) + ".html"
	err := os.MkdirAll(path.Dir(filename), 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	note.convertToHtml(f)
	err = f.Close()
	if err != nil {
		return err
	}

	// Save attachments
	for hash, res := range note.resources {
		filename = note.attachmentFileName(hash, false)
		err = os.MkdirAll(path.Dir(filename), 0755)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, res.Data.Body, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (note noteWithResources) convertToHtml(w io.Writer) {
	w.Write([]byte(`<!doctype html>
<html>
<head>
</head>
`))
	z := html.NewTokenizer(strings.NewReader(*note.note.Content))
	for {
		if z.Next() == html.ErrorToken {
			if z.Err() == io.EOF {
				// Returning io.EOF indicates success.
				break;
			}
			log.Fatal(z.Err())
		}
		tok := z.Token();
		if tok.Type == html.DoctypeToken || tok.Type == html.CommentToken {
			continue;
		}
		switch (tok.Data) {
		case "en-note":
			if tok.Type == html.StartTagToken {
				w.Write([]byte("<body>"))
			} else if tok.Type == html.EndTagToken {
				w.Write([]byte("</body>"))
			} else {
				log.Fatalf("Can't happen! token = %s", tok)
			}
		case "en-media":
			if tok.Type == html.EndTagToken {
				break;
			}
			// Start end SelfClosing, generate the complete link
			t, _ := findAttribute(tok, "type");
			h, _ := findAttribute(tok, "hash");
			filename := note.attachmentFileName(h, true)
			if isImage(t) {
				height := ""
				width := ""
				if w, ok := findAttribute(tok, "width"); ok {
					width = fmt.Sprintf("width=\"%s\"", w)
				}
				if h, ok := findAttribute(tok, "height"); ok {
					height = fmt.Sprintf("height=\"%s\"", h)
				}
				w.Write([]byte(fmt.Sprintf("<img src=\"%s\" %s %s>", filename, width, height)))
			} else {
				displayName := note.resources[h].Attributes.FileName
				if displayName == nil {
					displayName = &filename
				}
				w.Write([]byte(fmt.Sprintf("<a href=\"%s\">%s</a>", filename, *displayName)))
			}
		default:
			w.Write([]byte(tok.String()))
		}
	}
	w.Write([]byte("</html>"))
}

func isImage(mimetype string) bool {
	return strings.HasPrefix(mimetype, "image/")
}

func findAttribute(token html.Token, name string) (vak string, found bool) {
	for _, attr := range token.Attr {
		if attr.Key == name {
			return attr.Val, true
		}
	}
	return "", false
}

func makeFilename(raw string) string {
	// Bad as defined by wikipedia: https://en.wikipedia.org/wiki/Filename#Reserved_characters_and_words
	// Also have to escape the backslash
	res := ""
	for _, ch := range raw {
		switch (ch) {
		case '/', '\\', '?', '%', '*', ':', '|', '"', '<', '>', '.':
			res += "_"
		default:
			res += string(ch)
		}
	}
	return res
}

func fetchNote(guid string, ctx context.Context) (noteWithResources, error) {
	var err error

	note := noteWithResources{}
	nrs := &edam.NoteResultSpec{
		IncludeContent:                boolVal(true),
		IncludeResourcesData:          boolVal(true),
		IncludeResourcesAlternateData: boolVal(true),
	}
	note.note, err = ns.GetNoteWithResultSpec(ctx, client.authToken, edam.GUID(guid), nrs)
	if err != nil {
		return note, err
	}

	if len(note.note.Resources) > 0 {
		note.resources = make(map[string]*edam.Resource)
		for _, res := range note.note.Resources {
			if r, err := ns.GetResource(ctx, client.authToken, *res.GUID, /* withData= */ true, /* withRecognition= */ true, /* withAttributes= */ true, /* withAlternateData */ true); err == nil {
				note.resources[hex.EncodeToString(r.Data.BodyHash)] = res
			} else {
				return note, err
			}
		}
	}
	return note, nil
}

func boolVal(val bool) *bool {
	b := val
	return &b
}
