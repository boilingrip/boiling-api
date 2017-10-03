package api

import (
	"fmt"
	"sync"

	"github.com/boilingrip/boiling-api/db"
)

type SyncedLookupTable struct {
	l LookupTable
	sync.RWMutex
}

type Cache struct {
	formats           *SyncedLookupTable
	leechTypes        *SyncedLookupTable
	media             *SyncedLookupTable
	releaseGroupTypes *SyncedLookupTable
	releaseProperties *SyncedLookupTable
	releaseRoles      *SyncedLookupTable
	privileges        *SyncedLookupTable
}

func NewCache(db db.BoilingDB) (Cache, error) {
	c := Cache{
		formats:           new(SyncedLookupTable),
		leechTypes:        new(SyncedLookupTable),
		media:             new(SyncedLookupTable),
		releaseGroupTypes: new(SyncedLookupTable),
		releaseProperties: new(SyncedLookupTable),
		releaseRoles:      new(SyncedLookupTable),
		privileges:        new(SyncedLookupTable),
	}

	err := c.RefreshFormats(db)
	if err != nil {
		return Cache{}, err
	}

	err = c.RefreshLeechTypes(db)
	if err != nil {
		return Cache{}, err
	}

	err = c.RefreshMedia(db)
	if err != nil {
		return Cache{}, err
	}

	err = c.RefreshReleaseGroupTypes(db)
	if err != nil {
		return Cache{}, err
	}

	err = c.RefreshReleaseProperties(db)
	if err != nil {
		return Cache{}, err
	}

	err = c.RefreshReleaseRoles(db)
	if err != nil {
		return Cache{}, err
	}

	err = c.RefreshPrivileges(db)
	if err != nil {
		return Cache{}, err
	}

	return c, nil
}

func (c Cache) RefreshFormats(db db.BoilingDB) error {
	c.formats.Lock()
	defer c.formats.Unlock()

	formats, err := db.GetAllFormats()
	if err != nil {
		return err
	}

	m := make(map[int]string)
	for k, v := range formats {
		m[k] = fmt.Sprintf("%s$%s", v.Format, v.Encoding)
	}

	t := BuildLookupTable(m)
	c.formats.l = t

	return nil
}

func (c Cache) RefreshLeechTypes(db db.BoilingDB) error {
	c.leechTypes.Lock()
	defer c.leechTypes.Unlock()

	leechTypes, err := db.GetAllLeechTypes()
	if err != nil {
		return err
	}

	t := BuildLookupTable(leechTypes)
	c.leechTypes.l = t

	return nil
}

func (c Cache) RefreshMedia(db db.BoilingDB) error {
	c.media.Lock()
	defer c.media.Unlock()

	media, err := db.GetAllMedia()
	if err != nil {
		return err
	}

	t := BuildLookupTable(media)
	c.media.l = t

	return nil
}

func (c Cache) RefreshReleaseGroupTypes(db db.BoilingDB) error {
	c.releaseGroupTypes.Lock()
	defer c.releaseGroupTypes.Unlock()

	releaseGroupTypes, err := db.GetAllReleaseGroupTypes()
	if err != nil {
		return err
	}

	t := BuildLookupTable(releaseGroupTypes)
	c.releaseGroupTypes.l = t

	return nil
}

func (c Cache) RefreshReleaseProperties(db db.BoilingDB) error {
	c.releaseProperties.Lock()
	defer c.releaseProperties.Unlock()

	releaseProperties, err := db.GetAllReleaseProperties()
	if err != nil {
		return err
	}

	t := BuildLookupTable(releaseProperties)
	c.releaseProperties.l = t

	return nil
}

func (c Cache) RefreshReleaseRoles(db db.BoilingDB) error {
	c.releaseRoles.Lock()
	defer c.releaseRoles.Unlock()

	releaseRoles, err := db.GetAllReleaseGroupRoles()
	if err != nil {
		return err
	}

	t := BuildLookupTable(releaseRoles)
	c.releaseRoles.l = t

	return nil
}

func (c Cache) RefreshPrivileges(db db.BoilingDB) error {
	c.privileges.Lock()
	defer c.privileges.Unlock()

	privileges, err := db.GetAllPrivileges()
	if err != nil {
		return err
	}

	t := BuildLookupTable(privileges)
	c.privileges.l = t

	return nil
}
