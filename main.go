package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var textSuffix = map[bool]string{
	false: "at the start",
	true:  "at the end",
}

var textMany = map[bool]string{
	false: "a key",
	true:  "keys",
}

var rootCmd = &cobra.Command{
	Use:     "pkdns-vanity <text>",
	Example: "  pkdns-vanity woop\n  pkdns-vanity 1234 --suffix",
	Short:   "Generate vanity PKDNS domain names",
	Long:    "This tool abuses your CPU to generate vanity PKDNS domains starting (or ending) with characters of your choosing. More than 4 characters will take minutes, more than 5 can take hours. More than 6 and you'll only generate heat.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}

		text := args[0]
		if _, err := zBase32.DecodeString(text); err != nil {
			chars := []rune(zBase32Chars)
			sort.Slice(chars, func(i, j int) bool {
				return chars[i] < chars[j]
			})
			fmt.Fprintf(os.Stderr, "Error: '%s' includes invalid characters (not in %s)\n", text, string(chars))
			os.Exit(1)
		}

		suffix, err := cmd.Flags().GetBool("suffix")
		if err != nil {
			suffix = false
		}

		many, err := cmd.Flags().GetBool("many")
		if err != nil {
			many = false
		}

		match := func(pub string) bool { return strings.HasPrefix(pub, text) }
		if suffix {
			match = func(pub string) bool { return strings.HasSuffix(pub, text) }
		}

		fmt.Printf("Looking for %s with '%s' %s…\n", textMany[many], text, textSuffix[suffix])

		keys := make(chan key)

		for i := 1; i <= runtime.NumCPU(); i++ {
			go func() {
				panic(findKeys(match, keys))
			}()
		}

		foundAny := false
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)

	Listener:
		for {
			select {
			case k := <-keys:
				fmt.Printf("  %s <- %s\n", k.pub, k.prv)
				foundAny = true
				if !many {
					break Listener
				}
			case <-sigs:
				break Listener
			}
		}

		explain(foundAny)
		return nil
	},
}

func explain(foundAny bool) {
	if foundAny {
		fmt.Println("  Public Key <- Private Key")
		fmt.Println()
		fmt.Println("Put a Private Key in ~/.pkdns/seed.txt to use it with pkdns-cli")
	} else {
		fmt.Println()
		fmt.Println("No keys were found.")
		fmt.Println("Perhaps try a shorter key?")
	}
}

type key struct {
	prv string
	pub string
}

func main() {
	rootCmd.Flags().BoolP("suffix", "s", false, "seek keys with the vanity phrase at the end, not the start")
	rootCmd.Flags().BoolP("many", "m", false, "keep going after finding one")

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

var zBase32Chars = "ybndrfg8ejkmcpqxot1uwisza345h769"
var zBase32 = base32.NewEncoding(zBase32Chars).WithPadding(base32.NoPadding)

func findKeys(match func(string) bool, keys chan<- key) error {
	for {
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}

		pub := zBase32.EncodeToString(pubKey)
		if match(pub) {
			keys <- key{
				pub: pub,
				prv: zBase32.EncodeToString(privKey.Seed()),
			}
		}
	}
}
