package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

const baseURL = "https://www.toptal.com/developers/gitignore/api/"

func banner() {
	fmt.Println(`
	___ ____ ___ _____ 
	|_ _/ ___|_ _|_   _|
	 | | |  _ | |  | |  
	 | | |_| || |  | |  
	|___\____|___| |_|  
						

	`)
}

func main() {
	banner()

	helpFlag := flag.Bool("h", false, "Display help manual")
	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	gitignoreCreation()
}

func showHelp() {
	fmt.Println("Help:")
	fmt.Println("    -h: Display help manual")
	fmt.Printf("Example: \n")
	fmt.Println("Enter the required tech stack for your project:")
	fmt.Println("node")
	fmt.Println("Enter the location for the .gitignore file:")
	fmt.Println("Enter '.' to indicate the current directory")
	fmt.Println(".")
	fmt.Println(color.GreenString(("Successfully generated .gitignore file at _")))
}

func generateGitignore(techStack string, location string) error {
	if _, err := os.Stat(location); os.IsNotExist(err) {
		err = os.MkdirAll(location, os.ModePerm)
		if err != nil {
			return err
		}
	}

	gitignorePath := filepath.Join(location, ".gitignore")
	if _, err := os.Stat(gitignorePath); !os.IsNotExist(err) {
		os.Remove(gitignorePath)
	}

	url := fmt.Sprintf("%s%s", baseURL, techStack)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMsg := "Failed to retrieve .gitignore file. Status: " + resp.Status
		log.Fatalf(color.RedString(errorMsg))
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return err
	}

	file, err := os.Create(gitignorePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(buf.String())
	if err != nil {
		return err
	}

	successMessgae := "Successfully generated .gitignore file at %s\n"

	fmt.Println(color.GreenString(successMessgae, gitignorePath))
	return nil
}

func gitignoreCreation() {
	fmt.Println("Enter the required tech stack of your project:")
	reader := bufio.NewReader(os.Stdin)
	techStack, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %s", err)
	}
	techStack = strings.TrimSpace(techStack)

	fmt.Println("Enter the location for the .gitignore file:")
	fmt.Println("Enter '.' to indicate the current directory")

	location, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %s", err)
	}
	location = strings.TrimSpace(location)

	if location == "." {
		location, err = os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %s", err)
		}
	} else {
		location = strings.TrimSuffix(location, "\n")
	}

	err = generateGitignore(techStack, location)
	if err != nil {
		log.Fatalf("Failed to generate .gitignore file: %s", err)
	}
}
