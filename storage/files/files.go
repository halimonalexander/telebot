package files

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"telebot/lib/e"
	"telebot/storage"
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {}()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, 0774); err != nil {
		return e.WrapError("Unable to create folder", err)
	}

	fileName, err := filename(page)
	if err != nil {
		return err
	}

	fullName := filepath.Join(filePath, fileName)
	file, err := os.Create(fullName)
	if err != nil {
		return e.WrapError("Unable to create file", err)
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.WrapError("encoding error", err)
	}

	return nil
}

func (s Storage) PickRandom(username string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, username)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, e.WrapError("Unable to read dir", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrorNoSavedPages
	}

	n := rand.Intn(len(files))

	return s.decodePage(filepath.Join(path, files[n].Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fileName, err := filename(page)
	if err != nil {
		return e.WrapError("Unable to remove page file", err)
	}

	filePath := filepath.Join(s.basePath, page.UserName, fileName)
	if err := os.Remove(filePath); err != nil {
		return e.WrapError(fmt.Sprintf("Unable to delete file %s", filePath), err)
	}

	return nil
}

func (s Storage) Prune(username string) error {
	filePath := filepath.Join(s.basePath, username)
	if err := os.Remove(filePath); err != nil {
		return e.WrapError(fmt.Sprintf("Unable to delete folder %s", filePath), err)
	}

	return nil
}

func (s Storage) Exists(page *storage.Page) (bool, error) {
	fileName, err := filename(page)
	if err != nil {
		return false, e.WrapError("Unable to locate file", err)
	}

	filePath := filepath.Join(s.basePath, page.UserName, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, e.WrapError("Unable to open file", err)
	} else {
		return true, nil
	}
}

func (s Storage) decodePage(filename string) (*storage.Page, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, e.WrapError("Unable to open file", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.WrapError("Unable to decode file", err)
	}

	return &p, nil
}

func filename(p *storage.Page) (string, error) {
	hash, err := p.Hash()
	if err != nil {
		return "", err
	}

	return hash + ".txt", nil
}
