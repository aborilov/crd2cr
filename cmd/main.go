package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aborilov/crd2cr"
	"github.com/spf13/cobra"
)

type flags struct {
	file string
}

func main() {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "crd2cr",
		Short: "Generate simple CR from CRD",
		RunE: func(cmd *cobra.Command, args []string) error {
			var data []byte
			var err error
			if flags.file == "" {
				data, err = ioutil.ReadAll(bufio.NewReader(os.Stdin))
				if err != nil {
					return err
				}
			} else {
				data, err = ioutil.ReadFile(flags.file)
				if err != nil {
					return err
				}
			}
			res, err := crd2cr.Convert(data)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().StringVar(&flags.file, "file", "", "file path. use STDIN by default")
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
