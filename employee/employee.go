package employee

import (
	"context"
	"github.com/HRMonitorr/PasetoprojectBackend"
	"github.com/HRMonitorr/UsersBackend"
	"github.com/HRMonitorr/githubwrapper"
	"github.com/HRMonitorr/monitoring-backend/structure"
	"net/http"
	"os"
)

func GetDataCommitsAll(PublicKey, MongoEnv, dbname, colname, personalToken, reponame, ownername string, r *http.Request) string {
	req := new(structure.Creds)
	conn := PasetoprojectBackend.MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = http.StatusNotAcceptable
		req.Message = "Header Login Not Found"
	} else {
		cekadmin := UsersBackend.IsAdmin(tokenlogin, PublicKey)
		if cekadmin != true {
			req.Status = http.StatusNotAcceptable
			req.Message = "IHHH Kamu bukan admin"
		}
		checktoken, err := PasetoprojectBackend.DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = http.StatusNotAcceptable
			req.Message = "tidak ada data username : " + tokenlogin
		}
		compared := PasetoprojectBackend.CompareUsername(conn, colname, checktoken)
		if compared != true {
			req.Status = http.StatusNotAcceptable
			req.Message = "Data User tidak ada"
		} else {
			datacomms, err := githubwrapper.ListCommitALL(context.Background(), personalToken, reponame, ownername)
			if err != nil {
				req.Status = http.StatusBadRequest
				req.Message = err.Error()
			}
			req.Status = http.StatusOK
			req.Message = "data User berhasil diambil"
			req.Data = datacomms
		}
	}
	return PasetoprojectBackend.ReturnStringStruct(req)
}
