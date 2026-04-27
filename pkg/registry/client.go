package registry

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

type Client struct {
	RegistryURL string
}

func NewClient(registryURL string) *Client {
	return &Client{RegistryURL: registryURL}
}

type ImageMetadata struct {
	Digest    string `json:"digest"`
	MediaType string `json:"mediaType"`
	Platform  string `json:"platform"`
}

func (c *Client) ListTags(ctx context.Context, repository string) ([]string, error) {
	repoName, err := c.normalizeRepository(repository)
	if err != nil {
		return nil, err
	}

	tags, err := remote.List(repoName, remote.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (c *Client) InspectImage(ctx context.Context, imageRef string) (*ImageMetadata, error) {
	ref, err := name.ParseReference(imageRef, name.WithDefaultRegistry(c.RegistryURL))
	if err != nil {
		return nil, err
	}

	desc, err := remote.Get(ref, remote.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	platform, err := platformFromDescriptor(desc)
	if err != nil {
		return nil, err
	}

	return &ImageMetadata{
		Digest:    desc.Digest.String(),
		MediaType: string(desc.MediaType),
		Platform:  platform,
	}, nil
}

func (c *Client) normalizeRepository(repository string) (name.Repository, error) {
	if strings.Contains(repository, "/") && strings.HasPrefix(repository, c.RegistryURL+"/") {
		return name.NewRepository(repository)
	}

	return name.NewRepository(fmt.Sprintf("%s/%s", c.RegistryURL, repository))
}

func platformFromDescriptor(desc *remote.Descriptor) (string, error) {
	switch desc.MediaType {
	case types.OCIImageIndex, types.DockerManifestList:
		index, err := desc.ImageIndex()
		if err != nil {
			return "", err
		}

		manifest, err := index.IndexManifest()
		if err != nil {
			return "", err
		}

		if len(manifest.Manifests) == 0 || manifest.Manifests[0].Platform == nil {
			return "", nil
		}

		platform := manifest.Manifests[0].Platform
		return formatPlatform(platform.OS, platform.Architecture), nil
	default:
		image, err := desc.Image()
		if err != nil {
			return "", err
		}

		cfg, err := image.ConfigFile()
		if err != nil {
			return "", err
		}

		return formatPlatform(cfg.OS, cfg.Architecture), nil
	}
}

func formatPlatform(osName, arch string) string {
	if osName == "" || arch == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s", osName, arch)
}
