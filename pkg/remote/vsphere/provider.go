package vsphere

import (
	"fmt"
	"net/url"
	"os"

	"github.com/cloudskiff/driftctl/pkg/output"

	"github.com/cloudskiff/driftctl/pkg/remote/terraform"
	tf "github.com/cloudskiff/driftctl/pkg/terraform"
)

type VSphereTerraformProvider struct {
	*terraform.TerraformProvider
}

type vsphereConfig struct {
	User     string
	Password string
	Server   string
}

func NewVSphereTerraformProvider(version string, progress output.Progress, configDir string) (*VSphereTerraformProvider, error) {
	p := &VSphereTerraformProvider{}
	providerKey := "vsphere"
	installer, err := tf.NewProviderInstaller(tf.ProviderConfig{
		Key:       providerKey,
		Version:   version,
		ConfigDir: configDir,
	})
	if err != nil {
		return nil, err
	}
	tfProvider, err := terraform.NewTerraformProvider(installer, terraform.TerraformProviderConfig{
		Name:         providerKey,
		DefaultAlias: p.GetConfig().getDefaultServer(),
		GetProviderConfig: func(server string) interface{} {
			return vsphereConfig{
				Server: p.GetConfig().getDefaultServer(),
			}
		},
	}, progress)
	if err != nil {
		return nil, err
	}
	p.TerraformProvider = tfProvider
	return p, err
}

func (c vsphereConfig) getDefaultServer() string {
	return c.Server
}

func (c vsphereConfig) url() (*url.URL, error) {
	u, err := url.Parse("https://" + c.Server + "/sdk")
	if err != nil {
		return nil, fmt.Errorf("Error parse url: %s", err)
	}

	u.User = url.UserPassword(c.User, c.Password)

	return u, nil
}

func (p VSphereTerraformProvider) GetConfig() vsphereConfig {
	return vsphereConfig{
		User:     os.Getenv("VSPHERE_USER"),
		Password: os.Getenv("VSPHERE_PASSWORD"),
		Server:   os.Getenv("VSPHERE_SERVER"),
	}
}
