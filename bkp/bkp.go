/********************************************************************************
* La funzione iterara su ogni file e carpella al interno di PATH e fare due cose:
* 1. limpiara i contenuti dei file di tutti i id uniche forniti da notion
* 2. Rimovera la id unica dal nome dil file.
* ******************************************/

package bkp

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

const PATH2 ="notion_12.22.24"

func Bkp(){ 

	cleanNotionLinks__(PATH2)
}


func cleanNotionLinks__(path string) error {
    // Regular expression to match Notion's unique ID pattern
    // This matches strings that look like " 1234abcd" at the end of text
    	idPatternWithSpace := regexp.MustCompile(`\b[a-zA-Z0-9]*\.?[a-zA-Z0-9]{16,}\b`)
   	idPattern := regexp.MustCompile(`^[a-zA-Z0-9]{16,}$`)
	totalLinksFound := 0

    // Walk through all files in the workspace
    err := filepath.Walk(path, 
	func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Skip if not an HTML file
        if !strings.HasSuffix(strings.ToLower(path), ".html") {
            return nil
        }

        // Open the file
        file, err := os.Open(path)
        if err != nil {
            return fmt.Errorf("error opening file %s: %v", path, err)
        }
        defer file.Close()

        // Parse HTML
        doc, err := html.Parse(file)
        if err != nil {
            return fmt.Errorf("error parsing HTML in %s: %v", path, err)
        }

        // Track if we made any changes
        modified := false

        // Function to recursively process nodes
        var processNode func(*html.Node)
        processNode = func(n *html.Node) {
            if n.Type == html.ElementNode {
			fmt.Println(n.Type)
			return

                // Process href attributes in anchor tags
                if n.Data == "a" {
                    for i, attr := range n.Attr {
				// Original href value 
				ogName := n.Attr[i].Val

                        if attr.Key == "href" {



					// if the href is a link to a file and not a dir will have .html extension
					isHTMLFile := false

					// new name formed without the id 
					var newNameParts = make([]string, 0)

					// Check for an ID at the herf level by dir "/"
					fmt.Printf("◉ All parts in href: %s \n", strings.Join(strings.Split(attr.Val, "/"), ", "))
					
					for _, part := range strings.Split(attr.Val, "/") {
						fmt.Printf("  ◎ Processing part: %s \n", part)
						if idPatternWithSpace.MatchString(part){
							// If it matches, then check by space level
							fmt.Printf("    ❌ %s part is dirty! Checking by space now \n", part)
							
							nameParts := make([]string, 0)
							
							fmt.Printf("      ◦ Parts by space: %s \n", strings.Join(strings.Split(part, " "), ", "))
							for _, word := range strings.Split(part, " ") {
								fmt.Printf("        ・Processing word: %s \n", word)

								// remove html from the string if this is a file file
								wordWithoutExtension, hasExt := strings.CutSuffix(word, ".html")
								isHTMLFile = hasExt

								if idPattern.MatchString(wordWithoutExtension) {
									totalLinksFound++
									fmt.Printf("          ❌ %s is dirty! removing... \n", wordWithoutExtension)
									modified = true
									
								} else{
									fmt.Printf("          ✅ %s is clean! skipping... \n", wordWithoutExtension)
									nameParts = append(nameParts, wordWithoutExtension)
								}
							}

							// since this was originally split by space, join them by space
							newName := strings.Join(nameParts, " ")

							// if it is a link to an html file, add the extension back
							if isHTMLFile {
								newName+= ".html"
							}

							// add this name to the newNameParts 
							newNameParts = append(newNameParts, newName)
						} else {
							fmt.Printf("    ✅ %s is clean! skipping... \n", part)
							newNameParts = append(newNameParts, part)
						}
					}

					// since this was originally split by /, join them by /
					ogName = strings.Join(newNameParts, "/")
					n.Attr[i].Val = ogName

					fmt.Printf("📝 New name will be:  %s \n", ogName)
                        }

				// change the displayable text too
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					fmt.Println("----", c.Data)
					c.Data = ogName
				}
                    }
                }
            }

            // Process child nodes
            for c := n.FirstChild; c != nil; c = c.NextSibling {
                processNode(c)
            }
        }

        // Process the document
        processNode(doc)

        // If modifications were made, write the changes back to the file
        if modified {
            // Create a temporary file
            tempFile, err := os.CreateTemp(filepath.Dir(path), "temp-*.html")
            if err != nil {
                return fmt.Errorf("error creating temp file for %s: %v", path, err)
            }
            tempPath := tempFile.Name()

		fmt.Printf("🪄🦄 Temp file %s created", tempPath)

            // Write the modified HTML to the temp file
            err = html.Render(tempFile, doc)
            tempFile.Close()
            if err != nil {
                os.Remove(tempPath)
                return fmt.Errorf("error writing to temp file for %s: %v", path, err)
            }

            // Replace the original file with the temp file
            err = os.Rename(tempPath, path)
            if err != nil {
                os.Remove(tempPath)
                return fmt.Errorf("error replacing file %s: %v", path, err)
            }

            fmt.Printf("✅✨ File Processed: %s\n", path)
        }

        return nil
    })

    if err != nil { 
	return err
    }

	// how many links were found?
	fmt.Printf("🔢 total links found: %d \n", totalLinksFound)
	return nil
}
