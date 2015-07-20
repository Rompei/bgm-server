package main

import (
	"bgm-server/Godeps/_workspace/src/github.com/astaxie/beego/orm"
	"bgm-server/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//LikeをPOSTする関数。もし動画がなければ作成する。
func likeUpdate(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	checkError(w, err)
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(400)
		utils.CheckError(w, err)
	}

	w = utils.SetJSONHeader(w)

	var video Video

	//jsonをパース
	if err := json.Unmarshal(body, &video); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			utils.CheckError(w, err)
		}
	}

	//DBに保存
	if created, _, err := o.ReadOrCreate(&video, "videoId"); err == nil {
		if created {
			if video.HighThumbnail != nil {
				o.Insert(video.HighThumbnail)
			}
			if video.MediumThumbnail != nil {
				o.Insert(video.MediumThumbnail)
			}
		}
		video.Liked = video.Liked + 1
		if _, err := o.Update(&video); err == nil {
			w.WriteHeader(200)
			response, err := json.Marshal(video)
			utils.CheckError(w, err)
			fmt.Fprintln(w, string(response))
		}
	}
}
