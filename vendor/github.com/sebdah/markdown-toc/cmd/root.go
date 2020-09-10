package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sebdah/markdown-toc/toc"
	"github.com/spf13/cobra"
)

var (
	// depth indicates how many levels of headers should be included in the ToC.
	//
	// If 0, all headers are included.
	depth int

	// header is the string injected as a header for the table of contents.
	header string

	// inline indicates whether we should do an inline replacement of the ToC or
	// if we should print to stdout. Default is to print to stdout.
	inline bool

	// noHeader is indicating whether or not a header should be injected in the
	// table of contents.
	noHeader bool

	// skipHeaders is indicating how many of headers we should skip from the
	// top. This is useful for projects that has e.g. the project name as the
	// first header, but they don't want that to go in the ToC.
	//
	// The flag is ignoring the header size (H1, H2, etc).
	//
	// If this is set to 0 no headers would be skipped.
	skipHeaders = 0

	// replaceToC is indicating whether we should replace the table of contents
	// in the input file. This assumes that there are two tags indicating where
	// the ToC starts and where it ends:
	//
	// Start:   <!-- ToC start -->
	// End:     <!-- ToC end -->
	//
	// If these tags are not found, the table of contents will be injected on
	// top all existing content in the markdown file.
	replaceToC bool
)

var RootCmd = &cobra.Command{
	Use:   "markdown-toc <file>",
	Short: "Generate a table of contents for your markdown file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}

		t, err := toc.Build(d, header, depth, skipHeaders, !noHeader)
		if err != nil {
			return err
		}

		if replaceToC {
			t = toc.Replace(d, t)
		}

		if inline {
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}

			i, err := f.Stat()
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(args[0], []byte(fmt.Sprintf("%s\n", strings.Join(t, "\n"))), i.Mode())
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("%s\n", strings.Join(t, "\n"))
		}

		return nil
	},
}

func init() {
	RootCmd.Flags().IntVar(&depth, "depth", 0, "Depth of headers to include. Set to 0 for all headers")
	RootCmd.Flags().StringVar(&header, "header", "# Table of Contents", "Text to use for the header for the ToC")
	RootCmd.Flags().BoolVar(&noHeader, "no-header", false, "If this is set there will be no header for the ToC")
	RootCmd.Flags().IntVar(&skipHeaders, "skip-headers", 0, "Number of headers to skip. Useful if you don't want e.g. the first header to be included in the ToC")
	RootCmd.Flags().BoolVar(&replaceToC, "replace", false, "If the replace flag is set the full markdown will be returned and any existing ToC replaced")
	RootCmd.Flags().BoolVar(&inline, "inline", false, "Overwrite the input file with the output from this command. Should be used together with --replace")
}
