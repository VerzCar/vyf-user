package cache

import (
	"context"
	"fmt"
	"gitlab.vecomentman.com/vote-your-face/service/user/utils"
	"strconv"
)

func (c *redisCache) RegisterNewCompany(
	ctx context.Context,
	verificationKey string,
	companyId int64,
) error {
	err := c.set(
		ctx,
		verificationKey,
		companyId,
		utils.FormatDuration(c.config.Ttl.Token.Account.Verification),
	)

	if err != nil {
		c.log.Errorf("error setting company verification key data: %s", err)
		return err
	}

	return nil
}

// CompanyInRegistration gets the company that is in registration process.
// Returns the company id if found otherwise an error
func (c *redisCache) CompanyInRegistration(
	ctx context.Context,
	verificationKey string,
) (int64, error) {

	entry, err := c.get(ctx, verificationKey)

	if err != nil {
		c.log.Errorf("error reading company cached verification data: %s", err)
		return 0, err
	}

	if !entry.Exists {
		c.log.Infof("company cannot be verified anymore. New token required.")
		return 0, fmt.Errorf("company cannot be verified anymore. new token required")
	}

	// convert cached data to current data type
	companyId, err := strconv.ParseInt(entry.Val, 10, 64)

	if err != nil {
		c.log.Errorf("could not parse company id from cahed entry: %s", err)
		return 0, err
	}

	return companyId, nil
}

// CleanUpCompanyRegistration deletes all the cached company registration data
// Returns nil if deleting was sucessfull otherwise an error
func (c *redisCache) CleanUpCompanyRegistration(
	ctx context.Context,
	verificationKey string,
) error {
	err := c.redis.Del(ctx, verificationKey).Err()

	if err != nil {
		c.log.Errorf("error deleting cached company registration entries: %s", err)
		return err
	}

	return nil
}
