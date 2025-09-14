package cmd

import (
	"fmt"

	"github.com/Lily-404/todo/internal/config"
	"github.com/Lily-404/todo/internal/i18n"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version = "v1.2.2"
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: i18n.GetMessage(config.GetConfig().Language, "cmd_root_short"),
	Long:  color.HiWhiteString(i18n.GetMessage(config.GetConfig().Language, "cmd_root_long")),
	Version: fmt.Sprintf(`
 _______   _______           _________  _______   _______   _______ 
|\   ____\|\   __  \        |\___   ___\\   __  \|\   ___ \|\   __  \    
\ \  \___|\ \  \|\  \       \|___ \  \_\ \  \|\  \ \  \_|\ \ \  \|\  \   
 \ \  \  __\ \  \\\  \           \ \  \ \ \  \\\  \ \  \ \\ \ \  \\\  \  
  \ \  \|\  \ \  \\\  \           \ \  \ \ \  \\\  \ \  \_\\ \ \  \\\  \ 
   \ \_______\ \_______\           \ \__\ \ \_______\ \_______\ \_______\
    \|_______|\|_______|            \|__|  \|_______|\|_______|\|_______|                                                                                                                                                                                      
Version: %s
`, version),
}

func init() {
	// 设置帮助和版本标志的描述
	rootCmd.PersistentFlags().BoolP("help", "h", false, i18n.GetMessage(config.GetConfig().Language, "help_flag_short"))
	rootCmd.PersistentFlags().BoolP("version", "v", false, i18n.GetMessage(config.GetConfig().Language, "version_flag_short"))

	// 隐藏 completion 和 help 命令
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// 设置自定义的帮助模板
	helpTemplate := `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if .HasAvailableSubCommands}}` + i18n.GetMessage(config.GetConfig().Language, "available_commands") + `{{range .Commands}}{{if (and (not .Hidden) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

{{end}}{{if .HasAvailableLocalFlags}}` + i18n.GetMessage(config.GetConfig().Language, "flags") + `
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

` + i18n.GetMessage(config.GetConfig().Language, "help_usage") + `
`
	rootCmd.SetHelpTemplate(helpTemplate)

	// 设置一个隐藏的帮助命令
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "help [command]",
		Hidden: true,
	})
}

func Execute() error {
	return rootCmd.Execute()
}
