package internal

import (
	"go.uber.org/zap"
	"os"
)

type Program struct {
	FilePath string
}

func (p *Program) Load(dest *[]byte, offset uint) {
	file, err := os.Open(p.FilePath)

	if err != nil {
		panic(err)
	}

	read, err := file.Read((*dest)[offset:])

	if err != nil {
		panic(err)
	}

	Logger.Info("loaded program", zap.Int("programBytes", read))
}
