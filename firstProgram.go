package main

import (
  	"crypto/sha256"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"encoding/hex"
	"encoding/base64"
  	"net/http"
  	"strings"
  	"golang.org/x/crypto/pbkdf2"
)
const (
	password = "banana"
	salt = "salt"
)
var dk []byte
func handle404(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "404: /" + message +" not found"
	w.Write([]byte(message))
}
func returnFile(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	file, err := ioutil.ReadFile("red.jpg")
	check(err);

	data := base64.StdEncoding.EncodeToString(file)

	block, err := aes.NewCipher(dk)
	check(err);
	
	ciphertext := make([]byte, len(data))
	iv := make([]byte,aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, []byte(data))

	final := hex.EncodeToString(iv)+":"+hex.EncodeToString(ciphertext)
    fmt.Println(len(string(data)))

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(final))
}

func main() {

	dk = pbkdf2.Key([]byte(password), []byte(salt), 100, 16, sha256.New)
	http.HandleFunc("/", handle404)
  	http.HandleFunc("/red", returnFile)

  	if err := http.ListenAndServe(":3000", nil); err != nil {
    	panic(err)
  	}
}
func padRight(str []byte) (string) {
	m := string(str)
	for i := 0; i < (16-(len(str)%16)); i++ {
		m += " "
	}
	return m
}
func check(e error) {
    if e != nil {
        panic(e)
    }
}