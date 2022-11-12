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
	Username   string `db:"username" form:"username"`
	Password   string `db:"password" form:"password"`
	Root       bool   `db:"root" form:"root"`
	Blind      bool   `db:"blind" form:"blind"`
	Onetime    bool   `db:"onetime" form:"onetime"`
	Restricted bool   `db:"restricted" form:"restricted"`
}

type UserFromForm struct {
	Username   string `form:"username"`
	Password   string `form:"password"`
	Root       string `form:"root"`
	Blind      string `form:"blind"`
	Onetime    string `form:"onetime"`
	Restricted string `form:"restricted"`
}

type File struct {
	File       string    `db:"file"`
	Author     string    `db:"author"`
	Timestamp  time.Time `db:"timestamp"`
	Short      string    `db:"short"`
	Ip         string    `db:"ip"`
	Restricted bool      `db:"restricted"`
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
