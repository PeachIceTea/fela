package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/PeachIceTea/fela/conf"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func main() {
	c := conf.Init()
	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		JSONResponse(w, http.StatusOK, conf.M{"msg": "Hello, World!"})
	})

	r.POST("/upload", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uploadFile, header, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			JSONResponse(w, http.StatusInternalServerError, conf.M{"msg": "internal server error"})
			return
		}
		defer uploadFile.Close()

		path := fmt.Sprintf("%s/%s", c.FilePath, header.Filename)

		osFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println(err)
			JSONResponse(w, http.StatusInternalServerError, conf.M{"msg": "internal server error"})
			return
		}
		defer osFile.Close()

		info, err := AudiobookMetadata(path)
		if err != nil {
			fmt.Println(err)
			JSONResponse(w, http.StatusInternalServerError, conf.M{"msg": "internal server error"})
			return
		}

		tx, err := c.DB.Beginx()
		if err != nil {
			fmt.Println(err)
			JSONResponse(w, http.StatusInternalServerError, conf.M{"msg": "internal server error"})
			return
		}
		defer tx.Commit()

		sqlStr, _ := c.TemplateString("file_insert", nil)

		infoStr, err := json.Marshal(info)
		if err != nil {
			fmt.Println(err)
			JSONResponse(w, http.StatusInternalServerError, conf.M{"msg": "internal server error"})
			return
		}

		res, err := tx.Exec(sqlStr, header.Filename, infoStr, "audiobook")
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			JSONResponse(w, http.StatusInternalServerError, conf.M{"msg": "internal server error"})
			return
		}
		if affected, _ := res.RowsAffected(); affected == 0 {
			JSONResponse(w, http.StatusConflict, conf.M{"error": "username already in use"})
			return
		}

		JSONResponse(w, 200, conf.M{"msg": fmt.Sprintf("%s was successfully uploaded.", header.Filename), "info": info})
	})

	r.ServeFiles("/files/*filepath", http.Dir(c.FilePath))

	n := negroni.Classic()
	n.Use(cors.New(cors.Options{}))
	n.UseHandler(r)

	n.Run(":8080")
}

func AudiobookMetadata(path string) (info AudiobookInfo, err error) {
	cmd := exec.Command("ffprobe", "-i", path, "-print_format", "json", "-v", "quiet", "-show_format", "-show_chapters", "-show_streams")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	if err = json.NewDecoder(stdout).Decode(&info); err != nil {
		return
	}

	err = cmd.Wait()
	return
}

// JSONBody - Get JSON Body from request
func JSONBody(e interface{}, r *http.Request) (err error) {
	err = json.NewDecoder(r.Body).Decode(e)
	return
}

// JSONResponse - Respond with JSON
func JSONResponse(w http.ResponseWriter, status int, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(b)
	if err != nil {
		panic(err)
	}
}

type AudiobookInfo struct {
	Streams []struct {
		Index          int    `json:"index"`
		CodecName      string `json:"codec_name"`
		CodecLongName  string `json:"codec_long_name"`
		Profile        string `json:"profile,omitempty"`
		CodecType      string `json:"codec_type"`
		CodecTimeBase  string `json:"codec_time_base,omitempty"`
		CodecTagString string `json:"codec_tag_string"`
		CodecTag       string `json:"codec_tag"`
		SampleFmt      string `json:"sample_fmt,omitempty"`
		SampleRate     string `json:"sample_rate,omitempty"`
		Channels       int    `json:"channels,omitempty"`
		ChannelLayout  string `json:"channel_layout,omitempty"`
		BitsPerSample  int    `json:"bits_per_sample,omitempty"`
		RFrameRate     string `json:"r_frame_rate"`
		AvgFrameRate   string `json:"avg_frame_rate"`
		TimeBase       string `json:"time_base"`
		StartPts       int    `json:"start_pts"`
		StartTime      string `json:"start_time"`
		DurationTs     int64  `json:"duration_ts"`
		Duration       string `json:"duration"`
		BitRate        string `json:"bit_rate,omitempty"`
		MaxBitRate     string `json:"max_bit_rate,omitempty"`
		NbFrames       string `json:"nb_frames,omitempty"`
		Disposition    struct {
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Original        int `json:"original"`
			Comment         int `json:"comment"`
			Lyrics          int `json:"lyrics"`
			Karaoke         int `json:"karaoke"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			VisualImpaired  int `json:"visual_impaired"`
			CleanEffects    int `json:"clean_effects"`
			AttachedPic     int `json:"attached_pic"`
			TimedThumbnails int `json:"timed_thumbnails"`
		} `json:"disposition"`
		Tags struct {
			Language    string `json:"language"`
			HandlerName string `json:"handler_name"`
		} `json:"tags,omitempty"`
		Width              int    `json:"width,omitempty"`
		Height             int    `json:"height,omitempty"`
		CodedWidth         int    `json:"coded_width,omitempty"`
		CodedHeight        int    `json:"coded_height,omitempty"`
		HasBFrames         int    `json:"has_b_frames,omitempty"`
		SampleAspectRatio  string `json:"sample_aspect_ratio,omitempty"`
		DisplayAspectRatio string `json:"display_aspect_ratio,omitempty"`
		PixFmt             string `json:"pix_fmt,omitempty"`
		Level              int    `json:"level,omitempty"`
		ColorRange         string `json:"color_range,omitempty"`
		ColorSpace         string `json:"color_space,omitempty"`
		ChromaLocation     string `json:"chroma_location,omitempty"`
		Refs               int    `json:"refs,omitempty"`
		BitsPerRawSample   string `json:"bits_per_raw_sample,omitempty"`
	} `json:"streams"`
	Chapters []struct {
		ID        int    `json:"id"`
		TimeBase  string `json:"time_base"`
		Start     int    `json:"start"`
		StartTime string `json:"start_time"`
		End       int    `json:"end"`
		EndTime   string `json:"end_time"`
		Tags      struct {
			Title string `json:"title"`
		} `json:"tags"`
	} `json:"chapters"`
	Format struct {
		Filename       string `json:"filename"`
		NbStreams      int    `json:"nb_streams"`
		NbPrograms     int    `json:"nb_programs"`
		FormatName     string `json:"format_name"`
		FormatLongName string `json:"format_long_name"`
		StartTime      string `json:"start_time"`
		Duration       string `json:"duration"`
		Size           string `json:"size"`
		BitRate        string `json:"bit_rate"`
		ProbeScore     int    `json:"probe_score"`
		Tags           struct {
			MajorBrand       string `json:"major_brand"`
			MinorVersion     string `json:"minor_version"`
			CompatibleBrands string `json:"compatible_brands"`
			Title            string `json:"title"`
			Artist           string `json:"artist"`
			AlbumArtist      string `json:"album_artist"`
			Composer         string `json:"composer"`
			Album            string `json:"album"`
			Date             string `json:"date"`
			Encoder          string `json:"encoder"`
			Comment          string `json:"comment"`
			Genre            string `json:"genre"`
			Copyright        string `json:"copyright"`
		} `json:"tags"`
	} `json:"format"`
}
