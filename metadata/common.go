package metadata

var Verbose bool = false

var DefaultOutputPath struct {
	AWS   string
	GCP   string
	Azure string
}

func Initialise() {
	DefaultOutputPath.AWS = "./out/aws.json"
	DefaultOutputPath.GCP = "./out/gcp.json"
	DefaultOutputPath.Azure = "./out/azure.json"
}

type Encoder func(string) string

func PassthroughEncoder(endpoint string) string {
	return endpoint
}

func spinner() func(verbose bool) string {
	spinner := "-"
	counter := 0
	return func(verbose bool) string {
		if verbose {
			return "\n[" + spinner + "]"
		} else {
			mod := counter % 4
			switch mod {
			case 0:
				spinner = "-"
			case 1:
				spinner = "\\"
			case 2:
				spinner = "|"
			case 3:
				spinner = "/"
			}
			counter = counter + 1
			return "\r\033[K[" + spinner + "]"
		}
	}
}

var spin = spinner()
