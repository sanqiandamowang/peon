package gui

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.Flags().StringVarP(cmdConfig.ConfigDir, "dir", "d", "./", "打开文件目录下所有的符合格式的文件(config/config.json存在则读取配置),若无使用默认配置")
	rootCmd.AddCommand(completionCmd)
}

var rootCmd = &cobra.Command{
	Use:   "peon",
	Short: "可视化编辑cli",
	Long:  "可视化编辑相应文件的cli工具，配置文件在config/config.json中。",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 验证并确保目录存在
		if err := ensureDir(cmdConfig.ConfigDir); err != nil {
			return err
		}

		DisBase()
		return nil
	},
}
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "生成shell自动补全脚本",
	Long: `生成相应shell的自动补全脚本。 
要启用自动补全功能，请执行以下命令:
Bash:
$ source <(yourprogram completion bash)

Zsh:
$ source <(yourprogram completion zsh)
$ compdef _yourprogram yourprogram

Fish:
$ yourprogram completion fish | source

PowerShell:
$ yourprogram completion powershell | Out-String | Invoke-Expression
`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "bash":
			err = rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			err = rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			err = rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			err = fmt.Errorf("unsupported shell type: %s", args[0])
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating completion script: %v\n", err)
			os.Exit(1)
		}
	},
}

func ensureDir(dir *string) error {
	// 使用 os.Stat 检查目录是否存在
	info, err := os.Stat(*dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("目录不存在: %s", dir)
	}
	if err != nil {
		return err
	}
	// 检查路径是否是目录
	if !info.IsDir() {
		return fmt.Errorf("路径不是目录: %s", dir)
	}
	return nil
}

func Run() {
	if err := LoadConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 从命令行标志中获取参数并更新 cmdConfig
	dir, err := rootCmd.Flags().GetString("dir")
	if err != nil {
		fmt.Println("Error retrieving 'dir' flag:", err)
		os.Exit(1)
	}
	// 更新 cmdConfig 的 ConfigDir 字段
	if cmdConfig.ConfigDir == nil || *cmdConfig.ConfigDir == "" {
		*cmdConfig.ConfigDir = "./"
	}
	*cmdConfig.ConfigDir = dir
	cobra.CheckErr(rootCmd.Execute())
}
