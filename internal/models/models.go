package models

import "time"

type Config struct {
	Url         string `yaml:"url"`
	Port        string `yaml:"port"`
	Type        string `yaml:"dbtype"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Address     string `yaml:"address"`
	JWTKey      string `yaml:"jwtkey"`
	ShortLength int    `yaml:"shortlength"`
	UploadLimit int    `yaml:"uploadLimit"`
}

type User struct {
	Username   string `db:"Username" form:"username"`
	Password   string `db:"Password" form:"password"`
	Root       bool   `db:"Root" form:"root"`
	Blind      bool   `db:"Blind" form:"blind"`
	Onetime    bool   `db:"Onetime" form:"onetime"`
	Restricted bool   `db:"Restricted" form:"restricted"`
}

type UserFromForm struct {
	Username   string `form:"username"`
	Password   string `form:"password"`
	Root       string `form:"root"`
	Blind      string `form:"blind"`
	Onetime    string `form:"onetime"`
	Restricted string `form:"restricted"`
}

type Short struct {
	Short      string    `db:"Short"`
	Type       string    `db:"Type"`
	Author     string    `db:"Author"`
	Timestamp  time.Time `db:"Timestamp"`
	Ip         string    `db:"Ip"`
	Restricted bool      `db:"Restricted"`
	Downloads  int       `db:"Downloads"`
}

type File struct {
	Filename   string    `db:"Filename"`
	Short      string    `db:"Short"`
	Author     string    `db:"Author"`
	Timestamp  time.Time `db:"Timestamp"`
	Ip         string    `db:"Ip"`
	Restricted bool      `db:"Restricted"`
	Downloads  int       `db:"Downloads"`
}

type Paste struct {
	Text       string    `db:"Text"`
	Short      string    `db:"Short"`
	Author     string    `db:"Author"`
	Timestamp  time.Time `db:"Timestamp"`
	Ip         string    `db:"Ip"`
	Restricted bool      `db:"Restricted"`
	Downloads  int       `db:"Downloads"`
}

type Login struct {
	Username string `form:"username"`
	Password string `form:"password"`
}
type JWT struct {
	Username   string
	Root       bool
	Blind      bool
	Onetime    bool
	Restricted bool
}
