package oci

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/tetratelabs/wabin/binary"
	"github.com/tetratelabs/wabin/wasm"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry/remote"
)

type Metadata struct {
	ConfigLayer   []byte
	CustomSection []*wasm.CustomSection
}

func PullOCIRef(ctx context.Context, targetRef, tag string, log logr.Logger) ([]byte, *Metadata, error) {
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
	md := new(Metadata)

	for _, l := range layers {
		switch l.MediaType {
		case "application/vnd.wasmcloud.actor.archive.config":
			md.ConfigLayer, err = content.FetchAll(ctx, repo, l)
			if err != nil {
				panic(err)
			}
		case "application/vnd.module.wasm.content.layer.v1+wasm":
			log.V(1).Info(fmt.Sprintf("downloading %s:%s", targetRef, tag))
			img, err = content.FetchAll(ctx, repo, l)
			if err != nil {
				panic(err)
			}
			log.V(1).Info("download complete")
		}
	}

	if mod, err := binary.DecodeModule(img, wasm.CoreFeaturesV2); err != nil {
		log.Error(err, "failed to decode wasm module")
	} else {
		md.CustomSection = mod.CustomSections
		return img, md, nil
	}

	return nil, nil, errors.New("did not find artifact")
}
