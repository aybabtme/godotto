package do

import (
	"sync"

	"github.com/digitalocean/godo"
)

type database struct {
	droplets *dropletsDB
}

func newDatabase() *database {
	return &database{
		droplets: newDropletsDB(),
	}
}

func (db *database) Droplets() *dropletsDB { return db.droplets }

type dropletsDB struct {
	db *database

	mu       sync.Mutex
	droplets map[int]*godo.Droplet
}

func newDropletsDB() *dropletsDB {
	return &dropletsDB{
		droplets: make(map[int]*godo.Droplet),
	}
}

func (db *dropletsDB) Create(name, region, size, image string) (*godo.Droplet, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	d := &godo.Droplet{
		ID:       len(db.droplets) + 1,
		Name:     name,
		Region:   &godo.Region{Slug: region},
		Size:     &godo.Size{Slug: size},
		SizeSlug: size,
		Image:    &godo.Image{Slug: image},
	}
	db.droplets[d.ID] = d
	return d, nil
}

func (db *dropletsDB) Delete(id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.droplets, id)
	return nil
}
