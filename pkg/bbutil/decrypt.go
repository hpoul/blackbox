package bbutil

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
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
	//
	//
	//block, err := armor.Decode(privateKeyFile)
	//if err != nil {
	//	return err
	//}
	//
	//privateKeyReader := packet.NewReader(block.Body)
	//privateKeyPacket, err := privateKeyReader.Next()
	////entity := openpgp.EntityList{}
	//entity := openpgp.Entity{}
	//for {
	//	next, err := privateKeyReader.Next()
	//	if err != nil {
	//		return err
	//	}
	//	privateKey, ok := next.(*packet.PrivateKey)
	//	if ok {
	//		entity.PrivateKey = privateKey
	//	}
	//
	//}
	//if !ok {
	//	return errors.New("no private key found")
	//}
	//
	//next, err := privateKeyReader.Next()
	//log.Infof("next: %v, err: %v", next, err)

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
	file, err := os.Open(strings.Join([]string{bbu.RepoBaseDir, "/", filename, ".gpg"}, ""))
	if err != nil {
		return err
	}
	md, err := openpgp.ReadMessage(file, key, nil, nil)
	if err != nil {
		return err
	}

	n, err := io.Copy(os.Stdout, md.UnverifiedBody)
	log.Infof("Read %d bytes", n)

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
