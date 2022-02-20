package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/meedoed/telegram-bot/internal/lib/e"
	"github.com/meedoed/telegram-bot/internal/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return e.Wrap("can't save", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return e.Wrap("can't save", err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return e.Wrap("can't save", err)
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap("can't save", err)
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, e.Wrap("can't pick random page", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPage
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't chek if file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't chek if file %s exists", path)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer f.Close()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
