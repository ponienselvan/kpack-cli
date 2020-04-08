package image

import (
	"github.com/ghodss/yaml"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewGetCommand(kpackClient versioned.Interface, defaultNamespace string) *cobra.Command {
	var (
		namespace string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an image configuration",
		Long:  "Get an image configuration by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			image, err := kpackClient.BuildV1alpha1().Images(namespace).Get(args[0], metav1.GetOptions{})
			if err != nil {
				return errors.Errorf("image \"%s\" not found", args[0])
			}

			bytes, err := yaml.Marshal(image)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write(bytes)
			return err
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", defaultNamespace, "the namespace of the image")
	return cmd
}
