package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with somebody",
	Run: func(cmd *cobra.Command, args []string) {
		personality, _ := cmd.Flags().GetString("personality")
		fn := strings.Fields(personality)[0]

		prompt := "You answer in the speaking style of " + personality + "."
		if pr, _ := cmd.Flags().GetString("prompt"); pr != "" {
			prompt = pr
			fn = "Response"
		}

		colors, _ := cmd.Flags().GetBool("colors")
		sess := c.NewChatSession(prompt)

		r := bufio.NewReader(os.Stdin)
		go func() {
			<-ctx.Done()
			os.Stdin.SetDeadline(time.Now())
		}()

		colorize := func(text string, code int) string {
			if colors {
				return fmt.Sprintf("\033[%dm%s\033[0m", code, text)
			}
			return text
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			fmt.Print(colorize("You: ", 32))
			msg, err := r.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			res, err := sess.Complete(ctx, msg)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(colorize(fn+": ", 33))
			fmt.Println(func() string {
				if colors {
					var b bytes.Buffer
					err := quick.Highlight(&b, res, "markdown", "terminal16m", "monokai")
					if err != nil {
						log.Fatal(err)
					}

					return b.String()
				}

				return res
			}())
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().String("prompt", "", "A prompt to override the default")
	chatCmd.Flags().Bool("colors", true, "Colorize the output")
	chatCmd.Flags().String("personality", "Sigmund Freud", "A prompt to override the default")
}
