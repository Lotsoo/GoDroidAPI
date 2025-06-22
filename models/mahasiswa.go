package models

import "gorm.io/gorm"

type Mahasiswa struct {
	gorm.Model
	NIM           string `gorm:"uniqueIndex;not null; size=10" json:"nim" binding:"required,len=10"`
	NamaMahasiswa string `gorm:"not null;column:nama_mahasiswa" json:"nama_mahasiswa" binding:"required"`
	Alamat        string `gorm:"not null" binding:"required" json:"alamat"`
	Jurusan       string `gorm:"not null" binding:"required" json:"jurusan"`
}
