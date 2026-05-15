// Package output provides buffered, format-aware writers for logslice output.
//
// It supports writing extracted log lines to any io.Writer destination
// (defaulting to os.Stdout) with optional line numbering.
//
// # Formats
//
// FormatRaw (default) writes each line as-is followed by a newline.
// FormatNumbered prefixes each line with a tab-separated line counter,
// useful for debugging or post-processing pipelines.
//
// # Usage
//
//	w := output.New(output.Options{
//		Dest:   os.Stdout,
//		Format: output.FormatNumbered,
//	})
//	defer w.Flush()
//
//	for _, line := range lines {
//		if err := w.WriteLine(line); err != nil {
//			log.Fatal(err)
//		}
//	}
package output
