package badgerkv

//https://dgraph.io/docs/badger/get-started/
import (
	"encoding/json"
	"github.com/dgraph-io/badger/v2"
	"log"
)

type MigratedCsvFile struct {
	FileName string `json:"file_name"`
	UserID   int    `json:"user_id"`
	Known    bool   `json:"known_category"`
	Status   string `json:"status"`
	FilePath string `json:"file_path"`
	Index    string `json:"index"`
}

type BadgerStorage struct {
	VolatileDataStore *badger.DB
	HybridDataStore   *badger.DB
}

func Init() *BadgerStorage {
	tempDB, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		log.Printf("error init badger .... : %+v", err)
	}
	hbDB, err := badger.Open(badger.DefaultOptions("storage/persistence"))
	if err != nil {
		log.Printf("error init badger .... : %+v", err)
	}
	return &BadgerStorage{
		VolatileDataStore: tempDB,
		HybridDataStore:   hbDB,
	}

}

func (b *BadgerStorage) AddNewItem(index string, data MigratedCsvFile) {
	inputData, err := json.Marshal(data)
	if err != nil {
		log.Printf("[AddNewItem] error marshall input data : %+v", err)
	}

	err = b.HybridDataStore.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(index), inputData)
	})
	if err != nil {
		log.Printf("[AddNewItem] error save data to HybridDataStore kv : %+v", err)
	}

	err = b.VolatileDataStore.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(index), inputData)
	})
	if err != nil {
		log.Printf("[AddNewItem] error save data to VolatileDataStore kv : %+v", err)
	}
}

func (b *BadgerStorage) UpdateItem(data MigratedCsvFile) {
	inputData, err := json.Marshal(data)
	if err != nil {
		log.Printf("[AddNewItem] error marshall input data : %+v", err)
	}

	err = b.HybridDataStore.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(data.Index), inputData)
	})
	if err != nil {
		log.Printf("[AddNewItem] error save data to HybridDataStore kv : %+v", err)
	}

	err = b.VolatileDataStore.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(data.Index), inputData)
	})
	if err != nil {
		log.Printf("[AddNewItem] error save data to VolatileDataStore kv : %+v", err)
	}
}

func (b *BadgerStorage) GetVolatileItemList() (data []*MigratedCsvFile) {
	//volatile data
	if err := b.VolatileDataStore.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			err := it.Item().Value(func(v []byte) error {
				csvData := MigratedCsvFile{}
				err := json.Unmarshal(v, &csvData)
				if err != nil {
					log.Printf("[GetVolatileItemList] error retrieve data : %+v", err)
				}

				data = append(data, &csvData)
				return nil
			})
			if err != nil {
				log.Printf("[GetVolatileItemList] error retrieve data : %+v", err)
			}
		}
		return nil
	}); err != nil {
		log.Printf("[GetVolatileItemList] error retrieve data : %+v", err)
	}
	return
}

func (b *BadgerStorage) GetPersistenceItemList() (data []*MigratedCsvFile) {
	//volatile data
	if err := b.HybridDataStore.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			err := it.Item().Value(func(v []byte) error {
				csvData := MigratedCsvFile{}
				err := json.Unmarshal(v, &csvData)
				if err != nil {
					log.Printf("[GetVolatileItemList] error retrieve data : %+v", err)
				}
				data = append(data, &csvData)
				return nil
			})
			if err != nil {
				log.Printf("[GetVolatileItemList] error retrieve data : %+v", err)
			}
		}
		return nil
	}); err != nil {
		log.Printf("[GetVolatileItemList] error retrieve data : %+v", err)
	}
	return
}
