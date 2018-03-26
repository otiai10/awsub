package core

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/otiai10/daap"
	"github.com/otiai10/ternary"
	"golang.org/x/sync/errgroup"
)

// Fetch downloads resources from cloud storage services,
// localize those URLs to local path,
// and, additionally, ensure the output directories exist.
func (job *Job) Fetch() error {

	eg := new(errgroup.Group)

	for _, input := range job.Parameters.Inputs {
		i := input
		eg.Go(func() error { return job.fetch(i) })
	}

	for _, output := range job.Parameters.Outputs {
		o := output
		eg.Go(func() error { return job.ensure(o) })
	}

	for _, env := range job.Parameters.Envs {
		job.addContainerEnv(env)
	}

	return eg.Wait()
}

// fetch
func (job *Job) fetch(input *Input) error {
	log.Println(job.Identity.Name, "fetch", input.Name, input.URL)

	if err := input.Localize(AWSUB_CONTAINERROOT); err != nil {
		return err
	}

	fetch := &daap.Execution{
		Inline:  "/lifecycle/download.sh",
		Env:     input.EnvForFetch(),
		Inspect: true,
	}

	ctx := context.Background()
	stream, err := job.Container.Routine.Exec(ctx, fetch)
	if err != nil {
		return err
	}

	for payload := range stream {
		fmt.Printf("&%d> %s\n", payload.Type, payload.Text())
	}

	if fetch.ExitCode != 0 {
		return fmt.Errorf(
			"failed to download `%s` with status %d, please use --verbose option",
			input.URL, fetch.ExitCode,
		)
	}

	job.addContainerEnv(input.Env())

	return nil
}

// ensure the output directories exist on the workflow container.
func (job *Job) ensure(output *Output) error {
	// log.Println(job.Identity.Name, "ensure", output.URL)
	if err := output.Localize(AWSUB_CONTAINERROOT); err != nil {
		return err
	}

	dir := ternary.If(output.Recursive).String(output.LocalPath, filepath.Dir(output.LocalPath))

	ensure := &daap.Execution{
		Inline:  fmt.Sprintf("mkdir -p %s", dir),
		Inspect: true,
	}

	ctx := context.Background()
	stream, err := job.Container.Routine.Exec(ctx, ensure)
	if err != nil {
		return err
	}

	for payload := range stream {
		fmt.Printf("&%d> %s\n", payload.Type, payload.Text())
	}

	if ensure.ExitCode != 0 {
		return fmt.Errorf(
			"failed to download `%s` with status %d, please use --verbose option",
			output.URL, ensure.ExitCode,
		)
	}

	job.addContainerEnv(output.Env())
	return nil
}

func (job *Job) addContainerEnv(env Env) {
	job.Container.Envs = append(job.Container.Envs, env)
}