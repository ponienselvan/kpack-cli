package builder

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	expv1alpha1 "github.com/pivotal/kpack/pkg/apis/experimental/v1alpha1"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned"
	"github.com/spf13/cobra"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewApplyCommand(kpackClient versioned.Interface, defaultNamespace string) *cobra.Command {
	var (
		path string
	)

	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply a builder configuration",
		Long:    "Apply a builder configuration by filename.\nThe builder will be created if it does not yet exist.\nOnly YAML files are accepted.",
		Example: "tbctl builder apply -f ./builder.yaml\ncat ./builder.yaml | tbctl builder apply -f -",
		RunE: func(cmd *cobra.Command, args []string) error {
			builderConfig, err := getBuilderConfig(path)
			if err != nil {
				return err
			}

			if builderConfig.Namespace == "" {
				builderConfig.Namespace = defaultNamespace
			}

			_, err = kpackClient.ExperimentalV1alpha1().CustomBuilders(builderConfig.Namespace).Get(builderConfig.Name, metav1.GetOptions{})
			if err != nil && !k8serrors.IsNotFound(err) {
				return err
			} else if k8serrors.IsNotFound(err) {
				_, err = kpackClient.ExperimentalV1alpha1().CustomBuilders(builderConfig.Namespace).Create(builderConfig)
			} else {
				_, err = kpackClient.ExperimentalV1alpha1().CustomBuilders(builderConfig.Namespace).Update(builderConfig)
			}
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "\"%s\" applied\n", builderConfig.Name)
			return err
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringVarP(&path, "file", "f", "", "path to the builder configuration file")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func getBuilderConfig(path string) (*expv1alpha1.CustomBuilder, error) {
	var (
		file io.ReadCloser
		err  error
	)

	if path == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(path)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var builderConfig expv1alpha1.CustomBuilder
	err = yaml.Unmarshal(buf, &builderConfig)
	if err != nil {
		return nil, err
	}
	return &builderConfig, nil
}
