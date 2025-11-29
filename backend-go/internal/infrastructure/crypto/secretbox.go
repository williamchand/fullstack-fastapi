package crypto

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "golang.org/x/crypto/nacl/secretbox"
)

func Encrypt(key []byte, plaintext []byte) (string, error) {
    var k [32]byte
    if len(key) != 32 {
        return "", errors.New("invalid key length")
    }
    copy(k[:], key)
    var nonce [24]byte
    if _, err := rand.Read(nonce[:]); err != nil {
        return "", err
    }
    out := secretbox.Seal(nonce[:], plaintext, &nonce, &k)
    return base64.StdEncoding.EncodeToString(out), nil
}

func Decrypt(key []byte, ciphertext string) ([]byte, error) {
    var k [32]byte
    if len(key) != 32 {
        return nil, errors.New("invalid key length")
    }
    copy(k[:], key)
    raw, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return nil, err
    }
    var nonce [24]byte
    copy(nonce[:], raw[:24])
    decrypted, ok := secretbox.Open(nil, raw[24:], &nonce, &k)
    if !ok {
        return nil, errors.New("decryption failed")
    }
    return decrypted, nil
}
