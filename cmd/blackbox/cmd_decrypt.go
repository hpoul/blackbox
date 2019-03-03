package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/StackExchange/blackbox/pkg/bbutil"
	"github.com/pkg/errors"
)

func cmdDecrypt(allFiles bool, filenames []string, group string) error {
	bbu, err := bbutil.New()
	if err != nil {
		return err
	}

	// prepare_keychain

	fnames, valid, err := bbu.FileIterator(allFiles, filenames)
	if err != nil {
		return errors.Wrap(err, "decrypt")
	}
	for i, filename := range fnames {
		if valid[i] {
			if err := bbu.DecryptFile(filename, group, true) ; err != nil {
				logrus.WithError(err).Panicf("Error decrypting file")
			}
		} else {
			fmt.Fprintf(os.Stderr, "SKIPPING: %q\n", filename)
		}
	}

	return nil
}
