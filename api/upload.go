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

// File represents a File database row.
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

// Upload - POST /audiobook/upload - Creates a new audiobook.
// Requires the "file" field which has to include 1 or more files.
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
			panic(err)
		}

		var eg errgroup.Group
		files := make([]*File, len(data.Files))
		for i, header := range data.Files {
			eg.Go(func(i int, header *multipart.FileHeader) func() error {
				return func() (err error) {
					file := &File{
						Name:        header.Filename,
						Path:        fmt.Sprintf("%s/%s", dirPath, header.Filename),
						AudiobookID: &audiobookID,
					}

					err = storeFile(file.Path, header, c)
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

					if i == 0 {
						// Extract cover from first file. Errors are ignored, the
						// client is expected to show a cover from the server.
						// Also uses
						coverPath := path.Clean(fmt.Sprintf(
							"%s/cover/%d.jpg",
							c.FilesPath,
							audiobookID,
						))
						extractCover(file.Path, coverPath)

						// Tries to get the name of the book by looking up the
						// album name. The reason there might be none is because
						// the entire book is contained in a single file. We use
						// the title of the file in that case.
						title := file.Metadata.Format.Tags.Album
						if title == "" {
							title = file.Metadata.Format.Tags.Title
						}

						// Extract title and author from first file.
						a := Audiobook{
							ID:     audiobookID,
							Title:  &title,
							Author: &file.Metadata.Format.Tags.Artist,
						}

						if *a.Title != "" && *a.Author != "" {
							_, err = tx.NamedExec(
								c.TemplateWithData("update_audiobook", a),
								a,
							)
						}
					}

					files[i] = file
					return
				}
			}(i, header))
		}
		err = eg.Wait()
		if err != nil {
			tx.Rollback()
			os.RemoveAll(dirPath)
			panic(err)
		}

		a := &Audiobook{}
		err = tx.Get(a, c.TemplateString("get_audiobook"), audiobookID)
		if err != nil {
			panic(err)
		}

		ctx.JSON(
			http.StatusOK,
			conf.M{"audiobook": a},
		)
	})
}

// storeFile writes a single multipart.File to disk. Returns an error if the
// file already exists.
func storeFile(
	path string,
	h *multipart.FileHeader,
	c *conf.Config,
) (err error) {
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

// extractMetadata uses ffprobe to get information about the streams it
// includes as well as any other metadata included in the file, like author,
// title or chapter markers.
// This information can later be used to display the length of a given
// audiobook, provide chapter markers with the player etc.
func extractMetadata(path string) (meta *Metadata, err error) {
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

	meta = &Metadata{}
	if err = json.NewDecoder(stdout).Decode(&meta); err != nil {
		return
	}

	err = cmd.Wait()
	return
}

// extractCover uses ffmpeg to extract the cover included in many audiobook
// files and saves it as a jpg.
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

// usePlaceholderCover copies a placeholder cover to the path. Used when the
// audiobook does not come with one by default.
func usePlaceholderCover(coverPath string) (err error) {
	f, err := os.Create(coverPath)
	if err != nil {
		return err
	}
	defer f.Close()

	placeholder, err := os.Open("placeholder-cover.jpg")
	if err != nil {
		return err
	}
	defer placeholder.Close()

	_, err = io.Copy(f, placeholder)
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

// Scan turns a string or []byte into a Metadata struct.
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

// Value turns struct into string.
func (fm Metadata) Value() (v driver.Value, err error) {
	v, err = json.Marshal(fm)
	return
}

// Codec extracts the codec of the audiostream.
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

// Duration returns the length of the audiostream.
func (fm *Metadata) Duration() (duration float64, err error) {
	for _, stream := range fm.Streams {
		if stream.CodecType == "audio" {
			duration, err = strconv.ParseFloat(stream.Duration, 64)
			return
		}
	}

	err = ErrNoAudioStream
	return
}
