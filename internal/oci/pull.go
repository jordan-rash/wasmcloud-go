package oci

import (
	"context"
	"errors"
	"log"

	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry/remote"
)

func PullOCIRef(ctx context.Context, targetRef, tag string) ([]byte, error) {
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

	for _, l := range layers {
		if l.MediaType == "application/vnd.module.wasm.content.layer.v1+wasm" {
			log.Printf("downloading %s:%s", targetRef, tag)
			img, err := content.FetchAll(ctx, repo, l)
			if err != nil {
				panic(err)
			}
			log.Printf("download complete")
			return img, nil
		}
	}

	return nil, errors.New("did not find artifact")
}
