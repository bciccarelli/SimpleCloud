package main

import (
  	"crypto/sha256"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"encoding/hex"
	"encoding/base64"
  	"net/http"
  	"strings"
  	"golang.org/x/crypto/pbkdf2"
)
const (
	password = "This is my password"
)
var dk []byte
var salt string
var err error
func handle404(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "404: /" + message +" not found"
	w.Write([]byte(message))
}
func returnFile(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")

	salt,err = GenerateRandomString(8)
	check(err)
	dk = pbkdf2.Key([]byte(password), []byte(salt), 100, 16, sha256.New)
	
	out := ""
	files := readFS()
	for _, files := range files {
        f, err := ioutil.ReadFile(strings.TrimPrefix("files/"+files, "/"))
		check(err);
		if(strings.HasSuffix(files, ".jpg")||strings.HasSuffix(files, ".png")){
			out += base64.StdEncoding.EncodeToString(f)
		} else {
			out += string(f)

		}
		out += "!fileName" + strings.TrimPrefix(files, "/")+"!fileName"
    }
	

	data := padRight([]byte(out))

	block, err := aes.NewCipher(dk)
	check(err);
	
	ciphertext := make([]byte, len(data))
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, []byte(data))

	final := hex.EncodeToString(iv)+":"+hex.EncodeToString(ciphertext)+":"+salt

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(final))
	fmt.Println("Sent.")
}

func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    // Note that err == nil only if we read len(b) bytes.
    if err != nil {
        return nil, err
    }

    return b, nil
}

func GenerateRandomString(s int) (string, error) {
    b, err := GenerateRandomBytes(s)
    return base64.URLEncoding.EncodeToString(b), err
}


func main() {
	http.HandleFunc("/", returnFile)


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
func readFS() ([]string) {
	m := make([]string,0)
	dirname := "./files"
    d, err := os.Open(dirname)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
    fi, err := d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _, fi := range fi {
        if fi.Mode().IsRegular() {
        	m = append(m,fi.Name())
        }
    }
    return m
}