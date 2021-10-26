package repo_mysql

import (
	"log"
	"sync"
	"time"

	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLFavoriteRepository struct {
	repository.IFavoriteRepository
	db       db.MySQLDB
	asyncFav AsyncAddFavorite
}

func NewMySQLFavoriteRepository(DB db.MySQLDB) *MySQLFavoriteRepository {
	return &MySQLFavoriteRepository{
		db:       DB,
		asyncFav: *NewAsyncAddFavorite(DB),
	}
}

func (r *MySQLFavoriteRepository) FindFavorite(TwID model.TwIDType, UsrID model.UsrIDType) (ret model.Favorite, err error) {
	err = r.db.DB.Get(&ret, "SELECT * FROM Favorite WHERE tw_id=? and usr_id=?", TwID, UsrID)
	return ret, db.InterpretMySQLError(err)
}
func (r *MySQLFavoriteRepository) FindFavorites(TwID model.TwIDType) (favs []model.Favorite, err error) {
	favs = []model.Favorite{}
	err = r.db.DB.Select(&favs, "SELECT * FROM Favorite WHERE tw_id=?", TwID)
	return favs, db.InterpretMySQLError(err)
}

type AsyncAddFavorite struct {
	addedFavorites []model.Favorite
	favChan        chan model.Favorite
	db             db.MySQLDB
	mutex          sync.Mutex
}

func NewAsyncAddFavorite(db db.MySQLDB) *AsyncAddFavorite {
	ret := AsyncAddFavorite{
		addedFavorites: []model.Favorite{},
		favChan:        make(chan model.Favorite, 1000000),
		mutex:          sync.Mutex{},
		db:             db,
	}
	go func() {
		for {
			fav := <-ret.favChan
			ret.mutex.Lock()
			ret.addedFavorites = append(ret.addedFavorites, fav)
			ret.mutex.Unlock()
		}
	}()
	time.AfterFunc(1*time.Second, ret.CollectFavorites)
	return &ret
}

func (f *AsyncAddFavorite) CollectFavorites() {
	f.mutex.Lock()
	buf := f.addedFavorites
	f.addedFavorites = []model.Favorite{}
	f.mutex.Unlock()
	if len(buf) > 0 {
		_, err := f.db.DB.NamedExec("INSERT IGNORE INTO Favorite (tw_id, usr_id) VALUES (:tw_id, :usr_id)", buf)
		if err != nil {
			log.Print("Failed to insert ")
		} else {
			log.Printf("%d favorites inserted", len(buf))
		}
	}

	time.AfterFunc(1*time.Second, f.CollectFavorites)
}

// ignores duplicate key
func (r *MySQLFavoriteRepository) AddFavorite(fav model.Favorite) (model.Favorite, error) {
	// _, err := r.db.DB.Exec("INSERT IGNORE INTO Favorite (tw_id, usr_id) VALUES (?,?)", fav.TwID, fav.UsrID)
	r.asyncFav.favChan <- fav
	return fav, nil
}

// ignores not found key
func (r *MySQLFavoriteRepository) DeleteFavorite(TwID model.TwIDType, UsrID model.UsrIDType) error {
	_, err := r.db.DB.Exec("DELETE IGNORE FROM Favorite where tw_id=? and usr_id=?", TwID, UsrID)
	return err
}
