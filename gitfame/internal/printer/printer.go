package printer

import (
	"gitlab.com/slon/shad-go/gitfame/internal/gitstatistics"
	"io"
)

type Printer interface {
	Print(writer io.Writer, statistics *[]gitstatistics.UserStatistics) error
}
