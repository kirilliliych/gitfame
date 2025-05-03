//go:build !solution

package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"gitlab.com/slon/shad-go/gitfame/internal/gitrepository"
	"gitlab.com/slon/shad-go/gitfame/internal/gitstatistics"
	"gitlab.com/slon/shad-go/gitfame/internal/printer"
	"gitlab.com/slon/shad-go/gitfame/internal/utility"
	"os"
)

func main() {
	repository := pflag.String("repository", ".", "путь до Git репозитория")
	revision := pflag.String("revision", "HEAD", "указатель на коммит")
	orderBy := pflag.String("order-by", "lines", "ключ сортировки результатов")
	useCommitter := pflag.Bool("use-committer", false, "заменить в расчетах автора на коммиттера")
	format := pflag.String("format", "tabular", "формат вывода")
	var extensions []string
	var languages []string
	var exclude []string
	var restrictTo []string
	pflag.StringSliceVar(&extensions, "extensions", []string{}, "список расширений файлов в расчете")
	pflag.StringSliceVar(&languages, "languages", []string{}, "список языков(программирования, разметки и др.)")
	pflag.StringSliceVar(&exclude, "exclude", []string{}, "набор Glob паттернов, исключающих файлы из расчета")
	pflag.StringSliceVar(&restrictTo, "restrict-to", []string{}, "набор Glob паттернов, исключающий не подходящее под них")
	pflag.Parse()

	handler := gitstatistics.Handler{Repository: gitrepository.Handler(*repository), Revision: *revision, UseCommitter: *useCommitter, Extensions: extensions, Languages: languages, Exclude: exclude, RestrictTo: restrictTo}
	statisticsByUser, err := handler.GetStatistics()
	if err != nil {
		fmt.Println("error while gathering statistics: ", err)
		os.Exit(1337)
	}

	printer, err := printer.Select(*format)
	if err != nil {
		fmt.Println("error while creating printer", err)
		os.Exit(1337)
	}
	statisticsByUser = utility.SortUserInfos(statisticsByUser, *orderBy)
	err = printer.Print(os.Stdout, &statisticsByUser)
	if err != nil {
		fmt.Println("", err)
		os.Exit(1337)
	}
}
