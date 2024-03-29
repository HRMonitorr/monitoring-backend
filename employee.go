package monitoring_backend

import (
	"context"
	"encoding/json"
	"github.com/HRMonitorr/PasetoprojectBackend"
	"github.com/HRMonitorr/UsersBackend"
	"github.com/HRMonitorr/githubwrapper"
	"github.com/HRMonitorr/monitoring-backend/employee"
	"github.com/HRMonitorr/monitoring-backend/structure"
	"net/http"
	"os"
)

func GetDataCommitsAll(MongoEnv, dbname, personalToken string, r *http.Request) string {
	req := new(structure.Creds)
	//conn := PasetoprojectBackend.MongoCreateConnection(MongoEnv, dbname)
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
			datacomms, err := githubwrapper.ListCommitALL(context.Background(), os.Getenv(personalToken), datauser.RepoName, datauser.OwnerName)
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
					Author:  v.Commit.Author.Name,
					Repos:   v.Commit.URL,
					Email:   v.Commit.Author.Email,
					Comment: v.Commit.Message,
				}
				datas = append(datas, data)
			}
			//
			//_, err = employee.InsertCommitsManyToDB(conn, datas)
			//if err != nil {
			//	req.Status = http.StatusBadRequest
			//	req.Message = err.Error()
			//}
			req.Status = http.StatusOK
			req.Message = "data Commit berhasil diambil"
			req.Data = datas
		}
	}
	return PasetoprojectBackend.ReturnStringStruct(req)
}

func GetAndInsertCommits(Publickey, MongoEnv, dbname, personalToken string, r *http.Request) string {
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
			checkadmin := UsersBackend.IsAdmin(tokenlogin, os.Getenv(Publickey))
			if !checkadmin {
				checkHR := UsersBackend.IsHR(tokenlogin, os.Getenv(Publickey))
				if !checkHR {
					req.Status = http.StatusNotAcceptable
					req.Message = "Anda tidak bisa Insert data karena bukan HR atau admin"
				}
			}
			datacomms, err := githubwrapper.ListCommitALL(context.Background(), os.Getenv(personalToken), datauser.RepoName, datauser.OwnerName)
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
					Author:  v.Commit.Author.Name,
					Repos:   v.Commit.URL,
					Email:   v.Commit.Author.Email,
					Comment: v.Commit.Message,
				}
				datas = append(datas, data)
			}

			_, err = employee.InsertCommitsManyToDB(conn, datas)
			if err != nil {
				req.Status = http.StatusBadRequest
				req.Message = err.Error()
			}
			req.Status = http.StatusOK
			req.Message = "data Commit berhasil diambil"
			req.Data = datas
		}
	}
	return PasetoprojectBackend.ReturnStringStruct(req)
}

func GetListRepositories(Publickey, personalToken string, r *http.Request) string {
	req := new(structure.Creds)
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
			checkadmin := UsersBackend.IsAdmin(tokenlogin, os.Getenv(Publickey))
			if !checkadmin {
				checkHR := UsersBackend.IsHR(tokenlogin, os.Getenv(Publickey))
				if !checkHR {
					req.Status = http.StatusNotAcceptable
					req.Message = "Anda tidak bisa Insert data karena bukan HR atau admin"
				}
			}
			datacomms, err := githubwrapper.ListRepositoriesOnlydDetail(context.Background(), os.Getenv(personalToken), datauser.OwnerName)
			if err != nil {
				req.Status = http.StatusBadRequest
				req.Message = err.Error()
			}
			if len(datacomms) == 0 {
				req.Status = http.StatusNotAcceptable
				req.Message = "data tidak ditemukan"
			}
			req.Status = http.StatusOK
			req.Message = "data Commit berhasil diambil"
			req.Data = datacomms
		}
	}
	return PasetoprojectBackend.ReturnStringStruct(req)
}
