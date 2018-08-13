package database

import (
	"os"
	"sync"
	"time"

	"github.com/go-ini/ini"
	"github.com/gwaylib/errors"
)

var (
	cacheLock = sync.Mutex{}
	cache     = map[string]*DB{}
)

func regCache(iniFileName, sectionName string, db *DB) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	key := iniFileName + sectionName
	cache[key] = db
}

func getCache(iniFileName, sectionName string) (*DB, error) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	key := iniFileName + sectionName

	// get from cache
	db, ok := cache[key]
	if ok {
		return db, nil
	}

	// create a new
	cfg, err := ini.Load(iniFileName)
	if err != nil {
		return nil, errors.As(err, iniFileName)
	}
	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return nil, errors.As(err, sectionName)
	}
	drvName, err := section.GetKey("driver")
	if err != nil {
		return nil, errors.As(err, "not found 'driver'")
	}
	dsn, err := section.GetKey("dsn")
	if err != nil {
		return nil, errors.As(err, "not found 'dsn'")
	}
	// http://techblog.en.klab-blogs.com/archives/31093990.html
	lifeTimeKey, err := section.GetKey("life_time")
	if err != nil {
		// ignore error, and make default value
		lifeTimeKey, err = section.NewKey("life_time", "0")
		if err != nil {
			return nil, errors.As(err)
		}
	}
	lifeTime, err := lifeTimeKey.Int64()
	if err != nil {
		return nil, errors.As(err, "error life_time value")
	}

	db, err = Open(drvName.String(), os.ExpandEnv(dsn.String()))
	if err != nil {
		return nil, errors.As(err)
	}
	if lifeTime > 0 {
		db.SetConnMaxLifetime(time.Duration(lifeTime) * time.Second)
	}
	cache[key] = db
	return db, nil
}

func closeCache() {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	for key, db := range cache {
		Close(db)
		delete(cache, key)
	}
}
