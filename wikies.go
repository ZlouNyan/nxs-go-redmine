package redmine

import (
	"net/http"
	"net/url"
	"strconv"
)

/* Get */

// WikiMultiObject struct used for wikies all get operations
type WikiMultiObject struct {
	Title     string            `json:"title"`
	Parent    *WikiParentObject `json:"parent"`
	Version   int               `json:"version"`
	CreatedOn string            `json:"created_on"`
	UpdatedOn string            `json:"updated_on"`
}

// WikiObject struct used for wiki get operations
type WikiObject struct {
	Title       string              `json:"title"`
	Parent      *WikiParentObject   `json:"parent"`
	Text        string              `json:"text"`
	Version     int                 `json:"version"`
	Author      IDName              `json:"author"`
	Comments    string              `json:"comments"`
	CreatedOn   string              `json:"created_on"`
	UpdatedOn   string              `json:"updated_on"`
	Attachments *[]AttachmentObject `json:"attachments"`
}

// WikiParentObject struct used for wikies get operations
type WikiParentObject struct {
	Title string `json:"title"`
}

/* Create */

// WikiCreateObject struct used for wiki create operations
type WikiCreateObject struct {
	Text     string                   `json:"text"`
	Comments string                   `json:"comments,omitempty"`
	Uploads  []AttachmentUploadObject `json:"uploads,omitempty"`
}

/* Update */

// WikiUpdateObject struct used for wiki update operations
type WikiUpdateObject struct {
	Text     string                   `json:"text"`
	Comments string                   `json:"comments,omitempty"`
	Version  int                      `json:"version,omitempty"`
	Uploads  []AttachmentUploadObject `json:"uploads,omitempty"`
}

/* Requests */

// WikiSingleGetRequest contains data for making request to get specified wiki
type WikiSingleGetRequest struct {
	Includes []string
}

/* Internal types */

type wikiAllResult struct {
	WikiPages []WikiMultiObject `json:"wiki_pages"`
}

type wikiSingleResult struct {
	WikiPage WikiObject `json:"wiki_page"`
}

type wikiCreate struct {
	WikiPage WikiCreateObject `json:"wiki_page"`
}

type wikiUpdate struct {
	WikiPage WikiUpdateObject `json:"wiki_page"`
}

// WikiAllGet gets info for all wikies for project with specified ID
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_WikiPages#Getting-the-pages-list-of-a-wiki
func (r *Context) WikiAllGet(projectID string) ([]WikiMultiObject, int, error) {

	var w wikiAllResult

	ur := url.URL{
		Path: "/projects/" + projectID + "/wiki/index.json",
	}

	status, err := r.Get(&w, ur, http.StatusOK)

	return w.WikiPages, status, err
}

// WikiSingleGet gets single wiki info by specific project ID and wiki title
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_WikiPages#Getting-a-wiki-page
//
// Available includes:
// * attachments
func (r *Context) WikiSingleGet(projectID, wikiTitle string, request WikiSingleGetRequest) (WikiObject, int, error) {

	var w wikiSingleResult

	urlParams := url.Values{}

	// Preparing includes
	urlIncludes(&urlParams, request.Includes)

	ur := url.URL{
		Path:     "/projects/" + projectID + "/wiki/" + wikiTitle + ".json",
		RawQuery: urlParams.Encode(),
	}

	status, err := r.Get(&w, ur, http.StatusOK)

	return w.WikiPage, status, err
}

// WikiSingleVersionGet gets single wiki info by specific project ID, wiki title and version
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_WikiPages#Getting-an-old-version-of-a-wiki-page
//
// Available includes:
// * attachments
func (r *Context) WikiSingleVersionGet(projectID, wikiTitle string, version int, request WikiSingleGetRequest) (WikiObject, int, error) {

	var w wikiSingleResult

	urlParams := url.Values{}

	// Preparing includes
	urlIncludes(&urlParams, request.Includes)

	ur := url.URL{
		Path:     "/projects/" + projectID + "/wiki/" + wikiTitle + "/" + strconv.Itoa(version) + ".json",
		RawQuery: urlParams.Encode(),
	}

	status, err := r.Get(&w, ur, http.StatusOK)

	return w.WikiPage, status, err
}

// WikiCreate creates new wiki
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_WikiPages#Creating-or-updating-a-wiki-page
func (r *Context) WikiCreate(projectID, wikiTitle string, wiki WikiCreateObject) (WikiObject, int, error) {

	var w wikiSingleResult

	ur := url.URL{
		Path: "/projects/" + projectID + "/wiki/" + wikiTitle + ".json",
	}

	status, err := r.Put(wikiCreate{WikiPage: wiki}, &w, ur, http.StatusCreated)

	return w.WikiPage, status, err
}

// WikiUpdate updates wiki page
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_WikiPages#Creating-or-updating-a-wiki-page
func (r *Context) WikiUpdate(projectID, wikiTitle string, wiki WikiUpdateObject) (int, error) {

	ur := url.URL{
		Path: "/projects/" + projectID + "/wiki/" + wikiTitle + ".json",
	}

	status, err := r.Put(wikiUpdate{WikiPage: wiki}, nil, ur, http.StatusNoContent)

	return status, err
}

// WikiDelete deletes wiki with specified project ID and title
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_WikiPages#Deleting-a-wiki-page
func (r *Context) WikiDelete(projectID, wikiTitle string) (int, error) {

	ur := url.URL{
		Path: "/projects/" + projectID + "/wiki/" + wikiTitle + ".json",
	}

	status, err := r.Del(nil, nil, ur, http.StatusNoContent)

	return status, err
}
