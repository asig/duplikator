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

	"github.com/asig/duplikator/edam"
)

type throttlingNoteStore struct {
	ns edam.NoteStore
}

func (t throttlingNoteStore) GetSyncState(ctx context.Context, authenticationToken string) (r *edam.SyncState, err error) {
	for {
		res, err := t.ns.GetSyncState(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetFilteredSyncChunk(ctx context.Context, authenticationToken string, afterUSN int32, maxEntries int32, filter *edam.SyncChunkFilter) (r *edam.SyncChunk, err error) {
	for {
		res, err := t.ns.GetFilteredSyncChunk(ctx, authenticationToken, afterUSN, maxEntries, filter)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetLinkedNotebookSyncState(ctx context.Context, authenticationToken string, linkedNotebook *edam.LinkedNotebook) (r *edam.SyncState, err error) {
	for {
		res, err := t.ns.GetLinkedNotebookSyncState(ctx, authenticationToken, linkedNotebook)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetLinkedNotebookSyncChunk(ctx context.Context, authenticationToken string, linkedNotebook *edam.LinkedNotebook, afterUSN int32, maxEntries int32, fullSyncOnly bool) (r *edam.SyncChunk, err error) {
	for {
		for {
			res, err := t.ns.GetLinkedNotebookSyncChunk(ctx, authenticationToken, linkedNotebook, afterUSN, maxEntries, fullSyncOnly)
			if maybeThrottle(err) {
				continue
			}
			return res, err
		}
	}
}

func (t throttlingNoteStore) ListNotebooks(ctx context.Context, authenticationToken string) (r []*edam.Notebook, err error) {
	for {
		res, err := t.ns.ListNotebooks(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) ListAccessibleBusinessNotebooks(ctx context.Context, authenticationToken string) (r []*edam.Notebook, err error) {
	for {
		res, err := t.ns.ListAccessibleBusinessNotebooks(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetNotebook(ctx context.Context, authenticationToken string, guid edam.GUID) (r *edam.Notebook, err error) {
	for {
		res, err := t.ns.GetNotebook(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetDefaultNotebook(ctx context.Context, authenticationToken string) (r *edam.Notebook, err error) {
	for {

		res, err := t.ns.GetDefaultNotebook(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) CreateNotebook(ctx context.Context, authenticationToken string, notebook *edam.Notebook) (r *edam.Notebook, err error) {
	for {
		res, err := t.ns.CreateNotebook(ctx, authenticationToken, notebook)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}

}

func (t throttlingNoteStore) UpdateNotebook(ctx context.Context, authenticationToken string, notebook *edam.Notebook) (r int32, err error) {
	for {
		res, err := t.ns.UpdateNotebook(ctx, authenticationToken, notebook)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}

}

func (t throttlingNoteStore) ExpungeNotebook(ctx context.Context, authenticationToken string, guid edam.GUID) (r int32, err error) {
	for {

		res, err := t.ns.ExpungeNotebook(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err

	}
}

func (t throttlingNoteStore) ListTags(ctx context.Context, authenticationToken string) (r []*edam.Tag, err error) {
	for {

		res, err := t.ns.ListTags(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err

	}
}

func (t throttlingNoteStore) ListTagsByNotebook(ctx context.Context, authenticationToken string, notebookGuid edam.GUID) (r []*edam.Tag, err error) {
	for {
 		res, err := t.ns.ListTagsByNotebook(ctx, authenticationToken, notebookGuid)
		if maybeThrottle(err) {
			continue
		}
		return res, err

	}
}

func (t throttlingNoteStore) GetTag(ctx context.Context, authenticationToken string, guid edam.GUID) (r *edam.Tag, err error) {
	for {
 		res, err := t.ns.GetTag(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err

	}
}

func (t throttlingNoteStore) CreateTag(ctx context.Context, authenticationToken string, tag *edam.Tag) (r *edam.Tag, err error) {
	for {
		res, err := t.ns.CreateTag(ctx, authenticationToken, tag)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) UpdateTag(ctx context.Context, authenticationToken string, tag *edam.Tag) (r int32, err error) {
	for {
		res, err := t.ns.UpdateTag(ctx, authenticationToken, tag)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) UntagAll(ctx context.Context, authenticationToken string, guid edam.GUID) (err error) {
	for {
		err = t.ns.UntagAll(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return err
	}
}

func (t throttlingNoteStore) ExpungeTag(ctx context.Context, authenticationToken string, guid edam.GUID) (r int32, err error) {
	for {
		res, err := t.ns.ExpungeTag(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) ListSearches(ctx context.Context, authenticationToken string) (r []*edam.SavedSearch, err error) {
	for {
		res, err := t.ns.ListSearches(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetSearch(ctx context.Context, authenticationToken string, guid edam.GUID) (r *edam.SavedSearch, err error) {
	for {
		res, err := t.ns.GetSearch(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) CreateSearch(ctx context.Context, authenticationToken string, search *edam.SavedSearch) (r *edam.SavedSearch, err error) {
	for {
		res, err := t.ns.CreateSearch(ctx, authenticationToken, search)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) UpdateSearch(ctx context.Context, authenticationToken string, search *edam.SavedSearch) (r int32, err error) {
	for {
		res, err := t.ns.UpdateSearch(ctx, authenticationToken, search)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) ExpungeSearch(ctx context.Context, authenticationToken string, guid edam.GUID) (r int32, err error) {
	for {
		res, err := t.ns.ExpungeSearch(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) FindNoteOffset(ctx context.Context, authenticationToken string, filter *edam.NoteFilter, guid edam.GUID) (r int32, err error) {
	for {
		res, err := t.ns.FindNoteOffset(ctx, authenticationToken, filter, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) FindNotesMetadata(ctx context.Context, authenticationToken string, filter *edam.NoteFilter, offset int32, maxNotes int32, resultSpec *edam.NotesMetadataResultSpec) (r *edam.NotesMetadataList, err error) {
	for {
		res, err := t.ns.FindNotesMetadata(ctx, authenticationToken, filter, offset, maxNotes, resultSpec)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) FindNoteCounts(ctx context.Context, authenticationToken string, filter *edam.NoteFilter, withTrash bool) (r *edam.NoteCollectionCounts, err error) {
	for {
		res, err := t.ns.FindNoteCounts(ctx, authenticationToken, filter, withTrash)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetNoteWithResultSpec(ctx context.Context, authenticationToken string, guid edam.GUID, resultSpec *edam.NoteResultSpec) (r *edam.Note, err error) {
	for {
		res, err := t.ns.GetNoteWithResultSpec(ctx, authenticationToken, guid, resultSpec)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetNote(ctx context.Context, authenticationToken string, guid edam.GUID, withContent bool, withResourcesData bool, withResourcesRecognition bool, withResourcesAlternateData bool) (r *edam.Note, err error) {
	for {
		res, err := t.ns.GetNote(ctx, authenticationToken, guid, withContent, withResourcesData, withResourcesRecognition, withResourcesAlternateData)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetNoteApplicationData(ctx context.Context, authenticationToken string, guid edam.GUID) (r *edam.LazyMap, err error) {
	for {
		res, err := t.ns.GetNoteApplicationData(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetNoteApplicationDataEntry(ctx context.Context, authenticationToken string, guid edam.GUID, key string) (r string, err error) {
	for {
		res, err := t.ns.GetNoteApplicationDataEntry(ctx, authenticationToken, guid, key)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) SetNoteApplicationDataEntry(ctx context.Context, authenticationToken string, guid edam.GUID, key string, value string) (r int32, err error) {
	for {
		res, err := t.ns.SetNoteApplicationDataEntry(ctx, authenticationToken, guid, key, value)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) UnsetNoteApplicationDataEntry(ctx context.Context, authenticationToken string, guid edam.GUID, key string) (r int32, err error) {
	for {
		res, err := t.ns.UnsetNoteApplicationDataEntry(ctx, authenticationToken, guid, key)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetNoteContent(ctx context.Context, authenticationToken string, guid edam.GUID) (r string, err error) {
	for {
		res, err := t.ns.GetNoteContent(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetNoteSearchText(ctx context.Context, authenticationToken string, guid edam.GUID, noteOnly bool, tokenizeForIndexing bool) (r string, err error) {
	for {
		res, err := t.ns.GetNoteSearchText(ctx, authenticationToken, guid, noteOnly, tokenizeForIndexing)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetResourceSearchText(ctx context.Context, authenticationToken string, guid edam.GUID) (r string, err error) {
	for {
		res, err := t.ns.GetResourceSearchText(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetNoteTagNames(ctx context.Context, authenticationToken string, guid edam.GUID) (r []string, err error) {
	for {
		res, err := t.ns.GetNoteTagNames(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) CreateNote(ctx context.Context, authenticationToken string, note *edam.Note) (r *edam.Note, err error) {
	for {
		res, err := t.ns.CreateNote(ctx, authenticationToken, note)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) UpdateNote(ctx context.Context, authenticationToken string, note *edam.Note) (r *edam.Note, err error) {
	for {
		res, err := t.ns.UpdateNote(ctx, authenticationToken, note)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) DeleteNote(ctx context.Context, authenticationToken string, guid edam.GUID) (r int32, err error) {
	for {
		res, err := t.ns.DeleteNote(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) ExpungeNote(ctx context.Context, authenticationToken string, guid edam.GUID) (r int32, err error) {
	for {
		res, err := t.ns.ExpungeNote(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) CopyNote(ctx context.Context, authenticationToken string, noteGuid edam.GUID, toNotebookGuid edam.GUID) (r *edam.Note, err error) {
	for {
		res, err := t.ns.CopyNote(ctx, authenticationToken, noteGuid, toNotebookGuid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) ListNoteVersions(ctx context.Context, authenticationToken string, noteGuid edam.GUID) (r []*edam.NoteVersionId, err error) {
	for {
		res, err := t.ns.ListNoteVersions(ctx, authenticationToken, noteGuid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetNoteVersion(ctx context.Context, authenticationToken string, noteGuid edam.GUID, updateSequenceNum int32, withResourcesData bool, withResourcesRecognition bool, withResourcesAlternateData bool) (r *edam.Note, err error) {
	for {
		res, err := t.ns.GetNoteVersion(ctx, authenticationToken, noteGuid, updateSequenceNum, withResourcesData, withResourcesRecognition, withResourcesAlternateData)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetResource(ctx context.Context, authenticationToken string, guid edam.GUID, withData bool, withRecognition bool, withAttributes bool, withAlternateData bool) (r *edam.Resource, err error) {
	for {
		res, err := t.ns.GetResource(ctx, authenticationToken, guid, withData, withRecognition, withAttributes, withAlternateData)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetResourceApplicationData(ctx context.Context, authenticationToken string, guid edam.GUID) (r *edam.LazyMap, err error) {
	for {
		res, err := t.ns.GetResourceApplicationData(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetResourceApplicationDataEntry(ctx context.Context, authenticationToken string, guid edam.GUID, key string) (r string, err error) {
	for {
		res, err := t.ns.GetResourceApplicationDataEntry(ctx, authenticationToken, guid, key)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) SetResourceApplicationDataEntry(ctx context.Context, authenticationToken string, guid edam.GUID, key string, value string) (r int32, err error) {
	for {
		res, err := t.ns.SetResourceApplicationDataEntry(ctx, authenticationToken, guid, key, value)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) UnsetResourceApplicationDataEntry(ctx context.Context, authenticationToken string, guid edam.GUID, key string) (r int32, err error) {
	for {
		res, err := t.ns.UnsetResourceApplicationDataEntry(ctx, authenticationToken, guid, key)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) UpdateResource(ctx context.Context, authenticationToken string, resource *edam.Resource) (r int32, err error) {
	for {
		res, err := t.ns.UpdateResource(ctx, authenticationToken, resource)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetResourceData(ctx context.Context, authenticationToken string, guid edam.GUID) (r []byte, err error) {
	for {
		res, err := t.ns.GetResourceData(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetResourceByHash(ctx context.Context, authenticationToken string, noteGuid edam.GUID, contentHash []byte, withData bool, withRecognition bool, withAlternateData bool) (r *edam.Resource, err error) {
	for {
		res, err := t.ns.GetResourceByHash(ctx, authenticationToken, noteGuid, contentHash, withData, withRecognition, withAlternateData)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetResourceRecognition(ctx context.Context, authenticationToken string, guid edam.GUID) (r []byte, err error) {
	for {
		res, err := t.ns.GetResourceRecognition(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetResourceAlternateData(ctx context.Context, authenticationToken string, guid edam.GUID) (r []byte, err error) {
	for {
		res, err := t.ns.GetResourceAlternateData(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetResourceAttributes(ctx context.Context, authenticationToken string, guid edam.GUID) (r *edam.ResourceAttributes, err error) {
	for {
		res, err := t.ns.GetResourceAttributes(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetPublicNotebook(ctx context.Context, userId edam.UserID, publicUri string) (r *edam.Notebook, err error) {
	for {
		res, err := t.ns.GetPublicNotebook(ctx, userId, publicUri)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) ShareNotebook(ctx context.Context, authenticationToken string, sharedNotebook *edam.SharedNotebook, message string) (r *edam.SharedNotebook, err error) {
	for {
		res, err := t.ns.ShareNotebook(ctx, authenticationToken, sharedNotebook, message)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) CreateOrUpdateNotebookShares(ctx context.Context, authenticationToken string, shareTemplate *edam.NotebookShareTemplate) (r *edam.CreateOrUpdateNotebookSharesResult_, err error) {
	for {
		res, err := t.ns.CreateOrUpdateNotebookShares(ctx, authenticationToken, shareTemplate)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) UpdateSharedNotebook(ctx context.Context, authenticationToken string, sharedNotebook *edam.SharedNotebook) (r int32, err error) {
	for {
		res, err := t.ns.UpdateSharedNotebook(ctx, authenticationToken, sharedNotebook)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) SetNotebookRecipientSettings(ctx context.Context, authenticationToken string, notebookGuid string, recipientSettings *edam.NotebookRecipientSettings) (r *edam.Notebook, err error) {
	for {
		res, err := t.ns.SetNotebookRecipientSettings(ctx, authenticationToken, notebookGuid, recipientSettings)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) ListSharedNotebooks(ctx context.Context, authenticationToken string) (r []*edam.SharedNotebook, err error) {
	for {
		res, err := t.ns.ListSharedNotebooks(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) CreateLinkedNotebook(ctx context.Context, authenticationToken string, linkedNotebook *edam.LinkedNotebook) (r *edam.LinkedNotebook, err error) {
	for {
		res, err := t.ns.CreateLinkedNotebook(ctx, authenticationToken, linkedNotebook)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) UpdateLinkedNotebook(ctx context.Context, authenticationToken string, linkedNotebook *edam.LinkedNotebook) (r int32, err error) {
	for {
		res, err := t.ns.UpdateLinkedNotebook(ctx, authenticationToken, linkedNotebook)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) ListLinkedNotebooks(ctx context.Context, authenticationToken string) (r []*edam.LinkedNotebook, err error) {
	for {
		res, err := t.ns.ListLinkedNotebooks(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) ExpungeLinkedNotebook(ctx context.Context, authenticationToken string, guid edam.GUID) (r int32, err error) {
	for {
		res, err := t.ns.ExpungeLinkedNotebook(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) AuthenticateToSharedNotebook(ctx context.Context, shareKeyOrGlobalId string, authenticationToken string) (r *edam.AuthenticationResult_, err error) {
	for {
		res, err := t.ns.AuthenticateToSharedNotebook(ctx, shareKeyOrGlobalId, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetSharedNotebookByAuth(ctx context.Context, authenticationToken string) (r *edam.SharedNotebook, err error) {
	for {
		res, err := t.ns.GetSharedNotebookByAuth(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) EmailNote(ctx context.Context, authenticationToken string, parameters *edam.NoteEmailParameters) (err error) {
	for {
		err = t.ns.EmailNote(ctx, authenticationToken, parameters)
		if maybeThrottle(err) {
			continue
		}
		return err
	}
}

func (t throttlingNoteStore) ShareNote(ctx context.Context, authenticationToken string, guid edam.GUID) (r string, err error) {
	for {
		res, err := t.ns.ShareNote(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) StopSharingNote(ctx context.Context, authenticationToken string, guid edam.GUID) (err error) {
	for {
		err = t.ns.StopSharingNote(ctx, authenticationToken, guid)
		if maybeThrottle(err) {
			continue
		}
		return err
	}
}

func (t throttlingNoteStore) AuthenticateToSharedNote(ctx context.Context, guid string, noteKey string, authenticationToken string) (r *edam.AuthenticationResult_, err error) {
	for {
		res, err := t.ns.AuthenticateToSharedNote(ctx, guid, noteKey, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) FindRelated(ctx context.Context, authenticationToken string, query *edam.RelatedQuery, resultSpec *edam.RelatedResultSpec) (r *edam.RelatedResult_, err error) {
	for {
		res, err := t.ns.FindRelated(ctx, authenticationToken, query, resultSpec)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) UpdateNoteIfUsnMatches(ctx context.Context, authenticationToken string, note *edam.Note) (r *edam.UpdateNoteIfUsnMatchesResult_, err error) {
	for {
		res, err := t.ns.UpdateNoteIfUsnMatches(ctx, authenticationToken, note)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) ManageNotebookShares(ctx context.Context, authenticationToken string, parameters *edam.ManageNotebookSharesParameters) (r *edam.ManageNotebookSharesResult_, err error) {
	for {
		res, err := t.ns.ManageNotebookShares(ctx, authenticationToken, parameters)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingNoteStore) GetNotebookShares(ctx context.Context, authenticationToken string, notebookGuid string) (r *edam.ShareRelationships, err error) {
	for {
		res, err := t.ns.GetNotebookShares(ctx, authenticationToken, notebookGuid)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}
