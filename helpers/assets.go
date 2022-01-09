package helper

import (
	bin "celme/bindata/blank"
	assetfs "github.com/elazarl/go-bindata-assetfs"
)

// Assets файловая система с ресурсами приложения
var Assets = &assetfs.AssetFS{
	Asset:     bin.Asset,
	AssetDir:  bin.AssetDir,
	AssetInfo: bin.AssetInfo,
	Prefix:    "",
}
