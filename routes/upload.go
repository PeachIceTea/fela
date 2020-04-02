package routes

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/julienschmidt/httprouter"

	"github.com/PeachIceTea/fela/conf"
)

// Upload - /upload - Handles file uploads
func Upload(r *httprouter.Router, c *conf.Config) {
	r.POST("/upload", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uploadFile, header, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			return
		}
		defer uploadFile.Close()

		filename, err := hashFile(uploadFile)
		if err != nil {
			fmt.Println(err)
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			return
		}

		path := fmt.Sprintf("%s/%s", c.FilePath, filename)

		_, err = os.Stat(path)
		if err == nil {
			conf.JSONResponse(w, http.StatusConflict, conf.M{"error": "file already exists"})
			return
		}

		if err := storeFile(path, uploadFile); err != nil {
			fmt.Println(err)
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			return
		}

		fileID, err := (func() (fileID int64, err error) {
			info, err := audiobookMetadata(path)
			if err != nil {
				return
			}

			tx, err := c.DB.Beginx()
			if err != nil {
				return
			}
			defer tx.Commit()

			infoStr, err := json.Marshal(info)
			if err != nil {
				return
			}

			res, err := tx.Exec(c.TemplateString("file_insert"), header.Filename, filename, infoStr, "audiobook")
			if err != nil {
				tx.Rollback()
				return
			}

			fileID, _ = res.LastInsertId()
			return
		})()
		if err != nil {
			fmt.Println(err)
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			os.Remove(path)
			return
		}
		if fileID == 0 {
			// Weird limbo where the file didn't exist on disk but did exist
			// in the database. Might happen if user deletes file on disk.
			conf.JSONResponse(w, http.StatusOK, conf.M{"msg": "file reupload", "file_id": 0})
			return
		}

		conf.JSONResponse(w, 200, conf.M{"file_id": fileID})

		return
	})

}

type file interface {
	io.Reader
	io.Seeker
}

func hashFile(f file) (s string, err error) {
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	sha := sha1.New()
	_, err = io.Copy(sha, f)
	if err != nil {
		return
	}

	return hex.EncodeToString(sha.Sum(nil)), nil
}

func storeFile(path string, f file) (err error) {
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	df, err := os.Create(path)
	if err != nil {
		return
	}
	defer df.Close()

	_, err = io.Copy(df, f)
	return
}

func audiobookMetadata(path string) (info conf.AudiobookInfo, err error) {
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
