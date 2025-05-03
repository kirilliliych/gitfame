package printer

import (
	"encoding/json"
	"gitlab.com/slon/shad-go/gitfame/internal/gitstatistics"
	"io"
)

type jsonlineprinter int

func (p *jsonlineprinter) Print(writer io.Writer, statistics *[]gitstatistics.UserStatistics) error {
	for _, curStatistics := range *statistics {
		line, err := json.Marshal(curStatistics)
		if err != nil {
			return err
		}
		line = append(line, []byte("\n")...)
		if _, err := writer.Write(line); err != nil {
			return err
		}
	}
	return nil
}
