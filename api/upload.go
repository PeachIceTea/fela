package api

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/PeachIceTea/fela/conf"
)

// Declare errors
var (
	ErrNoAudioStream = errors.New("no audio stream")
)

// File represents a File database row
type File struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Codec       string    `db:"codec" json:"codec"`
	Duration    float64   `db:"duration" json:"duration"`
	Metadata    *Metadata `db:"metadata" json:"metadata"`
	Path        string    `db:"path" json:"path"`
	AudiobookID *int64    `db:"audiobook" json:"audiobook_id"`
	CreatedAt   string    `db:"created_at" json:"created_at"`
	UpdatedAt   *string   `db:"updated_at" json:"updated_at"`
}

// Upload - POST /user/register - Creates new user
func Upload(r *gin.RouterGroup, c *conf.Config) {
	r.POST("/audiobook/upload", func(ctx *gin.Context) {
		claims := getClaims(ctx)
		if !claims.isUploader() {
			ctx.JSON(
				http.StatusUnauthorized,
				conf.M{"err": "no permission to upload"},
			)
			return
		}

		var data struct {
			Files []*multipart.FileHeader `form:"file"`
		}

		err := ctx.ShouldBind(&data)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				conf.M{"err": "invalid request body"},
			)
			return
		}
		if len(data.Files) == 0 {
			ctx.JSON(http.StatusBadRequest, conf.M{"err": "no files"})
			return
		}

		tx, err := c.DB.Beginx()
		if err != nil {
			panic(err)
		}
		defer tx.Commit()

		res, err := tx.Exec(c.TemplateString("new_audiobook"), claims.ID)
		if err != nil {
			panic(err)
		}

		audiobookID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		dirPath := path.Clean(
			fmt.Sprintf("%s/audio/%d", c.FilesPath, audiobookID),
		)
		err = os.Mkdir(dirPath, os.ModePerm|os.ModeDir)
		if err != nil {
			tx.Rollback()
			panic(err) // file permission issues?
		}

		var eg errgroup.Group
		files := make([]*File, len(data.Files))
		for i, h := range data.Files {
			eg.Go(func(i int, h *multipart.FileHeader) func() error {
				return func() (err error) {
					file := &File{
						Name:        h.Filename,
						Path:        fmt.Sprintf("%s/%s", dirPath, h.Filename),
						AudiobookID: &audiobookID,
					}

					err = storeFile(file.Path, h, c)
					if err != nil {
						return
					}

					file.Metadata, err = extractMetadata(file.Path)
					if err != nil {
						return
					}

					file.Codec, err = file.Metadata.Codec()
					if err != nil {
						return
					}

					file.Duration, err = file.Metadata.Duration()
					if err != nil {
						return
					}

					_, err = tx.NamedExec(c.TemplateString("new_file"), file)
					if err != nil {
						return
					}

					// Extract cover from first file
					if i == 0 {
						extractCover(
							file.Path,
							path.Clean(fmt.Sprintf(
								"%s/cover/%d.jpg",
								c.FilesPath,
								audiobookID,
							)),
						)
					}

					files[i] = file
					return
				}
			}(i, h))
		}
		err = eg.Wait()
		if err != nil {
			tx.Rollback()
			os.RemoveAll(dirPath)
			panic(err)
		}

		ctx.JSON(
			http.StatusOK,
			conf.M{"audiobook_id": audiobookID, "files": files},
		)
	})
}

func storeFile(path string, h *multipart.FileHeader, c *conf.Config) (err error) {
	f, err := h.Open()
	if err != nil {
		return
	}
	defer f.Close()

	disk, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm)
	if err != nil {
		return
	}

	_, err = io.Copy(disk, f)
	if err != nil {
		return
	}

	return
}

func extractMetadata(path string) (fm *Metadata, err error) {
	cmd := exec.Command(
		"ffprobe",
		"-i", path,
		"-print_format", "json",
		"-v", "quiet",
		"-show_format", "-show_chapters", "-show_streams",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	fm = &Metadata{}
	if err = json.NewDecoder(stdout).Decode(&fm); err != nil {
		return
	}

	err = cmd.Wait()
	return
}

func extractCover(filePath, coverPath string) (err error) {
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		filePath,
		"-an",
		"-vcodec",
		"copy",
		coverPath,
	)

	if err = cmd.Start(); err != nil {
		return
	}

	err = cmd.Wait()
	return
}

// Metadata stores the output of ffprobe
type Metadata struct {
	Streams []struct {
		CodecName     string `json:"codec_name"`
		CodecType     string `json:"codec_type"`
		SampleRate    string `json:"sample_rate,omitempty"`
		Channels      int    `json:"channels,omitempty"`
		ChannelLayout string `json:"channel_layout,omitempty"`
		TimeBase      string `json:"time_base"`
		StartPts      int    `json:"start_pts"`
		StartTime     string `json:"start_time"`
		DurationTs    int64  `json:"duration_ts"`
		Duration      string `json:"duration"`
		BitRate       string `json:"bit_rate,omitempty"`
		MaxBitRate    string `json:"max_bit_rate,omitempty"`
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
		Filename   string `json:"filename"`
		FormatName string `json:"format_name"`
		StartTime  string `json:"start_time"`
		Duration   string `json:"duration"`
		BitRate    string `json:"bit_rate"`
		Tags       struct {
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

// Scan turns string into AudiobookInfo struct - used by database/sql
func (fm *Metadata) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var source []byte

	switch src.(type) {
	case string:
		{
			source = []byte(src.(string))
		}

	case []byte:
		{
			source = src.([]byte)
		}

	default:
		{
			return errors.New("incompatible type for FileInfo")
		}
	}

	return json.Unmarshal(source, fm)
}

// Value turns struct into string - used by database/sql
func (fm Metadata) Value() (v driver.Value, err error) {
	v, err = json.Marshal(fm)
	return
}

// Codec extracts the codec of the audiostream
func (fm *Metadata) Codec() (c string, err error) {
	for _, stream := range fm.Streams {
		if stream.CodecType == "audio" {
			c = stream.CodecName
			return
		}
	}
	err = ErrNoAudioStream
	return
}

// Duration returns the duration of the audiostream
func (fm *Metadata) Duration() (d float64, err error) {
	for _, stream := range fm.Streams {
		if stream.CodecType == "audio" {
			d, err = strconv.ParseFloat(stream.Duration, 64)
			return
		}
	}

	err = ErrNoAudioStream
	return
}
