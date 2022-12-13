package main

import (
	"bufio"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	listenAddr string
	socks5Addr string
	filename   string
	urlRegexes []*regexp.Regexp
)

var rootCmd = &cobra.Command{
	Use:   "s2h",
	Short: "A simple tool to convert socks5 proxy protocol to http proxy protocol",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if filename != "" {
			urlRegexes, err = parseUrlRegexes(filename)

			if err != nil {
				logrus.Fatal(err)
			}
		} else {
			logrus.Info("No urls explicitly set, always using proxy.")
		}

		logrus.Info("Starting Socks5 Proxy Convert Server...")
		logrus.Infof("HTTP Listen Address: %s", listenAddr)
		logrus.Infof("Socks5 Server Address: %s", socks5Addr)

		err = http.ListenAndServe(listenAddr, http.HandlerFunc(serveHTTP))
		if err != nil {
			logrus.Error(err)
		}
	},
}

func parseUrlRegexes(filename string) ([]*regexp.Regexp, error) {
	var result []*regexp.Regexp
	f, err := os.Open(filename)
	if err != nil {
		return result, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		regex, err := regexp.Compile(line)
		if err != nil {
			return result, err
		}
		result = append(result, regex)
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func match(regexes []*regexp.Regexp, url string) bool {
	host := []byte(strings.Split(url, ":")[0])

	for _, regex := range regexes {
		if regex.Match(host) {
			return true
		}
	}
	return false
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-05-04 15:02:01",
	})

	rootCmd.PersistentFlags().StringVarP(&listenAddr, "listen", "l", "0.0.0.0:8081", "http listen address")
	rootCmd.PersistentFlags().StringVarP(&socks5Addr, "socks5", "s", "127.0.0.1:1080", "remote socks5 listen address")
	rootCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "file for positive url filters that go through the proxy")
}
