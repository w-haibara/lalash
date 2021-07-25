package lalash

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func (cmd Command) setInternalStringFamily() {
	cmd.Internal.Cmds.Store("s-compare", InternalCmd{
		Usage: "s-compare",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.Compare(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-contains", InternalCmd{
		Usage: "s-contains",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.Contains(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-contains-any", InternalCmd{
		Usage: "s-contains-any",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.ContainsAny(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-contains-rune", InternalCmd{
		Usage: "s-contains-rune",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			i, err := strconv.Atoi(argv[1])
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.ContainsRune(argv[0], rune(i)))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-count", InternalCmd{
		Usage: "s-count",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.Count(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-equal-fold", InternalCmd{
		Usage: "s-equal-fold",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.EqualFold(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-fields", InternalCmd{
		Usage: "s-fields",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}

			for _, v := range strings.Fields(argv[0]) {
				fmt.Fprintln(cmd.Stdout, v)
			}

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-fields-func", InternalCmd{
		Usage: "s-has-prefix",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, "unimplemented")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-has-prefix", InternalCmd{
		Usage: "s-has-prefix",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.HasPrefix(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-has-suffix", InternalCmd{
		Usage: "s-has-suffix",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.HasSuffix(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-index", InternalCmd{
		Usage: "s-index",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.Index(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-index-any", InternalCmd{
		Usage: "s-index-any",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.IndexAny(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-index-byte", InternalCmd{
		Usage: "s-index-byte",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			if len(argv[1]) != 1 {
				return fmt.Errorf("length of substr must be 1 byte")
			}

			fmt.Fprintln(cmd.Stdout, strings.IndexByte(argv[0], byte(argv[1][0])))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-index-func", InternalCmd{
		Usage: "s-index-func",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			fmt.Fprintln(cmd.Stdout, "unimplemented")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-index-rune", InternalCmd{
		Usage: "s-index-rune",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			i, err := strconv.Atoi(argv[1])
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.IndexRune(argv[0], rune(i)))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-join", InternalCmd{
		Usage: "s-join",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			f := flag.NewFlagSet("join", flag.ContinueOnError)
			sep := f.String("sep", "", "")
			if err := f.Parse(argv); err != nil {
				return err
			}

			if err := checkArgv(f.Args(), 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.Join(f.Args(), *sep))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-last-index", InternalCmd{
		Usage: "s-last-index",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.LastIndex(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-last-index-any", InternalCmd{
		Usage: "s-last-index-any",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.LastIndexAny(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-last-index-byte", InternalCmd{
		Usage: "s-last-index-byte",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			if len(argv[1]) != 1 {
				return fmt.Errorf("length of substr must be 1 byte")
			}

			fmt.Fprintln(cmd.Stdout, strings.LastIndexByte(argv[0], byte(argv[1][0])))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-last-index-func", InternalCmd{
		Usage: "s-last-index-func",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			fmt.Fprintln(cmd.Stdout, "unimplemented")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-map", InternalCmd{
		Usage: "s-map",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			fmt.Fprintln(cmd.Stdout, "unimplemented")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-repeat", InternalCmd{
		Usage: "s-repeat",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			i, err := strconv.Atoi(argv[1])
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.Repeat(argv[0], i))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-replace", InternalCmd{
		Usage: "s-replace",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			f := flag.NewFlagSet("replace", flag.ContinueOnError)
			n := f.Int("n", -1, "")
			if err := f.Parse(argv); err != nil {
				return err
			}

			if err := checkArgv(f.Args(), 3); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.Replace(f.Arg(0), f.Arg(1), f.Arg(2), *n))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-replace-all", InternalCmd{
		Usage: "s-replace-all",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 3); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.ReplaceAll(argv[0], argv[1], argv[2]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-split", InternalCmd{
		Usage: "s-split",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			for _, v := range strings.Split(argv[0], argv[1]) {
				fmt.Fprintln(cmd.Stdout, v)
			}

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-split-after", InternalCmd{
		Usage: "s-split-after",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			for _, v := range strings.SplitAfter(argv[0], argv[1]) {
				fmt.Fprintln(cmd.Stdout, v)
			}

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-split-after-n", InternalCmd{
		Usage: "s-split-after-n",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			f := flag.NewFlagSet("split after n", flag.ContinueOnError)
			n := f.Int("n", -1, "")
			if err := f.Parse(argv); err != nil {
				return err
			}

			if err := checkArgv(f.Args(), 2); err != nil {
				return err
			}

			for _, v := range strings.SplitAfterN(f.Arg(0), f.Arg(1), *n) {
				fmt.Fprintln(cmd.Stdout, v)
			}

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-split-n", InternalCmd{
		Usage: "s-split-n",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			f := flag.NewFlagSet("split n", flag.ContinueOnError)
			n := f.Int("n", -1, "")
			if err := f.Parse(argv); err != nil {
				return err
			}

			if err := checkArgv(f.Args(), 2); err != nil {
				return err
			}

			for _, v := range strings.SplitN(f.Arg(0), f.Arg(1), *n) {
				fmt.Fprintln(cmd.Stdout, v)
			}

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-title", InternalCmd{
		Usage: "s-title",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.Title(argv[0]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-to-lower", InternalCmd{
		Usage: "s-to-lower",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.ToLower(argv[0]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-to-lower-spechial", InternalCmd{
		Usage: "s-to-lower-spechial",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			fmt.Fprintln(cmd.Stdout, "unimplemented")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-to-title", InternalCmd{
		Usage: "s-to-title",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.ToTitle(argv[0]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-to-title-spechial", InternalCmd{
		Usage: "s-to-title-spechial",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			fmt.Fprintln(cmd.Stdout, "unimplemented")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-to-upper", InternalCmd{
		Usage: "s-to-upper",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.ToUpper(argv[0]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-to-upper-spechial", InternalCmd{
		Usage: "s-to-upper-spechial",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			fmt.Fprintln(cmd.Stdout, "unimplemented")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-to-valid-utf8", InternalCmd{
		Usage: "s-to-valid-utf8",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.ToValidUTF8(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-trim", InternalCmd{
		Usage: "s-trim",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.Trim(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-trim-func", InternalCmd{
		Usage: "s-trim-func",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			fmt.Fprintln(cmd.Stdout, "unimplemented")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-trim-left", InternalCmd{
		Usage: "s-trim-left",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.TrimLeft(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-trim-left-func", InternalCmd{
		Usage: "s-trim-left-func",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			fmt.Fprintln(cmd.Stdout, "unimplemented")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-trim-prefix", InternalCmd{
		Usage: "s-trim-prefix",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.TrimPrefix(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-trim-right", InternalCmd{
		Usage: "s-trim-right",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.TrimRight(argv[0], argv[1]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-trim-right-func", InternalCmd{
		Usage: "s-trim-right-func",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			fmt.Fprintln(cmd.Stdout, "unimplemented")

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-trim-space", InternalCmd{
		Usage: "s-trim-space",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 1); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.TrimSpace(argv[0]))

			return nil
		},
	})

	cmd.Internal.Cmds.Store("s-trim-suffix", InternalCmd{
		Usage: "s-trim-suffix",
		Fn: func(ctx context.Context, cmd Command, args string, argv ...string) error {
			if err := checkArgv(argv, 2); err != nil {
				return err
			}

			fmt.Fprintln(cmd.Stdout, strings.TrimSuffix(argv[0], argv[1]))

			return nil
		},
	})
}
