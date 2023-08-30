package controller

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hawks-atlanta/fs-prototype/utils"
	"github.com/stretchr/testify/assert"
)

func TestController_CanReadFile(t *testing.T) {
	t.Run("Owned file", func(t *testing.T) {
		assertions := assert.New(t)

		c, err := Default()
		assertions.Nil(err)
		defer c.Close()

		var (
			contents = "fmt.Println(`hello`)"
			cf       = CreateFile{
				Filename:  "hello-world.go",
				OwnerUUID: uuid.New(),
				Hash:      utils.Hash(contents),
				Size:      uint64(len(contents)),
			}
		)
		file, err := c.CreateFile(&cf)
		assertions.Nil(err)

		var crf = CanReadFile{
			UserUUID: cf.OwnerUUID,
			FileUUID: file.UUID,
		}
		err = c.CanReadFile(&crf)
		assertions.Nil(err)
	})
	t.Run("Not owned but shared", func(t *testing.T) {
		assertions := assert.New(t)

		c, err := Default()
		assertions.Nil(err)
		defer c.Close()

		var (
			contents = "fmt.Println(`hello`)"
			cf       = CreateFile{
				Filename:  "hello-world.go",
				OwnerUUID: uuid.New(),
				Hash:      utils.Hash(contents),
				Size:      uint64(len(contents)),
			}
		)
		file, err := c.CreateFile(&cf)
		assertions.Nil(err)

		var sr = ShareRequest{
			OwnerUUID:      cf.OwnerUUID,
			FileUUID:       file.UUID,
			TargetUserUUID: uuid.New(),
		}
		err = c.ShareFile(&sr)
		assertions.Nil(err)

		var crf = CanReadFile{
			UserUUID: sr.TargetUserUUID,
			FileUUID: file.UUID,
		}
		err = c.CanReadFile(&crf)
		assertions.Nil(err)
	})
	t.Run("Not owned but shared inside directory", func(t *testing.T) {
		assertions := assert.New(t)

		c, err := Default()
		assertions.Nil(err)
		defer c.Close()

		var (
			owner    = uuid.New()
			parentCf = CreateFile{
				Filename:  "Desktop",
				OwnerUUID: owner,
			}
		)
		parent, err := c.CreateFile(&parentCf)
		assertions.Nil(err)

		var (
			contents = "fmt.Println(`hello`)"
			cf       = CreateFile{
				Filename:        "hello-world.go",
				OwnerUUID:       owner,
				Hash:            utils.Hash(contents),
				ParentDirectory: &parent.UUID,
				Size:            uint64(len(contents)),
			}
		)
		file, err := c.CreateFile(&cf)
		assertions.Nil(err)

		// Share parent
		var sr = ShareRequest{
			OwnerUUID:      parentCf.OwnerUUID,
			FileUUID:       parent.UUID,
			TargetUserUUID: uuid.New(),
		}
		err = c.ShareFile(&sr)
		assertions.Nil(err)

		// Check if can read file inside shared parent
		var crf = CanReadFile{
			UserUUID: sr.TargetUserUUID,
			FileUUID: file.UUID,
		}
		err = c.CanReadFile(&crf)
		assertions.Nil(err)
	})
	t.Run("Zero Access", func(t *testing.T) {
		assertions := assert.New(t)

		c, err := Default()
		assertions.Nil(err)
		defer c.Close()

		var (
			owner    = uuid.New()
			parentCf = CreateFile{
				Filename:  "Desktop",
				OwnerUUID: owner,
			}
		)
		parent, err := c.CreateFile(&parentCf)
		assertions.Nil(err)

		var (
			contents = "fmt.Println(`hello`)"
			cf       = CreateFile{
				Filename:        "hello-world.go",
				OwnerUUID:       owner,
				Hash:            utils.Hash(contents),
				ParentDirectory: &parent.UUID,
				Size:            uint64(len(contents)),
			}
		)
		file, err := c.CreateFile(&cf)
		assertions.Nil(err)

		// Check if can read file inside shared parent
		var crf = CanReadFile{
			UserUUID: uuid.New(),
			FileUUID: file.UUID,
		}
		err = c.CanReadFile(&crf)
		assertions.NotNil(err)
	})
}
