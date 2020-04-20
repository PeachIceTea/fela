package api

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/PeachIceTea/fela/conf"
)

// User represents a single user in the database.
type User struct {
	ID        int64   `db:"id" json:"id"`
	Name      string  `db:"name" json:"name"`
	Password  []byte  `db:"password" json:"-"`
	Role      string  `db:"role" json:"role"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt *string `db:"updated_at" json:"updated_at"`

	Uploads []Audiobook `db:"-" json:"uploads,omitempty"`
}

// GetUsers - GET /user - Gets all users.
func GetUsers(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/user", func(ctx *gin.Context) {
		var users []User
		err := c.DB.Select(&users, c.TemplateString("all_users"))
		if err != nil {
			panic(err)
		}
		ctx.JSON(http.StatusOK, conf.M{"users": users})
	})
}

// GetUser - GET /user/:id - Gets a single user and his uploads.
func GetUser(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/user/:id", func(ctx *gin.Context) {
		id, err := getID(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, conf.M{"err": err.Error()})
			return
		}

		user := User{}
		err = c.DB.Get(&user, c.TemplateString("get_user"), id)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(
					http.StatusNotFound,
					conf.M{"err": "no user with that id"},
				)
				return
			}

			panic(err)
		}

		err = c.DB.Select(&user.Uploads, c.TemplateString("user_uploads"), id)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, conf.M{"user": user})
	})
}

// Register - POST /user/register - Creates new user.
// Requires "name" and "password" as fields. Accepts "role" field, which
// defaults to "user", "uploader" and "admin" are also valid values.
func Register(r *gin.RouterGroup, c *conf.Config) {
	r.POST("/user/register", func(ctx *gin.Context) {
		var data struct {
			Name     string `form:"name" json:"name"`
			Password string `form:"password" json:"password"`
			Role     string `form:"role" json:"role"`
		}

		err := ctx.ShouldBind(&data)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				conf.M{"err": "invalid request body"})
			return
		}
		if len(data.Name) == 0 || len(data.Password) == 0 {
			ctx.JSON(http.StatusBadRequest, conf.M{"err": "data missing"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(data.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			panic(err)
		}

		if len(data.Role) == 0 {
			data.Role = "user"
		}

		data.Name = strings.ToLower(data.Name)

		res, err := c.DB.Exec(
			c.TemplateString("register"), data.Name, hash, data.Role,
		)
		if err != nil {
			// MySQL error code https://mariadb.com/kb/en/mariadb-error-codes/
			if strings.Contains(err.Error(), "1062") {
				ctx.JSON(
					http.StatusConflict, conf.M{"err": "user already exists"},
				)
				return
			}

			panic(err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			panic(err) // Should not happen with MySQL
		}

		ctx.JSON(http.StatusOK, conf.M{"user_id": id})
	})
}

// Login - POST /user/login - Authenticates user and generates JWT.
// Requires "name" and "password" as fields.
func Login(r *gin.RouterGroup, c *conf.Config) {
	r.POST("/user/login", func(ctx *gin.Context) {
		var data struct {
			Name     string `form:"name" json:"name"`
			Password string `form:"password" json:"password"`
		}

		err := ctx.ShouldBind(&data)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest, conf.M{"err": "invalid request body"},
			)
			return
		}
		if len(data.Name) == 0 || len(data.Password) == 0 {
			ctx.JSON(http.StatusBadRequest, conf.M{"err": "data missing"})
			return
		}

		data.Name = strings.ToLower(data.Name)

		user := User{}
		err = c.DB.Get(&user, c.TemplateString("login"), data.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, conf.M{"err": "cannot find user"})
				return
			}

			panic(err)
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				ctx.JSON(
					http.StatusBadRequest,
					conf.M{"err": "passwords do not match"},
				)
				return
			}

			panic(err)
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			user.ID,
			user.Name,
			user.Role,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		}).SignedString(c.Secret)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, conf.M{"token": token, "user": user})
	})
}

// UpdateUser - PUT /user/:id - Updates user.
// Accepts "name", "password" and "role" fields. Only admins may update other
// users and update the role.
func UpdateUser(r *gin.RouterGroup, c *conf.Config) {
	r.PUT("/user/:id", func(ctx *gin.Context) {
		id, err := getID(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, conf.M{"err": err.Error()})
			return
		}

		var data struct {
			ID       int64   `form:"-" json:"-" db:"id"`
			Name     *string `form:"name" json:"name" db:"name"`
			Password *string `form:"password" json:"password" db:"password"`
			Role     *string `form:"role" json:"role" db:"role"`
		}

		err = ctx.ShouldBind(&data)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				conf.M{"err": "invalid request body"},
			)
			return
		}

		if data.Name == nil && data.Password == nil && data.Role == nil {
			ctx.JSON(
				http.StatusBadRequest,
				conf.M{"err": "no fields to update"},
			)
			return
		}

		claims := getClaims(ctx)
		isAdmin := claims.isAdmin()
		if claims.ID != id && !isAdmin {
			ctx.JSON(
				http.StatusUnauthorized,
				conf.M{"err": "no permission to update other users"},
			)
			return
		}

		if data.Role != nil && !isAdmin {
			ctx.JSON(
				http.StatusUnauthorized,
				conf.M{"err": "no permission to update admin column"},
			)
			return
		}

		//TODO: Add more checks when changing password:
		// 1: If user is not an admin require reauthentication
		// 2: If the password of an admin is changed, required reauthentication
		if data.Password != nil {
			hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}
			str := string(hash)
			data.Password = &str
		}

		if data.Name != nil {
			lower := strings.ToLower(*data.Name)
			data.Name = &lower
		}

		data.ID = id
		_, err = c.DB.NamedExec(c.TemplateWithData("update_user", data), data)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, conf.M{"msg": "user updated"})
	})
}

// DeleteUser - DELETE /user/:id - Deletes user.
func DeleteUser(r *gin.RouterGroup, c *conf.Config) {
	r.DELETE("/user/:id", func(ctx *gin.Context) {
		id, err := getID(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, conf.M{"err": err.Error()})
			return
		}

		claims := getClaims(ctx)
		if claims.ID != id && claims.isAdmin() {
			ctx.JSON(
				http.StatusUnauthorized,
				conf.M{"err": "no permission to delete other users"},
			)
			return
		}

		_, err = c.DB.Exec(c.TemplateString("delete_user"), id)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, conf.M{"msg": "user deleted"})
	})
}
