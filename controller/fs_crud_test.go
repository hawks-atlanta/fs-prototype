package controller

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hawks-atlanta/fs-prototype/models"
	"github.com/hawks-atlanta/fs-prototype/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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
	})
}

func TestController_DeleteFile(t *testing.T) {
	t.Run("Delete file", func(t *testing.T) {
		assertions := assert.New(t)

		c, err := Default()
		assertions.Nil(err)
		defer c.Close()

		// Create file
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

		// Delete file
		var df = DeleteFile{
			OwnerUUID: cf.OwnerUUID,
			FileUUID:  file.UUID,
		}
		err = c.DeleteFile(&df)
		assertions.Nil(err)

		// Verify it doesn't exists anymore
		var check models.File
		err = c.DB.
			Where("uuid = ?", file.UUID).
			First(&check).
			Error
		assertions.ErrorIs(err, gorm.ErrRecordNotFound)
	})
	t.Run("Delete directory", func(t *testing.T) {
		assertions := assert.New(t)

		c, err := Default()
		assertions.Nil(err)
		defer c.Close()

		// Create parent
		var (
			owner    = uuid.New()
			parentCf = CreateFile{
				Filename:  "Desktop",
				OwnerUUID: owner,
			}
		)
		parent, err := c.CreateFile(&parentCf)
		assertions.Nil(err)

		// Create file
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

		// Delete parent
		var df = DeleteFile{
			OwnerUUID: cf.OwnerUUID,
			FileUUID:  parent.UUID,
		}
		err = c.DeleteFile(&df)
		assertions.Nil(err)

		// Verify parent doesn't exists anymore
		var check models.File
		err = c.DB.
			Where("uuid = ?", parent.UUID).
			First(&check).
			Error
		assertions.ErrorIs(err, gorm.ErrRecordNotFound)

		// Verify child doesn't exists anymore
		check = models.File{}
		err = c.DB.
			Where("uuid = ?", file.UUID).
			First(&check).
			Error
		assertions.ErrorIs(err, gorm.ErrRecordNotFound)
	})
}

func TestController_QueryFile(t *testing.T) {
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

		var qf = QueryFile{
			UserUUID: cf.OwnerUUID,
			FileUUID: file.UUID,
		}
		archive, err := c.QueryFile(&qf)
		assertions.Nil(err)

		assertions.NotEmpty(archive.Hash)
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

		var qf = QueryFile{
			UserUUID: sr.TargetUserUUID,
			FileUUID: file.UUID,
		}
		archive, err := c.QueryFile(&qf)
		assertions.Nil(err)

		assertions.NotEmpty(archive.Hash)
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
		var qf = QueryFile{
			UserUUID: sr.TargetUserUUID,
			FileUUID: file.UUID,
		}
		archive, err := c.QueryFile(&qf)
		assertions.Nil(err)

		assertions.NotEmpty(archive.Hash)
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
		var qf = QueryFile{
			UserUUID: uuid.New(),
			FileUUID: file.UUID,
		}
		archive, err := c.QueryFile(&qf)
		assertions.NotNil(err)

		assertions.Empty(archive.Hash)
	})
}

func TestController_MoveFile(t *testing.T) {
	// TODO: implement me!
}
