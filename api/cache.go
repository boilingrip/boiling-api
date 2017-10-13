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

func (s *SyncedLookupTable) MustLookUp(ss string) int {
	i, err := s.LookUp(ss)
	if err != nil {
		panic(fmt.Sprintf("Key not found, cache out of sync? %s", err.Error()))
	}
	return i
}

func (s *SyncedLookupTable) LookUp(ss string) (int, error) {
	s.RLock()
	i, err := s.l.LookUp(ss)
	s.RUnlock()
	return i, err
}

func (s *SyncedLookupTable) MustReverseLookUp(i int) string {
	ss, err := s.ReverseLookUp(i)
	if err != nil {
		panic(fmt.Sprintf("Key not found, cache out of sync? %s", err.Error()))
	}
	return ss
}

func (s *SyncedLookupTable) ReverseLookUp(i int) (string, error) {
	s.RLock()
	ss, err := s.l.ReverseLookUp(i)
	s.RUnlock()
	return ss, err
}

func (s *SyncedLookupTable) Has(ss string) bool {
	s.RLock()
	has := s.l.Has(ss)
	s.RUnlock()
	return has
}

func (s *SyncedLookupTable) HasReverse(i int) bool {
	s.RLock()
	has := s.l.HasReverse(i)
	s.RUnlock()
	return has
}

func (s *SyncedLookupTable) Keys() []string {
	s.RLock()
	keys := s.l.Keys()
	s.RUnlock()
	return keys
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
