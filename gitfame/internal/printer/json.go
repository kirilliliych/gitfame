package printer

import (
	"encoding/json"
	"gitlab.com/slon/shad-go/gitfame/internal/gitstatistics"
	"io"
)

type jsonprinter int

func (p *jsonprinter) Print(writer io.Writer, statistics *[]gitstatistics.UserStatistics) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(*statistics)
}
