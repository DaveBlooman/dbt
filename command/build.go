package command

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/DaveBlooman/dbt/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/DaveBlooman/dbt/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
)

func CmdBuild(c *cli.Context) {

	srpmCommand := "BD=$(mktemp -d /tmp/docker.XXXXXX) && cd $BD && mkdir BUILD RPMS SOURCES SPECS SRPMS; cp /mnt/SPECS/*.spec $BD/SPECS/specfile; cp -r /mnt/SOURCES/* $BD/SOURCES/; sudo chown -R $(id -u):$(id -g) $BD; rpmbuild --define \"_topdir $BD\" -bs $BD/SPECS/specfile; sudo cp $BD/SRPMS/*.rpm /mnt/SRPMS/"
	rpmCommand := "BD=$(mktemp -d /tmp/docker.XXXXXX) && cd $BD && mkdir BUILD RPMS SOURCES SPECS SRPMS; cp /mnt/SRPMS/* $BD/SRPMS/; sudo chown -R $(id -u):$(id -g) $BD; rpmbuild --define \"_topdir $BD\" --rebuild $BD/SRPMS/*; sudo cp $BD/RPMS/*/*.rpm /mnt/RPMS/;"

	srpmOpts := dockerOptions(srpmCommand)
	rpmOpts := dockerOptions(rpmCommand)

	dockerImage := `
                            ##        .
                      ## ## ##       ==
                   ## ## ## ##      ===
               /\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\___/ ===
          ~~~ {~~ ~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~
               \______ o          __/
                 \    \        __/
                  \____\______/
                  | rpmbuild |
               __ |  __   __ | _  __   _
              /  \| /  \ /   |/  / _\ |
              \__/| \__/ \__ |\_ \__  |
				`

	fmt.Println(dockerImage)

	err := dockerContainer(srpmOpts)
	if err != nil {
		log.Fatal(err)
	}
	err = dockerContainer(rpmOpts)
	if err != nil {
		log.Fatal(err)
	}

}

func dockerContainer(opts docker.CreateContainerOptions) error {
	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}
	path := os.Getenv("DOCKER_CERT_PATH")
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)
	client, err := docker.NewTLSClient(dockerHost, cert, key, ca)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating client!!", err)
	}
	container, err := client.CreateContainer(opts)
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to create container - Exec setup failed - %v", err))
		os.Exit(1)
	}

	if err := client.StartContainer(container.ID, opts.HostConfig); err != nil {
		fmt.Println("error")
	}

	err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: &outBuf,
		ErrorStream:  &errBuf,
		Logs:         true,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(outBuf.String())

	err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID: container.ID,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to rm container - %s", err))
	}
	return err

}

func dockerOptions(command string) docker.CreateContainerOptions {
	volumePath := "/Users/dblooman/go/src/github.com/DaveBlooman/dbt"
	opts := docker.CreateContainerOptions{
		HostConfig: &docker.HostConfig{
			Binds: []string{
				volumePath + "/RPMS:/mnt/RPMS",
				volumePath + "/SOURCES:/mnt/SOURCES",
				volumePath + "/SRPMS:/mnt/SRPMS",
				volumePath + "/SPECS:/mnt/SPECS",
			},
		},
		Config: &docker.Config{
			Image:        "dbt-build",
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          []string{"sh", "-c", command},
			User:         "dbt",
		},
	}
	return opts
}
