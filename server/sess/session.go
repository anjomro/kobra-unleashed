package sess

import (
	"github.com/anjomro/kobra-unleashed/server/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func SetupSessionStore() {
	// Create a new session store
	store = session.New(session.Config{
		Storage: storage.New(storage.Config{
			Folder: "./.sessions",
		}),
	})
}

func GetSession(c *fiber.Ctx) *session.Session {
	// Get session from store
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	return sess
}
