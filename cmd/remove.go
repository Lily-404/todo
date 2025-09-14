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

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   i18n.GetMessage(config.GetConfig().Language, "cmd_remove_short"),
	RunE: func(cmd *cobra.Command, args []string) error {
		notes, err := storage.ListNotes()
		if err != nil {
			return err
		}

		// 过滤未完成的任务
		var unfinishedNotes []storage.Note
		for _, note := range notes {
			if note.Status != "done" {
				unfinishedNotes = append(unfinishedNotes, note)
			}
		}

		if len(unfinishedNotes) == 0 {
			color.Yellow(i18n.GetMessage(config.GetConfig().Language, "no_pending_tasks"))
			return nil
		}

		// 创建任务选择提示
		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "→ {{ .Title | cyan }} {{ .Content | white | bold }} {{ if .DueDate }}({{ .DueDate | magenta }}){{ end }} [{{ .Priority | red }}]",
			Inactive: "  {{ .Title | faint }} {{ .Content | faint }} {{ if .DueDate }}({{ .DueDate | faint }}){{ end }} [{{ .Priority | faint }}]",
			Selected: "✓ {{ .Title | green }} {{ .Content | green | bold }} {{ if .DueDate }}({{ .DueDate | green }}){{ end }} [{{ .Priority | green }}]",
		}

		prompt := promptui.Select{
			Label:     i18n.GetMessage(config.GetConfig().Language, "select_task_to_remove"),
			Items:     unfinishedNotes,
			Templates: templates,
			Size:      10,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return fmt.Errorf(i18n.GetMessage(config.GetConfig().Language, "task_select_failed"), err)
		}

		selectedNote := unfinishedNotes[idx]

		// 删除选中的任务
		if deleteErr := storage.DeleteNote(selectedNote.ID); deleteErr != nil {
			return deleteErr
		}

		logger.Success(i18n.GetMessage(config.GetConfig().Language, "task_deleted", selectedNote.Content))

		// 显示更新后的任务列表
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
	rootCmd.AddCommand(removeCmd)
}
