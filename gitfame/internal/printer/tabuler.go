package printer

import (
	"fmt"
	"gitlab.com/slon/shad-go/gitfame/internal/gitstatistics"
	"io"
	"text/tabwriter"
)

type tabularprinter int

func (p *tabularprinter) Print(writer io.Writer, statistics *[]gitstatistics.UserStatistics) error {
	tabularWriter := tabwriter.NewWriter(writer, 1, 1, 1, ' ', 0)
	_, err := fmt.Fprint(tabularWriter, "Name\tLines\tCommits\tFiles\n")
	if err != nil {
		return err
	}
	for _, curStatistics := range *statistics {
		_, err := fmt.Fprintf(tabularWriter, "%s\t%d\t%d\t%d\n", curStatistics.Name, curStatistics.Lines, curStatistics.Commits, curStatistics.Files)
		if err != nil {
			return err
		}
	}
	if err := tabularWriter.Flush(); err != nil {
		return err
	}
	return nil
}
