package utils

import "regexp"

var (
	pattern = `^(([a-zA-Z0-9.-]+(:[0-9]+)?/)?([a-z0-9]+(?:[._-][a-z0-9]+)*/)*[a-z0-9]+(?:[._-][a-z0-9]+)*)?(?::[a-zA-Z0-9._-]+)?(?:@[A-Za-z0-9:]+)?$`
	re      = regexp.MustCompile(pattern)
)

func ValidateImageName(image string) bool {
	// Docker image regex pattern

	// Return true if the image matches the pattern, otherwise false
	return re.MatchString(image)
}
