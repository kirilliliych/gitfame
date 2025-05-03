package printer

import (
	"encoding/csv"
	"gitlab.com/slon/shad-go/gitfame/internal/gitstatistics"
	"io"
	"strconv"
)

type csvprinter int

func (p *csvprinter) Print(writer io.Writer, statistics *[]gitstatistics.UserStatistics) error {
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write([]string{"Name", "Lines", "Commits", "Files"}); err != nil {
		return err
	}
	for _, curStatistics := range *statistics {
		if curStatistics.Files == 0 {
			continue
		}
		row := []string{
			curStatistics.Name,
			strconv.Itoa(curStatistics.Lines),
			strconv.Itoa(curStatistics.Commits),
			strconv.Itoa(curStatistics.Files),
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}
	csvWriter.Flush()

	return nil
}
