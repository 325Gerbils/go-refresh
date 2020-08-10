package refresh

import (
	"html/template"
	"log"

	"github.com/fsnotify/fsnotify"
)

// Template accepts a list of paths to html template files
// It parses the templates and returns a *template.Template object for each one.
// It sets up a goroutine that watches for file changes in each and refreshes the
// templates when they change
func Template(paths ...string) (templates []*template.Template, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	for i, path := range paths {
		templates[i], err = template.ParseFiles(path)
	}
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("Event:", event)
				for i, path := range paths {
					templates[i], err = template.ParseFiles(path)
				}
				if err != nil {
					log.Fatal(err)
				}
				log.Println("Updated templates")
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Fatal(err)
			}
		}
	}()

	for _, path := range paths {
		err = watcher.Add(path)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
