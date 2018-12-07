package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golang/glog"

	"jingoal.com/dfs-client-golang/client"
	"jingoal.com/dfs-client-golang/cmd"
	"jingoal.com/dfs-client-golang/proto/transfer"
)

const (
	mainPage = `<html><head></head><body>
	Example<br><br>
	PutFile  http://localhost:2121/putfile?domain=2&size=1111 <br>
	GetFile  http://localhost:2121/getfile?domain=2&id=59af5fbdbccf53149a47216c <br>
	Dupl     http://localhost:2121/dupl?domain=2&id=59af5fbdbccf53149a47216c <br>
	GetByMd5 http://localhost:2121/getbymd5?domain=2&md5=b8c3850a7f8bb9d1b39bad3f11301db7 <br>
	Delete   http://localhost:2121/delfile?domain=2&id=59af5fbdbccf53149a47216c <br>
	</body></html>`
)

var (
	listenAddr = flag.String("listen-addr", ":2121", "listen address")
)

func main() {
	flag.Parse()
	go client.Initialize()

	http.HandleFunc("/", mainpage)
	http.HandleFunc("/getfile", getfile)
	http.HandleFunc("/putfile", putfile)
	http.HandleFunc("/delfile", delfile)
	http.HandleFunc("/dupl", dupl)
	http.HandleFunc("/exists", exists)
	http.HandleFunc("/getbymd5", getbymd5)
	http.HandleFunc("/existbymd5", existbymd5)
	http.HandleFunc("/stat", stat)

	glog.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func mainpage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(mainPage))
}

func dupl(w http.ResponseWriter, r *http.Request) {
	info, err := parseInfo(r.URL.RawQuery)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Duplicate file error, %s.\n", err)))
		return
	}

	did, err := client.Duplicate(info.Id, info.Domain, 0)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Duplicate file error, %s.\n", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Duplicate file ok %s [%+v].\n", did, info)))
}

func exists(w http.ResponseWriter, r *http.Request) {
	info, err := parseInfo(r.URL.RawQuery)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Exists file error, %s.\n", err)))
		return
	}

	result, err := client.Exists(info.Id, info.Domain, 0)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Exists file error, %s.\n", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Exists file ok %t [%+v].\n", result, info)))
}

func getbymd5(w http.ResponseWriter, r *http.Request) {
	info, err := parseInfo(r.URL.RawQuery)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Get by md5 error, %s.\n", err)))
		return
	}

	did, err := client.GetByMd5(info.Md5, info.Domain, info.Size, 0)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Get by md5 error, %s.\n", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Get by md5 ok %s [%+v].\n", did, info)))
}

func existbymd5(w http.ResponseWriter, r *http.Request) {
	info, err := parseInfo(r.URL.RawQuery)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Exist by md5 error, %s.\n", err)))
		return
	}

	result, err := client.ExistByMd5(info.Md5, info.Domain, info.Size, 0)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Exists by md5 error, %s.\n", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Exists by md5 ok %t [%+v].\n", result, info)))
}

func stat(w http.ResponseWriter, r *http.Request) {
	info, err := parseInfo(r.URL.RawQuery)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Stat file error, %s.\n", err)))
		return
	}

	inf, err := client.Stat(info.Id, info.Domain, 0)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Stat file error, %s.\n", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Stat file ok %+v.\n", inf, info)))
}

func delfile(w http.ResponseWriter, r *http.Request) {
	info, err := parseInfo(r.URL.RawQuery)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Delete file error, %s.\n", err)))
		return
	}

	err = client.Delete(info.Id, info.Domain, 0)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Delete file error, %s.\n", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Delete file ok %+v.\n", info)))
}

func putfile(w http.ResponseWriter, r *http.Request) {
	info, err := parseInfo(r.URL.RawQuery)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Put file error, %s.\n", err)))
		return
	}

	inf, err := cmd.WriteFile(info.Size, 1048576 /*chunk size */, info.Domain, info.Biz, 0)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Put file error, %s.\n", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Put file ok %+v\n", inf)))
}

func getfile(w http.ResponseWriter, r *http.Request) {
	info, err := parseInfo(r.URL.RawQuery)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Get file error, %s.\n", err)))
		return
	}

	err = cmd.ReadFile(info, 0)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Get file error, %s.\n", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Get file ok %+v.\n", info)))
}

func parseInfo(query string) (*transfer.FileInfo, error) {
	params, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}
	info := transfer.FileInfo{
		Id:   first(params["id"]),
		Name: first(params["fn"]),
		Md5:  first(params["md5"]),
		Biz:  first(params["biz"]),
	}

	domain, err := strconv.ParseInt(first(params["domain"]), 10, 64)
	if err == nil {
		info.Domain = domain
	}
	size, err := strconv.ParseInt(first(params["size"]), 10, 64)
	if err == nil {
		info.Size = size
	}
	user, err := strconv.ParseInt(first(params["user"]), 10, 64)
	if err == nil {
		info.User = user
	}

	return &info, nil
}

func first(ss []string) string {
	if len(ss) > 0 {
		return ss[0]
	}

	return ""
}
