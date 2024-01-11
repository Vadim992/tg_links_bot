package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"links_tg-bot/lib/e"
	"links_tg-bot/storage"
	"math/rand"
	"os"
	"path/filepath"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func NewStorage(basePath string) Storage {
	return Storage{basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)

	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)

	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (p *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	isExistDir, err := s.isExistDir(path)

	if err != nil {
		return nil, err
	}

	if !isExistDir {
		return nil, storage.ErrorNoSavedPages
	}

	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrorNoSavedPages
	}

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) (err error) {

	fName, err := fileName(p)

	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	fPath := filepath.Join(s.basePath, p.UserName, fName)

	if err := os.Remove(fPath); err != nil {
		msg := fmt.Sprintf("can't remove file: %s", fPath)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fName, err := fileName(p)

	if err != nil {
		return false, e.Wrap("can't check if file exist", err)
	}

	fPath := filepath.Join(s.basePath, p.UserName, fName)
	switch _, err := os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exist", fPath)
		return false, e.Wrap(msg, err)

	}

	return true, nil
}

func (s Storage) isExistDir(path string) (bool, error) {
	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if directory %s exist", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	defer func() { _ = file.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
