package cache

import (
	"context"
	"fmt"
	"gitlab.vecomentman.com/vote-your-face/service/user/utils"
)

// StartResetUserPassword sets the password activation key
// associated with the user id and the activation of duration.
// Returns nil if set was sucessfull otherwise an error
func (c *redisCache) StartResetUserPassword(
	ctx context.Context,
	passwordActivationKey string,
	userId string,
) error {
	resetEntryDuration := utils.FormatDuration(c.config.Ttl.Token.Account.Password)
	err := c.set(ctx, passwordActivationKey, userId, resetEntryDuration)

	if err != nil {
		c.log.Errorf("error setting user registration data: %s", err)
		return err
	}

	return nil
}

// UserInPasswordReset gets the user id that is in password reset process.
// Returns the user id if found otherwise an error
func (c *redisCache) UserInPasswordReset(
	ctx context.Context,
	resetPasswordKey string,
) (string, error) {
	entry, err := c.get(ctx, resetPasswordKey)

	if err != nil {
		c.log.Errorf("error reading password reset key: %s", err)
		return "", err
	}

	if !entry.Exists {
		return "", fmt.Errorf("no entry exists for the given reset key")
	}

	return entry.Val, nil
}
