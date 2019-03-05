package bbutil

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
	"path"
	"strings"
)

func (bbu *RepoInfo) PostDeploy(privateKeyFile io.Reader) error {
	fnames, valid, err := bbu.FileIterator(true, nil)
	if err != nil {
		return err
	}

	el, err := openpgp.ReadArmoredKeyRing(privateKeyFile)
	if err != nil {
		return err
	}

	for i, filename := range fnames {
		if valid[i] {
			if err := bbu.decryptFileWith(filename, &el) ; err != nil {
				return err
			}
		}
	}
	return nil
}

func (bbu *RepoInfo) decryptFileWith(filename string, key *openpgp.EntityList) error {
	origPath := path.Join(bbu.RepoBaseDir, filename)
	file, err := os.Open(strings.Join([]string{origPath, ".gpg"}, ""))
	if err != nil {
		return err
	}
	md, err := openpgp.ReadMessage(file, key, nil, nil)
	if err != nil {
		return err
	}

	writer, err := os.Create(origPath)
	if err != nil {
		return err
	}
	n, err := io.Copy(writer, md.UnverifiedBody)
	log.Infof("Decrypted %d bytes into %s", n, origPath)

	//compressed, err := gzip.NewReader(md.UnverifiedBody)
	//if err != nil {
	//	return err
	//}
	//defer compressed.Close()
	//
	//n, err := io.Copy(os.Stdout, compressed)
	//log.Infof("Decrypted %d bytes", n)

	return nil
}

// DecryptFile decrypts a single file.
func (bbu *RepoInfo) DecryptFile(filename, group string, overwrite bool) error {


	log.Infof("Decrypted %q", filename)


	//fmt.Fprintf(os.Stderr, "WOULD DECRYPT: %v %q %q\n", overwrite, group, filename)

	// export PATH=/usr/bin:/bin:"$PATH"

	// # Decrypt:
	// echo '========== Decrypting new/changed files: START'
	// while IFS= read <&99 -r unencrypted_file; do
	//   encrypted_file=$(get_encrypted_filename "$unencrypted_file")
	//   decrypt_file_overwrite "$encrypted_file" "$unencrypted_file"
	//   cp_permissions "$encrypted_file" "$unencrypted_file"
	//   if [[ ! -z "$FILE_GROUP" ]]; then
	//     chmod g+r "$unencrypted_file"
	//     chgrp "$FILE_GROUP" "$unencrypted_file"
	//   fi
	// done 99<"$BB_FILES"

	// echo '========== Decrypting new/changed files: DONE'

	return nil
}
