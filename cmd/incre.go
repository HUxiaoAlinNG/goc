package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"golang.org/x/tools/cover"
	"k8s.io/test-infra/gopherage/pkg/cov"
	"k8s.io/test-infra/gopherage/pkg/util"
)

var increCmd = &cobra.Command{
	Use:   "incre [files...]",
	Short: "incre multiple coherent Go coverage files into a single file.",
	Long: `merge will merge two Go coverage files into a single coverage file.
merge requires that the files are 'coherent', meaning that if they both contain references to the
same paths, then the contents of those source files were identical for the binary that generated
each file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		runIncre(args, outputIncreProfile)
	},
}

var outputIncreProfile string

func init() {
	increCmd.Flags().StringVarP(&outputIncreProfile, "output", "o", "incrementprofile.cov", "output file")

	rootCmd.AddCommand(increCmd)
}

func runIncre(args []string, output string) {

	if len(args) != 2 {
		log.Fatalln("Expected only two coverage file.")
		return
	}

	profiles := make([][]*cover.Profile, 0)
	for _, path := range args {
		profile, err := util.LoadProfile(path)
		if err != nil {
			log.Fatalf("failed to open %s: %v", path, err)
			return
		}
		profiles = append(profiles, profile)
	}

	diffProfiles, err := cov.DiffProfiles(profiles[0], profiles[1])
	if err != nil {
		log.Fatalf("failed to merge files: %v", err)
		return
	}

	err = util.DumpProfile(output, diffProfiles)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
