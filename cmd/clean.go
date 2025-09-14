package cmd

import (
	"fmt"

	"github.com/Lily-404/todo/internal/config"
	"github.com/Lily-404/todo/internal/i18n"
	"github.com/Lily-404/todo/internal/renderer"
	"github.com/Lily-404/todo/internal/storage"
	"github.com/Lily-404/todo/pkg/logger"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:     "clean",
	Aliases: []string{"c"},
	Short:   i18n.GetMessage(config.GetConfig().Language, "cmd_clean_short"),
	RunE: func(cmd *cobra.Command, args []string) error {
		notes, err := storage.ListNotes()
		if err != nil {
			return err
		}

		var unfinishedNotes []storage.Note
		var finishedNotes []storage.Note
		for _, note := range notes {
			if note.Status == "done" {
				finishedNotes = append(finishedNotes, note)
			} else {
				unfinishedNotes = append(unfinishedNotes, note)
			}
		}

		if len(finishedNotes) == 0 {
			color.Yellow(i18n.GetMessage(config.GetConfig().Language, "no_completed_tasks"))
			return nil
		}

		color.HiCyan("\n" + i18n.GetMessage(config.GetConfig().Language, "tasks_to_clean"))
		for i, note := range finishedNotes {
			color.HiBlack(fmt.Sprintf("  %d. %s", i+1, note.Content))
		}
		fmt.Println()

		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Are you sure you want to delete %d completed tasks?", len(finishedNotes)),
			IsConfirm: true,
		}

		if _, err := prompt.Run(); err != nil {
			fmt.Println("Clean cancelled.")
			return nil
		}

		if saveErr := storage.SaveNotes(unfinishedNotes); saveErr != nil {
			return saveErr
		}

		logger.Success(i18n.GetMessage(config.GetConfig().Language, "cleaned_tasks", len(finishedNotes)))

		fmt.Println("\n" + i18n.GetMessage(config.GetConfig().Language, "current_tasks"))
		updatedNotes, err := storage.ListNotes()
		if err != nil {
			return err
		}
		renderer.RenderNotes(updatedNotes, false, "")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}