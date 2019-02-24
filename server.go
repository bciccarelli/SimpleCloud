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
    //"encoding/json"
  	"net/http"
  	"strings"
  	"golang.org/x/crypto/pbkdf2"
  	"time"
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

type write struct {
    test string
}

func returnFile(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	if strings.TrimPrefix(message, "/write")==message {
		(w).Header().Set("Access-Control-Allow-Origin", "*")

		salt,err = GenerateRandomString(8)
		check(err)
		dk = pbkdf2.Key([]byte(password), []byte(salt), 100, 16, sha256.New)
		
		out := ""
		files := readFS("running")
		for _, files := range files {
	        f, err := ioutil.ReadFile("running/"+files)
			check(err);
			if(strings.HasSuffix(files, ".jpg")||strings.HasSuffix(files, ".png")){
				out += base64.StdEncoding.EncodeToString(f)
			} else {
				//out += string(f)

			}
			out += "!.!"
			out += files
			out += "!.!"
	    }
	    out+="!s!"
	    files = readFS("stopped")
		for _, files := range files {
	        f, err := ioutil.ReadFile("stopped/"+files)
			check(err);
			if(strings.HasSuffix(files, ".jpg")||strings.HasSuffix(files, ".png")){
				out += base64.StdEncoding.EncodeToString(f)
			} else {
				//out += string(f)

			}
			out += "!.!"
			out += files
			out += "!.!"
	    }
	    out += "........."
		data := padRight(out)

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
	} else {
		a := (r.Header["T"][0])
		c,err := hex.DecodeString(strings.Split(a,":")[1])
		check(err);
		s,err:= hex.DecodeString(strings.Split(a,":")[2])
		check(err);
		dk = pbkdf2.Key([]byte(password), []byte(s), 100, 16, sha256.New)
		block, err := aes.NewCipher(dk)
		iv,err := hex.DecodeString(strings.Split(a,":")[0])
		check(err);

		mode := cipher.NewCBCDecrypter(block, iv)

		mode.CryptBlocks(c, c)
		fileName := c[:strings.Index(string(c), ":")]
		write := c[strings.Index(string(c), ":")+1:]
		err = ioutil.WriteFile(string(fileName), write, 0644)
    	check(err)
	}
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
	fmt.Println("Started")
  	go manage()
	http.HandleFunc("/", returnFile)
  	if err := http.ListenAndServe(":3000", nil); err != nil {
    	panic(err)
  	}
}
func manage() {
	for (true) {
		time.Sleep(time.Second)
		if() {
				
		}
	}
}
func padRight(str string) (string) {
	m := str
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
func readFS(path string) ([]string) {
	m := make([]string,0)
	dirname := "./"+path
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