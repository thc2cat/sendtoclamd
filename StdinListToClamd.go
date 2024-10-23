package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/thc2cat/go-clamd"
)

var (
	clamavHostPtr        = flag.String("h", "clamd.inyourdns.net", "remote clamd host")
	clamavPortPtr        = flag.Int("p", 3310, "clamd port to connect")
	defaultMaxStreamSize = (int64)(52428800)
)

// StreamMaxLength 50M
// You also have to change the ClamClient instance with a higher MaxStreamSize:

// var client = new ClamClient("localhost", 3310)
//
//	{
//		     MaxStreamSize = 52428800
//		};
func main() {

	flag.Parse()

	// ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	ctx := context.TODO()
	// defer cancel()

	c := NewClamd()

	if !PingClamdOk(c, ctx) {
		os.Exit(-1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		filename := scanner.Text()
		errs := sendToClamav(filename, c, ctx)
		if errs != nil {
			fmt.Printf("%s: %v\n", filename, errs)
		}
	}

}

func NewClamd() *clamd.Clamd {
	c := clamd.NewClamd(clamd.WithTCP(*clamavHostPtr, *clamavPortPtr))
	if c == nil {
		fmt.Fprintf(os.Stderr, "Error in clamd.NewClamd(clamd.WithTCP(host, port)) \n")
	}
	return c
}
func PingClamdOk(c *clamd.Clamd, ctx context.Context) bool {
	if ok, cerr := c.Ping(ctx); !ok {
		fmt.Fprintf(os.Stderr, "503 Clamd Ping  error : %v !", cerr)
		return false
	}
	return true
}

func CheckFilesize(thePath string) error {
	fileInfo, err := os.Stat(thePath)
	if err != nil {
		return fmt.Errorf("404 File vanished: %s %v", thePath, err)
	}
	if fileInfo.Size() > defaultMaxStreamSize {
		return fmt.Errorf("413 File size of %s is larger (%d) than default Max Stream Size (%v)", thePath, fileInfo.Size(), defaultMaxStreamSize)
	}
	if fileInfo.Mode().IsDir() {
		return nil
	}
	return nil
}

// Fonction pour envoyer un fichier vers l'antivirus
func sendToClamav(path string, c *clamd.Clamd, ctx context.Context) error {

	for retries := 3; !PingClamdOk(c, ctx) && retries > 0; retries-- {
		fmt.Fprintf(os.Stderr, "503 Clamd server error retrying")
		c = NewClamd()
	}

	if ok, cerr := c.Ping(ctx); !ok {
		fmt.Fprintf(os.Stderr, "503 Clamd server error : %v !", cerr)
		return cerr
	}

	if errCkFS := CheckFilesize(path); errCkFS != nil {
		return nil // This is not a Clamav ERROR but acces error
	}

	file, err := os.Open(path)
	if err != nil { // Specifics errors ( bad links and others )
		// fmt.Fprintln(os.Stderr, err)
		return nil
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	_, scanerr := c.ScanStream(ctx, reader)

	return scanerr
}
