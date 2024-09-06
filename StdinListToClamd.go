package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/thc2cat/go-clamd"
)

func main() {

	clamavHostPtr := flag.String("host", "lorraine.ens.uvsq.fr", "Adresse IP du démon clamd")
	clamavPortPtr := flag.Int("port", 3310, "Port du démon clamd")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	c := clamd.NewClamd(clamd.WithTCP(*clamavHostPtr, *clamavPortPtr))
	if c == nil {
		fmt.Printf("Erreur lors de la création du client ClamAV\n")
		return
	}

	if ok, cerr := c.Ping(ctx); !ok {
		fmt.Println("clamd server error !", cerr)
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

// Fonction pour envoyer un fichier vers l'antivirus
func sendToClamav(path string, c *clamd.Clamd, ctx context.Context) error {

	if ok, cerr := c.Ping(ctx); !ok {
		fmt.Println("clamd server error !", cerr)
		return cerr
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	_, scanerr := c.ScanStream(ctx, reader)

	return scanerr
}
