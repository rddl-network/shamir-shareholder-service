package service

import "github.com/syndtr/goleveldb/leveldb"

var (
	key = []byte("mnemonic")
)

func InitDB(dbPath string) (db *leveldb.DB, err error) {
	return leveldb.OpenFile(dbPath, nil)
}

func (ss *ShamirService) PutMnemonic(mnemonic []byte) (err error) {
	return ss.db.Put(key, mnemonic, nil)
}

func (ss *ShamirService) GetMnemonic() (mnemonic []byte, err error) {
	return ss.db.Get(key, nil)
}
