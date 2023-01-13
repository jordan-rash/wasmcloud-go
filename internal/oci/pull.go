package oci

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry/remote"
)

func PullOCIRef(ctx context.Context, targetRef, tag string, log logr.Logger) ([]byte, []byte, error) {
	repo, err := remote.NewRepository(targetRef)
	if err != nil {
		panic(err)
	}

	des, err := repo.Resolve(ctx, tag)
	if err != nil {
		panic(err)
	}

	layers, err := content.Successors(ctx, repo, des)
	if err != nil {
		panic(err)
	}

	var img []byte
	var metadata []byte

	for _, l := range layers {
		switch l.MediaType {
		case "application/vnd.wasmcloud.actor.archive.config":
			metadata, err = content.FetchAll(ctx, repo, l)
			if err != nil {
				panic(err)
			}
		case "application/vnd.module.wasm.content.layer.v1+wasm":
			log.V(8).Info("downloading %s:%s", targetRef, tag)
			img, err = content.FetchAll(ctx, repo, l)
			if err != nil {
				panic(err)
			}
			log.V(8).Info("download complete")
		}
	}

	if img != nil {
		return img, metadata, nil
	}

	return nil, nil, errors.New("did not find artifact")
}
