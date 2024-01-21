package employee

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/HRMonitorr/PasetoprojectBackend"
	"github.com/HRMonitorr/UsersBackend"
	"github.com/HRMonitorr/githubwrapper"
	"github.com/HRMonitorr/monitoring-backend/structure"
	"net/http"
	"os"
)

func GetDataCommitsAll(PublicKey, MongoEnv, dbname, colname, personalToken string, r *http.Request) string {
	req := new(structure.Creds)
	conn := PasetoprojectBackend.MongoCreateConnection(MongoEnv, dbname)
	var datauser structure.BodyReq
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		req.Message = "error parsing application/json: " + err.Error()
	} else {
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
				datacomms, err := githubwrapper.ListCommitALL(context.Background(), personalToken, datauser.RepoName, datauser.OwnerName)
				if err != nil {
					req.Status = http.StatusBadRequest
					req.Message = err.Error()
				}
				req.Status = http.StatusOK
				req.Message = fmt.Sprintf("data Commit berhasil diambil"+
					"body : %+v\n", datauser)
				req.Data = datacomms
			}
		}
	}

	return PasetoprojectBackend.ReturnStringStruct(req)
}
