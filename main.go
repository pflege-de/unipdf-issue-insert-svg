package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	licenseKey := os.Getenv("UNIPDF_LICENSE_KEY")
	customerName := os.Getenv("UNIPDF_CUSTOMER_NAME")

	if err := license.SetLicenseKey(licenseKey, customerName); err != nil {
		panic(fmt.Errorf("cannot set unipdf license: %w", err))
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("syntax: %s <input.pdf>", os.Args[0])
	}

	target := os.Args[1]

	if _, err := os.Stat(target); err != nil {
		return fmt.Errorf("cannot stat target file: %w", err)
	}

	files, err := os.ReadDir("testdata")

	if err != nil {
		return fmt.Errorf("cannot list image files: %w", err)
	}

	for _, ent := range files {
		name := filepath.Join("testdata", ent.Name())

		if filepath.Ext(name) != ".svg" {
			continue
		}

		fmt.Printf("=> test: %s\n", name)

		if err := insert(target, name); err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}

	return nil
}

func insert(target string, name string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("recover: %v\n", err)
		}
	}()

	r, f, err := model.NewPdfReaderFromFile(target, &model.ReaderOpts{})

	if err != nil {
		return fmt.Errorf("cannot open target file: %w", err)
	}

	defer f.Close()

	n, err := r.GetNumPages()

	if err != nil {
		return fmt.Errorf("cannot list pages: %w", err)
	}

	c := creator.New()

	for i := 0; i < n; i++ {
		page, err := r.GetPage(i + 1)

		if err != nil {
			return fmt.Errorf("cannot get page: %w", err)
		}

		c.AddPage(page)

		img, err := creator.NewGraphicSVGFromFile(name)

		if err != nil {
			return fmt.Errorf("cannot create image: %w", err)
		}

		img.ScaleToWidth(150)
		img.SetPos(50, 100)

		if err := c.Draw(img); err != nil {
			return fmt.Errorf("cannot draw image: %w", err)
		}
	}

	dest := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	dest = filepath.Join("out", fmt.Sprintf("%s.pdf", dest))

	if err := c.WriteToFile(dest); err != nil {
		return fmt.Errorf("cannot write output file: %w", err)
	}

	fmt.Printf("wrote output file: %s\n", dest)

	return nil
}
