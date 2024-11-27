package main
import (
  "github.com/cli/go-gh"
  "fmt"
)

func main() {
  args := []string{"api", "user", "--jq", `"You are @\(.login) (\(.name))"` }
  stdOut, _, err := gh.Exec(args...)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(stdOut.String())
}
