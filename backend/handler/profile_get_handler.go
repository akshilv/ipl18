package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/smjn/ipl18/backend/db"
	"github.com/smjn/ipl18/backend/errors"
	"github.com/smjn/ipl18/backend/models"
	"github.com/smjn/ipl18/backend/util"
)

// UserGetHandler .
type UserGetHandler struct {
}

const (
	qFetchUserDetails = "select firstname, lastname, coin, alias, piclocation, inumber, points from ipluser where inumber=$1"
)

var (
	errUserNotFound = fmt.Errorf("user not found")
)

func (p UserGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("UserGetHandler: request to get profile")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked")
		}
	}()

	inumber, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "UserGetHandler: could not get username from token")

	pathVar := mux.Vars(r)
	if pathVar["inumber"] != inumber {
		errors.ErrWriterPanic(w, http.StatusForbidden, errors.ErrTokenInfoMismatch, errors.ErrTokenInfoMismatch, fmt.Sprintf("UserPutHandler: token info and path var mismatch %s-%s", pathVar["inumber"], inumber))
	}

	var pic sql.NullString
	info := models.ProfileViewModel{}

	if err := db.DB.QueryRow(qFetchUserDetails, inumber).Scan(&info.Firstname, &info.Lastname, &info.Coin, &info.Alias, &pic, &info.INumber, &info.Points); err == sql.ErrNoRows {
		errors.ErrWriterPanic(w, http.StatusForbidden, err, errUserNotFound, "UserGetHandler: user not found in db")
	}
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "UserGetHandler: could not query db")

	info.PicLocation = pic.String
	log.Println("UserGetHandler:", info)
	util.StructWriter(w, &info)
}
