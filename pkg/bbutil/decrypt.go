package bbutil

import (
	"github.com/proglottis/gpgme"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// DecryptFile decrypts a single file.
func (bbu *RepoInfo) DecryptFile(filename, group string, overwrite bool) error {

	origFilepath := filepath.Join(bbu.RepoBaseDir, filename)
	gpgFilepath := strings.Join([]string{filepath.Join(bbu.RepoBaseDir, filename), ".gpg"}, "")


	gpgFile, err := os.Open(gpgFilepath)
	if err != nil {
		log.Warnf("Error while opening file %v: %v", gpgFilepath, err)
		return err
	}
	fileInfo, err := gpgFile.Stat()

	decrypted, err := gpgme.Decrypt(gpgFile)
	if err != nil {
		log.Warnf("Error while decrypting file: %v", err)
		return err
	}

	flag := os.O_CREATE|os.O_WRONLY
	if overwrite {
		flag |= os.O_TRUNC

	} else {
		flag |= os.O_EXCL
	}


	outFile, err := os.OpenFile(origFilepath, flag, fileInfo.Mode())
	if err != nil {
		log.WithError(err).Warn("Error creating output file.")
		return err
	}

	if _, err := io.Copy(outFile, decrypted) ; err != nil {
		log.WithError(err).Warn("Error writing to output file.")
		return err
	}
	if err = gpgFile.Close() ; err != nil {
		return err
	}

	if err = outFile.Close() ; err != nil {
		return err
	}

	if err = decrypted.Close() ; err != nil {
		return err
	}

	if group != "" {
		group, err := user.LookupGroup(group)
		if err != nil {
			return err
		}
		gid, err := strconv.ParseInt(group.Gid, 10, 64)
		if err != nil {
			return err
		}

		if err = os.Chown(origFilepath, -1, int(gid)) ; err != nil {
			return err
		}
	}


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
