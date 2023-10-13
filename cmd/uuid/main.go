package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofrs/uuid/v5"
	"github.com/spf13/cobra"
)

var Version = "1.1.0"

var (
	count    int
	uuidType string
)

func main() {
	cmd := cobra.Command{
		Use:          "uuid [flags] [count]",
		Short:        "Generate UUIDs",
		ValidArgs:    []string{"count"},
		Version:      Version,
		RunE:         run,
		SilenceUsage: true,
	}

	cmd.Long = `Generates UUIDs.

Supported types:
 [1|v1]:            Version 1 (date-time and MAC address)
 [3|v3|md5]:        Version 3 (namespace name-based MD5)
 [4|v4|random]:     Version 4 (random)
 [5|v5|sha1|sha-1]: Version 5 (namespace name-based SHA-1)
 [6|v6]:            Version 6 (k-sortable and random, field-compatible with v1)
 [7|v7]:            Version 7 (k-sortable and random)
`

	cmd.Flags().IntVarP(&count, "count", "n", 1, "number of UUIDs to generate")
	cmd.Flags().StringVarP(&uuidType, "type", "t", "v4", "type of UUID")
	cmd.Flags().String("namespace", "", "namespace for UUID v3 and v5")
	cmd.Flags().String("name", "", "name for UUID v3 and v5")
	cmd.ParseFlags(os.Args[1:])

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if !cmd.Flag("count").Changed && len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("parse arg as count: %w", err)
		}
		count = n
	}

	// TODO: if we eventually support too many aliases, use a map[string]byte
	var uuidVersion byte
	switch uuidType {
	case "1", "v1":
		uuidVersion = uuid.V1
	case "3", "v3", "md5":
		uuidVersion = uuid.V3
	case "4", "v4", "random":
		uuidVersion = uuid.V4
	case "5", "v5", "sha1", "sha-1":
		uuidVersion = uuid.V5
	case "6", "v6":
		uuidVersion = uuid.V6
	case "7", "v7":
		uuidVersion = uuid.V7
	default:
		return fmt.Errorf("unknown UUID type: %q", uuidType)
	}

	var namespace uuid.UUID
	var name string
	if uuidVersion == uuid.V3 || uuidVersion == uuid.V5 {
		if !cmd.Flag("namespace").Changed {
			return fmt.Errorf("namespace is required for UUID v3 and v5")
		}
		ns, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return fmt.Errorf("get namespace flag: %w", err)
		}
		namespace, err = uuid.FromString(ns)
		if err != nil {
			return fmt.Errorf("parse namespace UUID: %w", err)
		}

		if !cmd.Flag("name").Changed {
			return fmt.Errorf("name is required for UUID v3 and v5")
		}
		name, err = cmd.Flags().GetString("name")
		if err != nil {
			return fmt.Errorf("get name flag: %w", err)
		}
	}

	for i := 0; i < count; i++ {
		var id uuid.UUID
		var err error
		switch uuidVersion {
		case uuid.V1:
			id, err = uuid.NewV1()
		case uuid.V3:
			id = uuid.NewV3(namespace, name)
		case uuid.V4:
			id, err = uuid.NewV4()
		case uuid.V5:
			id = uuid.NewV5(namespace, name)
		case uuid.V6:
			id, err = uuid.NewV6()
		case uuid.V7:
			id, err = uuid.NewV7()
		default:
			// the old "this should never happen"
			return fmt.Errorf("internal: unknown uuid type: %#v", uuidType)
		}

		if err != nil {
			return fmt.Errorf("new uuid: %v", err)
		}
		fmt.Println(id.String())
	}

	return nil
}
