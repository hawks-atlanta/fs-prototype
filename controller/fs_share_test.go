package controller

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hawks-atlanta/fs-prototype/models"
	"github.com/hawks-atlanta/fs-prototype/utils"
	"github.com/stretchr/testify/assert"
)

func TestShareWithMe_Check(t *testing.T) {
	t.Run("Succeed", func(t *testing.T) {
		assertions := assert.New(t)

		var swm = ShareWithMe{
			UserUUID: uuid.New(),
		}
		err := swm.Check()
		assertions.Nil(err)
	})
	t.Run("Failed", func(t *testing.T) {
		t.Run("UserUUID", func(t *testing.T) {
			assertions := assert.New(t)

			var swm = ShareWithMe{}
			err := swm.Check()
			assertions.NotNil(err)
		})
	})
}

func TestController_ShareWithMe(t *testing.T) {
	t.Run("Succeed", func(t *testing.T) {
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

		var swm = ShareWithMe{
			UserUUID: sr.TargetUserUUID,
		}
		shared, err := c.ShareWithMe(&swm)
		assertions.Nil(err)

		assertions.Len(shared, 1)
	})
	t.Run("Invalid request", func(t *testing.T) {
		assertions := assert.New(t)

		c, err := Default()
		assertions.Nil(err)
		defer c.Close()

		var swm = ShareWithMe{}
		_, err = c.ShareWithMe(&swm)
		assertions.NotNil(err)
	})
}

func TestShareWithWho_Check(t *testing.T) {
	t.Run("Succeed", func(t *testing.T) {
		assertions := assert.New(t)

		var sww = ShareWithWho{
			OwnerUUID: uuid.New(),
			FileUUID:  uuid.New(),
		}
		err := sww.Check()
		assertions.Nil(err)
	})
	t.Run("Failed", func(t *testing.T) {
		t.Run("OwnerUUID", func(t *testing.T) {
			assertions := assert.New(t)

			var sww = ShareWithWho{
				FileUUID: uuid.New(),
			}
			err := sww.Check()
			assertions.NotNil(err)
		})
		t.Run("OwnerUUID", func(t *testing.T) {
			assertions := assert.New(t)

			var sww = ShareWithWho{
				OwnerUUID: uuid.New(),
			}
			err := sww.Check()
			assertions.NotNil(err)
		})
	})
}

func TestController_ShareWithWho(t *testing.T) {
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

		var sr = ShareRequest{
			OwnerUUID:      cf.OwnerUUID,
			FileUUID:       file.UUID,
			TargetUserUUID: uuid.New(),
		}
		err = c.ShareFile(&sr)
		assertions.Nil(err)

		var sww = ShareWithWho{
			OwnerUUID: cf.OwnerUUID,
			FileUUID:  file.UUID,
		}
		shared, err := c.ShareWithWho(&sww)
		assertions.Nil(err)

		assertions.Len(shared, 1)
	})
	t.Run("Not owned file", func(t *testing.T) {
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

		var sww = ShareWithWho{
			OwnerUUID: uuid.New(),
			FileUUID:  file.UUID,
		}
		shared, err := c.ShareWithWho(&sww)
		assertions.NotNil(err)

		assertions.Len(shared, 0)
	})
	t.Run("Invalid request", func(t *testing.T) {
		assertions := assert.New(t)

		c, err := Default()
		assertions.Nil(err)
		defer c.Close()

		var sww = ShareWithWho{}
		_, err = c.ShareWithWho(&sww)
		assertions.NotNil(err)
	})
}

func TestShareRequest_Check(t *testing.T) {
	t.Run("Succeed", func(t *testing.T) {
		assertions := assert.New(t)
		var sr = ShareRequest{
			OwnerUUID:      uuid.New(),
			FileUUID:       uuid.New(),
			TargetUserUUID: uuid.New(),
		}
		err := sr.Check()
		assertions.Nil(err)
	})
	t.Run("Fail", func(t *testing.T) {
		t.Run("OwnerUUID", func(t *testing.T) {
			assertions := assert.New(t)
			var sr = ShareRequest{
				FileUUID:       uuid.New(),
				TargetUserUUID: uuid.New(),
			}
			err := sr.Check()
			assertions.NotNil(err)
		})
		t.Run("FileUUID", func(t *testing.T) {
			assertions := assert.New(t)
			var sr = ShareRequest{
				OwnerUUID:      uuid.New(),
				TargetUserUUID: uuid.New(),
			}
			err := sr.Check()
			assertions.NotNil(err)
		})
		t.Run("TargetUserUUID", func(t *testing.T) {
			assertions := assert.New(t)
			var sr = ShareRequest{
				OwnerUUID: uuid.New(),
				FileUUID:  uuid.New(),
			}
			err := sr.Check()
			assertions.NotNil(err)
		})
	})
}

func TestController_ShareFile(t *testing.T) {
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

		var sr = ShareRequest{
			OwnerUUID:      cf.OwnerUUID,
			FileUUID:       file.UUID,
			TargetUserUUID: uuid.New(),
		}
		err = c.ShareFile(&sr)
		assertions.Nil(err)

		var check models.SharedFile
		err = c.DB.
			Where("file_uuid = ? AND user_uuid = ?", file.UUID, sr.TargetUserUUID).
			First(&check).
			Error
		assertions.Nil(err)

		assertions.Equal(file.UUID, check.FileUUID)
		assertions.Equal(sr.TargetUserUUID, check.UserUUID)
	})
	t.Run("Not owned file", func(t *testing.T) {
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
			OwnerUUID:      uuid.New(),
			FileUUID:       file.UUID,
			TargetUserUUID: uuid.New(),
		}
		err = c.ShareFile(&sr)
		assertions.NotNil(err)

		var check models.SharedFile
		err = c.DB.
			Where("file_uuid = ? AND user_uuid = ?", file.UUID, sr.TargetUserUUID).
			First(&check).
			Error
		assertions.NotNil(err)
	})
	t.Run("Share twice to same user", func(t *testing.T) {
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

		var check models.SharedFile
		err = c.DB.
			Where("file_uuid = ? AND user_uuid = ?", file.UUID, sr.TargetUserUUID).
			First(&check).
			Error
		assertions.Nil(err)

		assertions.Equal(file.UUID, check.FileUUID)
		assertions.Equal(sr.TargetUserUUID, check.UserUUID)

		err = c.ShareFile(&sr)
		assertions.NotNil(err)
	})
	t.Run("Invalid Share request", func(t *testing.T) {
		assertions := assert.New(t)

		c, err := Default()
		assertions.Nil(err)
		defer c.Close()

		var sr = ShareRequest{}
		err = c.ShareFile(&sr)
		assertions.NotNil(err)
	})
}

func TestController_UnshareFile(t *testing.T) {
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

		var sr = ShareRequest{
			OwnerUUID:      cf.OwnerUUID,
			FileUUID:       file.UUID,
			TargetUserUUID: uuid.New(),
		}
		err = c.ShareFile(&sr)
		assertions.Nil(err)

		var check models.SharedFile
		err = c.DB.
			Where("file_uuid = ? AND user_uuid = ?", file.UUID, sr.TargetUserUUID).
			First(&check).
			Error
		assertions.Nil(err)

		assertions.Equal(file.UUID, check.FileUUID)
		assertions.Equal(sr.TargetUserUUID, check.UserUUID)

		err = c.UnshareFile(&sr)
		assertions.Nil(err)

		// Check share is removed
		check = models.SharedFile{}
		err = c.DB.
			Where("file_uuid = ? AND user_uuid = ?", file.UUID, sr.TargetUserUUID).
			First(&check).
			Error
		assertions.NotNil(err)
	})
	t.Run("Not owned file", func(t *testing.T) {
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

		var check models.SharedFile
		err = c.DB.
			Where("file_uuid = ? AND user_uuid = ?", file.UUID, sr.TargetUserUUID).
			First(&check).
			Error
		assertions.Nil(err)

		assertions.Equal(file.UUID, check.FileUUID)
		assertions.Equal(sr.TargetUserUUID, check.UserUUID)

		sr.OwnerUUID = uuid.New()
		err = c.UnshareFile(&sr)
		assertions.NotNil(err)

		// Check share is removed
		check = models.SharedFile{}
		err = c.DB.
			Where("file_uuid = ? AND user_uuid = ?", file.UUID, sr.TargetUserUUID).
			First(&check).
			Error
		assertions.Nil(err)

		assertions.Equal(file.UUID, check.FileUUID)
		assertions.Equal(sr.TargetUserUUID, check.UserUUID)
	})
	t.Run("Invalid Share request", func(t *testing.T) {
		assertions := assert.New(t)

		c, err := Default()
		assertions.Nil(err)
		defer c.Close()

		var sr = ShareRequest{}
		err = c.UnshareFile(&sr)
		assertions.NotNil(err)
	})
}
