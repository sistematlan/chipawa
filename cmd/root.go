package cmd

import (
	"os"

	"github.com/sistematlan/mistah/internal/i18n"
	"github.com/sistematlan/mistah/internal/wizard"
	"github.com/spf13/cobra"
)

// Global flags shared by every command.
var (
	// flagLang overrides locale autodetection. Empty = autodetect.
	flagLang string
	// flagAdvanced opts the user into technical phrasing and verbose paths.
	// Default behaviour is the simple (human-friendly) mode.
	flagAdvanced bool
)

var rootCmd = &cobra.Command{
	Use:   "mistah",
	Short: "Limpia tu Mac como desarrollador",
	Long:  "mistah — analiza disco, apps, caches y proyectos. Libera espacio con confirmación.",
	// PersistentPreRunE applies the locale before any subcommand executes.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		applyLangFlag()
		return nil
	},
	// RunE only fires when no subcommand is given. We use it to launch the
	// wizard so a bare `mistah` becomes the friendly entry point. Power
	// users still have every subcommand available.
	RunE: func(cmd *cobra.Command, args []string) error {
		return wizard.Run(os.Stdin, cmd.OutOrStdout())
	},
	// Disable cobra's automatic suggestions for misspelled subcommands when
	// no args are passed; otherwise `mistah` would print "did you mean…?"
	// because we now have RunE defined.
	SilenceUsage: true,
}

// applyLangFlag wires --lang into the i18n package. Accepts "es" / "en" /
// raw locale strings ("es_MX.UTF-8"); anything else triggers autodetect.
func applyLangFlag() {
	switch flagLang {
	case "":
		// no override; let i18n auto-detect from $LANG/$LC_ALL.
	case "es":
		i18n.Set(i18n.LangES)
	case "en":
		i18n.Set(i18n.LangEN)
	default:
		// Treat as locale string. Reuse the package's detection logic by
		// temporarily setting LANG.
		_ = os.Setenv("LANG", flagLang)
		i18n.Set("") // force re-detect
	}
}

// SimpleMode reports whether output should use non-technical phrasing.
// Subcommands query this via cmd.SimpleMode() when rendering.
func SimpleMode() bool { return !flagAdvanced }

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagLang, "lang", "",
		"Idioma (es | en). Por defecto autodetecta desde $LANG.")
	rootCmd.PersistentFlags().BoolVar(&flagAdvanced, "advanced", false,
		"Muestra detalles técnicos y rutas completas (modo desarrollador).")
}
