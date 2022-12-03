package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pot-code/web-cli/pkg/command"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type MassRenameConfig struct {
	Suffix []string `flag:"suffix" alias:"s" validate:"gt=0,required,ascii,dive" usage:"included file suffix"`
	Length int      `flag:"length" alias:"l" validate:"gt=0,lt=33" usage:"maximum length of the name"`
	Dry    bool     `flag:"dry" alias:"d" usage:"dry run"`
	Dir    string   `arg:"0" alias:"DIR" validate:"required"`
}

var MassRenameCmd = command.NewCliCommand("rename", "mass rename files to md5 string",
	&MassRenameConfig{
		Length: 32,
	},
	command.WithArgUsage("DIR"),
	command.WithAlias([]string{"re"}),
).AddHandlers(NewMassName()).BuildCommand()

type MassRename struct {
	suffixMap      map[string]struct{}
	hashCounterMap map[string]int
}

func NewMassName() *MassRename {
	return &MassRename{
		suffixMap:      make(map[string]struct{}),
		hashCounterMap: make(map[string]int),
	}
}

func (mr *MassRename) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*MassRenameConfig)
	root := config.Dir
	if !fileExists(root) {
		return errors.New("DIR not exists")
	}

	mr.initSuffixMap(config.Suffix)
	return filepath.WalkDir(root, func(oldpath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		filename := d.Name()
		if err != nil {
			log.WithFields(log.Fields{"dir": filename, "error": err}).Error("failed to walk into directory")
		}
		if !mr.matchAnySuffix(filename) {
			return nil
		}

		newname, err := mr.getHashedName(oldpath, config.Length)
		if err != nil {
			return err
		}

		dir := path.Dir(oldpath)
		newpath := path.Join(dir, newname)
		log.WithFields(log.Fields{"oldpath": oldpath, "newpath": newpath}).Info("renaming file")
		if !config.Dry {
			err = os.Rename(oldpath, newpath)
			if err != nil {
				return fmt.Errorf("rename file %s: %w", oldpath, err)
			}
		}
		return nil
	})
}

func (mr *MassRename) initSuffixMap(suffix []string) {
	sm := mr.suffixMap
	for _, s := range suffix {
		if !strings.HasPrefix(s, ".") {
			s = "." + s
		}
		log.WithFields(log.Fields{"suffix": s}).Debug("registered suffix")
		sm[s] = struct{}{}
	}
}

func (mr *MassRename) getHashedName(fp string, length int) (string, error) {
	hash, err := mr.hashFile(fp)
	if err != nil {
		return "", err
	}
	hash = hash[:length]

	name := hash
	ext := filepath.Ext(fp)
	if mr.inHashCounterMap(hash) {
		name = mr.getNextCollisionName(hash)
	}
	mr.addToHashCounterMap(hash)
	return name + ext, nil
}

func (mr *MassRename) addToHashCounterMap(hash string) {
	mr.hashCounterMap[hash]++
}

func (mr *MassRename) inHashCounterMap(hash string) bool {
	_, ok := mr.hashCounterMap[hash]
	return ok
}

func (mr *MassRename) getNextCollisionName(hash string) string {
	hcm := mr.hashCounterMap
	index := hcm[hash]
	return hash + "-" + strconv.Itoa(index)
}

func (mr *MassRename) matchAnySuffix(filename string) bool {
	ext := filepath.Ext(filename)
	_, ok := mr.suffixMap[ext]
	return ok
}

func (mr *MassRename) hashFile(fp string) (string, error) {
	fd, err := os.Open(fp)
	if err != nil {
		panic(fmt.Errorf("failed to hash file '%s', err: %v", fp, err))
	}
	defer fd.Close()

	h := md5.New()
	w, err := io.Copy(h, fd)
	if err != nil {
		panic(fmt.Errorf("failed to hash file '%s', err: %v", fp, err))
	}
	log.WithFields(log.Fields{"file": fp, "write": w}).Debug("copy file data to hash")
	return hex.EncodeToString(h.Sum(nil)), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
