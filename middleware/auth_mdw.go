package middleware

import (
	"admin-api/auth"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func CheckApiKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKey := os.Getenv("API_KEY")
		reqApiKey := c.Request().Header.Get("api-key")

		if apiKey != reqApiKey {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}

func AuthenticateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		incomingToken, err := extractIncomingToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		claim, err := auth.UnSign(incomingToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		accountID, _ := strconv.Atoi(claim.Subject)
		c.Set("AccountID", accountID)
		return next(c)
	}
}

func extractIncomingToken(c echo.Context) (string, error) {
	if val := c.Request().Header.Get("Authorization"); val != "" {
		if !strings.HasPrefix(strings.ToLower(val), "bearer ") {
			return "", fmt.Errorf(fmt.Sprintf("invalid or malformed %q header, expected 'Bearer JWT-token...'", val))
		}
		return strings.Split(val, " ")[1], nil
	}

	return "", fmt.Errorf("missing header or parameter Authorization")
}
