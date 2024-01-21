package monitoring_backend

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/HRMonitorr/PasetoprojectBackend"
	"github.com/HRMonitorr/UsersBackend"
	"github.com/HRMonitorr/githubwrapper"
	"github.com/HRMonitorr/monitoring-backend/employee"
	"github.com/HRMonitorr/monitoring-backend/structure"
	"net/http"
	"os"
	"time"
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
				datacomms, err := githubwrapper.ListCommitALL(context.Background(), os.Getenv(personalToken), "UsersBackend", "HRMonitorr")
				if err != nil {
					req.Status = http.StatusBadRequest
					req.Message = err.Error()
				}
				if len(datacomms) == 0 {
					req.Status = http.StatusNotAcceptable
					req.Message = "data tidak ditemukan"
				}
				datas := make([]structure.Commits, 0)
				for _, v := range datacomms {
					data := structure.Commits{
						Author:  *v.Author.Name,
						Repos:   *v.Commit.URL,
						Email:   *v.Author.Email,
						Comment: *v.Commit.Message,
						Date:    time.Now(),
					}
					datas = append(datas, data)
				}

				_, err = employee.InsertCommitsManyToDB(conn, datas)
				if err != nil {
					req.Status = http.StatusBadRequest
					req.Message = err.Error()
				}
				req.Status = http.StatusOK
				req.Message = fmt.Sprintf("data Commit berhasil diambil"+
					"%s ", os.Getenv(personalToken))
				req.Data = datas
			}
		}
	}

	return PasetoprojectBackend.ReturnStringStruct(req)
}
