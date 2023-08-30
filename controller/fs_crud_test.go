package controller

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hawks-atlanta/fs-prototype/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateFile_Check(t *testing.T) {
	t.Run("Succeed", func(t *testing.T) {
		assertions := assert.New(t)

		var (
			contents = "fmt.Println(`hello`)"
			cf       = CreateFile{
				Filename:  "hello-world.go",
				OwnerUUID: uuid.New(),
				Hash:      utils.Hash(contents),
				Size:      uint64(len(contents)),
			}
		)
		assertions.Nil(cf.Check())
	})
	t.Run("Invalid", func(t *testing.T) {
		t.Run("Filename", func(t *testing.T) {
			assertions := assert.New(t)
			var (
				contents = "fmt.Println(`hello`)"
				cf       = CreateFile{
					Filename:  "!hello-world.go",
					OwnerUUID: uuid.New(),
					Hash:      utils.Hash(contents),
					Size:      uint64(len(contents)),
				}
			)
			assertions.NotNil(cf.Check())
		})
		t.Run("OwnerUUID", func(t *testing.T) {
			assertions := assert.New(t)
			var (
				contents = "fmt.Println(`hello`)"
				cf       = CreateFile{
					Filename: "hello-world.go",
					Hash:     utils.Hash(contents),
					Size:     uint64(len(contents)),
				}
			)
			assertions.NotNil(cf.Check())
		})
		t.Run("Size", func(t *testing.T) {
			assertions := assert.New(t)

			var (
				contents = "fmt.Println(`hello`)"
				cf       = CreateFile{
					Filename:  "hello-world.go",
					OwnerUUID: uuid.New(),
					Hash:      utils.Hash(contents),
					Size:      0,
				}
			)
			assertions.NotNil(cf.Check())
		})
	})
}

func TestController_CreateFile(t *testing.T) {
	t.Run("Create files", func(t *testing.T) {
		t.Run("Without parent", func(t *testing.T) {
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

			assertions.NotEqual(uuid.Nil, file.UUID)
			assertions.Equal((*uuid.UUID)(nil), file.ParentUUID)
			assertions.Equal(cf.OwnerUUID, file.OwnerUUID)
			assertions.Equal(cf.Filename, file.Name)
		})
		t.Run("With valid parent", func(t *testing.T) {
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

			assertions.NotEqual(uuid.Nil, file.UUID)
			assertions.Equal(parent.UUID, *file.ParentUUID)
			assertions.Equal(cf.OwnerUUID, file.OwnerUUID)
			assertions.Equal(cf.Filename, file.Name)
		})
		t.Run("Create file in directory of other user", func(t *testing.T) {
			assertions := assert.New(t)

			c, err := Default()
			assertions.Nil(err)
			defer c.Close()

			var (
				parentCf = CreateFile{
					Filename:  "Desktop",
					OwnerUUID: uuid.New(),
				}
			)
			parent, err := c.CreateFile(&parentCf)
			assertions.Nil(err)

			var (
				contents = "fmt.Println(`hello`)"
				cf       = CreateFile{
					Filename:        "hello-world.go",
					OwnerUUID:       uuid.New(),
					Hash:            utils.Hash(contents),
					ParentDirectory: &parent.UUID,
					Size:            uint64(len(contents)),
				}
			)
			_, err = c.CreateFile(&cf)
			assertions.NotNil(err)
		})
		t.Run("Invalid request", func(t *testing.T) {
			assertions := assert.New(t)

			c, err := Default()
			assertions.Nil(err)
			defer c.Close()

			var (
				contents = "fmt.Println(`hello`)"
				cf       = CreateFile{
					Filename:  "!hello-world.go",
					OwnerUUID: uuid.New(),
					Hash:      utils.Hash(contents),
					Size:      uint64(len(contents)),
				}
			)
			_, err = c.CreateFile(&cf)
			assertions.NotNil(err)
		})
	})
}
