package login

import (
	"net/http"

	"github.com/zitadel/zitadel/v2/internal/domain"
)

const (
	tmplChangeUsername     = "changeusername"
	tmplChangeUsernameDone = "changeusernamedone"
)

type changeUsernameData struct {
	Username string `schema:"username"`
}

func (l *Login) renderChangeUsername(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, "UsernameChange.Title", "UsernameChange.Description", errID, errMessage)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplChangeUsername], data, nil)
}

func (l *Login) handleChangeUsername(w http.ResponseWriter, r *http.Request) {
	data := new(changeUsernameData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	_, err = l.command.ChangeUsername(setContext(r.Context(), authReq.UserOrgID), authReq.UserOrgID, authReq.UserID, data.Username)
	if err != nil {
		l.renderChangeUsername(w, r, authReq, err)
		return
	}
	l.renderChangeUsernameDone(w, r, authReq)
}

func (l *Login) renderChangeUsernameDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	var errType, errMessage string
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, "UsernameChangeDone.Title", "UsernameChangeDone.Description", errType, errMessage)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplChangeUsernameDone], data, nil)
}
