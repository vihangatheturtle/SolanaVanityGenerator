package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
)

var genned = 0
var genned5 = 0
var genned10 = 0
var genned30 = 0
var genned60 = 0
var threads = 16
var t1 = time.Now()
var refreshTime = time.Now()
var searchFor = "mjk134"
var closeThreads = false
var s5avg = time.Now()
var s10avg = time.Now()
var s30avg = time.Now()
var s60avg = time.Now()
var five = ""
var ten = ""
var thirty = ""
var sixty = ""

func genWallet() {
	for {
		if closeThreads {
			return
		}
		foundWallet := solana.NewWallet()
		if strings.HasPrefix(foundWallet.PublicKey().String(), searchFor) && !closeThreads {
			firstOffsetChar := strings.Split(foundWallet.PublicKey().String(), searchFor)[1][0:1]
			if firstOffsetChar == strings.ToUpper(firstOffsetChar) {
				fmt.Println("Found wallet:", foundWallet.PublicKey().String())
				fmt.Print("Private key: ")
				fmt.Println(foundWallet.PrivateKey)
				fmt.Print("Took " + strconv.Itoa(genned+1) + " attempts and ")
				fmt.Print(time.Since(t1))
				closeThreads = true
			}
		}
		genned++
		genned5++
		genned10++
		genned30++
		genned60++
		if genned%1000000 == 0 {
			fmt.Print("Generated " + strconv.Itoa(genned) + " wallets, took ")
			fmt.Println(time.Since(t1))
		}
	}
}

func timer() {
	for {
		if closeThreads {
			return
		}
		if time.Since(refreshTime).Seconds() > 1 {
			update := false
			if time.Since(s5avg).Seconds() >= 5 {
				five = strconv.Itoa(genned5 / int(time.Since(s5avg).Seconds()))
				genned5 = 0
				s5avg = time.Now()
			}
			if time.Since(s10avg).Seconds() >= 10 {
				ten = strconv.Itoa(genned10 / int(time.Since(s10avg).Seconds()))
				genned10 = 0
				s10avg = time.Now()
			}
			if time.Since(s30avg).Seconds() >= 30 {
				thirty = strconv.Itoa(genned30 / int(time.Since(s30avg).Seconds()))
				genned30 = 0
				s30avg = time.Now()
				update = true
			}
			if time.Since(s60avg).Seconds() >= 60 {
				sixty = strconv.Itoa(genned60 / int(time.Since(s60avg).Seconds()))
				genned60 = 0
				s60avg = time.Now()
			}
			if update {
				fmt.Println("Generating at (avg wallets per second) | 5sec:", five, "| 10sec:", ten, "| 30sec:", thirty, "| 60sec:", sixty)
				refreshTime = time.Now()
			}
		}
	}
}

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter search string: ")
		searchFor, _ = reader.ReadString('\n')
		break
	}
	fmt.Println("Searching for:", searchFor)
	for i := 0; i < threads; i++ {
		go genWallet()
	}
	go timer()
	fmt.Scanln()
}
