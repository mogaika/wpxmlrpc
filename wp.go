package wp

import (
	"encoding/base64"

	"github.com/kolo/xmlrpc"
)

type WordpressRpc struct {
	BlogId   int
	Username string
	Password string
	client   *xmlrpc.Client
}

func NewWordpressRpc(blog_id int, username, password string, url string) (*WordpressRpc, error) {
	rpcclient, err := xmlrpc.NewClient(url, nil)
	if err != nil {
		return nil, err
	}

	return &WordpressRpc{
		BlogId:   blog_id,
		Username: username,
		Password: password,
		client:   rpcclient,
	}, nil
}

func (wp *WordpressRpc) Rpc(method string, args []interface{}, reply interface{}) error {
	resp := make([]interface{}, 3)
	resp[0] = wp.BlogId
	resp[1] = wp.Username
	resp[2] = wp.Password
	if args == nil {
		return wp.client.Call(method, resp[:3], reply)
	} else {
		for i := range args {
			resp = append(resp, args[i])
		}
		return wp.client.Call(method, resp, reply)
	}
}

type ResponceWpUploadFile struct {
	Id   string `xmlrpc:"id"`
	File string `xmlrpc:"file"`
	Url  string `xmlrpc:"url"`
	Type string `xmlrpc:"type"`
}

func (wp *WordpressRpc) WpUploadFile(name, mime string, data []byte, optional map[string]interface{}) (*ResponceWpUploadFile, error) {
	args := NewOptional().
		SetIf(name != "", "name", name).
		SetIf(mime != "", "type", mime).
		Set("bits", xmlrpc.Base64(base64.StdEncoding.EncodeToString(data))).
		AddMap(optional)

	fresp := new(ResponceWpUploadFile)
	return fresp, wp.Rpc("wp.uploadFile", []interface{}{args}, &fresp)
}

type ResponceWpGetTaxonomy struct {
	Name         string            `xmlrpc:"name"`
	Label        string            `xmlrpc:"label"`
	Hierarchical bool              `xmlrpc:"hierarchical"`
	Public       bool              `xmlrpc:"public"`
	ShowUi       bool              `xmlrpc:"show_ui"`
	Builtin      bool              `xmlrpc:"_builtin"`
	Labels       map[string]string `xmlrpc:"labels"`
	Cap          map[string]string `xmlrpc:"cap"`
}

func (wp *WordpressRpc) WpGetTaxonomy(taxonomy string) (*ResponceWpGetTaxonomy, error) {
	fresp := new(ResponceWpGetTaxonomy)
	return fresp, wp.Rpc("wp.getTaxonomy", []interface{}{taxonomy}, &fresp)
}

type ResponceWpGetTaxonomies []ResponceWpGetTaxonomy

func (wp *WordpressRpc) WpGetTaxonomies() (ResponceWpGetTaxonomies, error) {
	fresp := make(ResponceWpGetTaxonomies, 0)
	return fresp, wp.Rpc("wp.getTaxonomies", nil, &fresp)
}

type ResponceWpGetTerm struct {
	TermId         string `xmlrpc:"term_id"`
	Name           string `xmlrpc:"name"`
	Slug           string `xmlrpc:"slug"`
	TermGroup      string `xmlrpc:"term_group"`
	TermTaxonomyId string `xmlrpc:"term_taxonomy_id"`
	Taxonomy       string `xmlrpc:"taxonomy"`
	Description    string `xmlrpc:"description"`
	Parent         string `xmlrpc:"parent"`
	Count          int    `xmlrpc:"count"`
}

func (wp *WordpressRpc) WpGetTerm(taxonomy string, term_id int) (*ResponceWpGetTerm, error) {
	fresp := new(ResponceWpGetTerm)
	return fresp, wp.Rpc("wp.getTerm", []interface{}{taxonomy, term_id}, &fresp)
}

type ResponceWpGetTerms []ResponceWpGetTerm

// filter - optional
func (wp *WordpressRpc) WpGetTerms(taxonomy string, filter map[string]interface{}) (ResponceWpGetTerms, error) {
	fresp := make(ResponceWpGetTerms, 0)
	args := []interface{}{taxonomy}
	if filter != nil {
		args = append(args, filter)
	}
	return fresp, wp.Rpc("wp.getTerms", args, &fresp)
}

// return id
// slug, desc and parent - optionals, can be "" and 0
func (wp *WordpressRpc) WpNewTerm(name, taxonomy string, optional map[string]interface{}) (string, error) {
	args := make(map[string]interface{})
	args["name"] = name
	args["taxonomy"] = taxonomy
	if optional != nil {
		for k, v := range optional {
			args[k] = v
		}
	}

	var fresp string
	return fresp, wp.Rpc("wp.newTerm", []interface{}{args}, &fresp)
}

// return id
func (wp *WordpressRpc) WpNewPost(post_type, status, title, content, slug string, optional map[string]interface{}) (string, error) {
	args := NewOptional().
		SetIf(post_type != "", "post_type", post_type).
		SetIf(status != "", "post_status", status).
		SetIf(title != "", "post_title", title).
		SetIf(content != "", "post_content", content).
		SetIf(slug != "", "post_name", slug).
		AddMap(optional)

	var fresp string
	return fresp, wp.Rpc("wp.newPost", []interface{}{args}, &fresp)
}

func (wp *WordpressRpc) WpEditPost(post_id int, changes map[string]interface{}) (bool, error) {
	var fresp bool
	return fresp, wp.Rpc("wp.editPost", []interface{}{post_id, changes}, &fresp)
}
