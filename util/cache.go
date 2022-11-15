package util

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
)

type xlsxHashFile struct {
	CRC32Map map[string]uint32
}

type xlsxFileCache struct {
	Name   string
	Sheets []*xlsxSheet
}

type xlsxSheet struct {
	Name  string
	Cells [][]string
}

type TableCache struct {
	z        *zip.ReadCloser
	name     string
	cacheDir string

	originFile *xlsx.File
}

func NewTableCache(name, cachedir string) *TableCache {
	return &TableCache{
		name:     name,
		cacheDir: cachedir,
	}
}

func (self *TableCache) UseCache() bool {
	return self.originFile == nil
}

func (self *TableCache) cacheFileName() string {
	return fmt.Sprintf("%s/%s.cache", self.cacheDir, self.name)
}

func (self *TableCache) hashFileName() string {
	return fmt.Sprintf("%s/%s.hash", self.cacheDir, self.name)
}

func (self *TableCache) Open() error {

	z, err := zip.OpenReader(self.name)
	if err != nil {
		return err
	}

	self.z = z

	return nil
}

func readJsonFile(filename string, m interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}

	return json.NewDecoder(f).Decode(m)
}

func writeJsonFile(filename string, m interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(m)
}

func (self *TableCache) readCache() (xf *xlsx.File, err error) {

	var hashFile xlsxHashFile
	err = readJsonFile(self.hashFileName(), &hashFile)
	if err != nil {
		return nil, nil
	}

	if len(self.z.File) != len(hashFile.CRC32Map) {
		return nil, nil
	}

	for _, f := range self.z.File {
		if hashFile.CRC32Map[f.Name] != f.CRC32 {
			return nil, nil
		}
	}

	var cacheFile xlsxFileCache
	err = readJsonFile(self.cacheFileName(), &cacheFile)

	if err != nil {
		return nil, nil
	}

	xf = xlsx.NewFile()

	for _, s := range cacheFile.Sheets {
		sheet, err := xf.AddSheet(s.Name)
		if err != nil {
			return nil, err
		}
		for _, srcRow := range s.Cells {
			tgtRow := sheet.AddRow()

			for _, c := range srcRow {
				cell := tgtRow.AddCell()
				cell.Value = c
			}

		}
	}

	return
}

func (self *TableCache) Load() (xf *xlsx.File, err error) {

	cfile, err := self.readCache()

	// cache未击中, 从原文件读取
	if err == nil && cfile == nil {
		xf, err = xlsx.ReadZipWithRowLimit(self.z, xlsx.NoRowLimit)
		self.originFile = xf
		return
	}

	return cfile, err
}

func (self *TableCache) Save() error {

	var hashFile xlsxHashFile
	hashFile.CRC32Map = make(map[string]uint32)

	for _, f := range self.z.File {
		hashFile.CRC32Map[f.Name] = f.CRC32
	}

	writeJsonFile(self.hashFileName(), &hashFile)

	var newFile xlsxFileCache
	newFile.Name = self.name

	for _, sheet := range self.originFile.Sheets {

		var newSheet xlsxSheet
		newSheet.Name = sheet.Name

		newSheet.Cells = make([][]string, 0, len(sheet.Rows))

		for _, row := range sheet.Rows {

			var rowData = make([]string, 0, len(row.Cells))
			for _, c := range row.Cells {
				rowData = append(rowData, c.Value)
			}

			newSheet.Cells = append(newSheet.Cells, rowData)
		}

		newFile.Sheets = append(newFile.Sheets, &newSheet)
	}

	writeJsonFile(self.cacheFileName(), &newFile)

	return nil
}
