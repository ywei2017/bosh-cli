package validation

import (
	"errors"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	bmerr "github.com/cloudfoundry/bosh-micro-cli/errors"
	bmrel "github.com/cloudfoundry/bosh-micro-cli/release"
)

type CpiValidator struct {
}

func NewCpiValidator() CpiValidator {
	return CpiValidator{}
}

func (v CpiValidator) Validate(release bmrel.Release) error {
	errs := v.validateCpiJob(release)
	if len(errs) > 0 {
		wrappedErrs := []error{}
		for _, err := range errs {
			wrappedErrs = append(wrappedErrs, bosherr.WrapError(err, "Validating CPI release"))
		}
		return bmerr.NewExplainableError(errs)
	}

	return nil
}

func (v CpiValidator) validateCpiJob(release bmrel.Release) []error {
	errs := []error{}

	job, ok := release.FindJobByName("cpi")
	if !ok {
		return append(errs, errors.New("Job `cpi' is missing from release"))
	}

	_, ok = job.FindTemplateByValue("bin/cpi")
	if !ok {
		errs = append(errs, errors.New("Job `cpi' is missing bin/cpi target"))
	}

	_, ok = job.FindTemplateByValue("bin/micro_discover_ip")
	if !ok {
		errs = append(errs, errors.New("Job `cpi' is missing bin/micro_discover_ip target"))
	}

	return errs
}