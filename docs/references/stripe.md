# Stripe References

Use official Stripe references for compatibility decisions.

## stripe-mock

Official Stripe mock HTTP server.

- https://github.com/stripe/stripe-mock

Use as a request-validation sanity oracle only. stripe-mock is OpenAPI-backed,
stateless, hardcoded, and does not support testing specific responses or errors.

## Stripe OpenAPI

Machine-readable endpoint, parameter, schema, fixture, and expansion reference.

- https://github.com/stripe/openapi

## API Errors

Official error envelope fields, error types, status-code classes, and
programmatic error codes.

- https://docs.stripe.com/api/errors
- https://docs.stripe.com/error-codes
- https://docs.stripe.com/api/idempotent_requests

## Testing Values

Sandbox payment success, decline, authentication, and failure references.

- https://docs.stripe.com/testing

## Stripe CLI

Useful for local webhook listening and triggering.

- https://docs.stripe.com/stripe-cli
- https://docs.stripe.com/stripe-cli/triggers
- https://github.com/stripe/stripe-cli/tree/master/pkg/fixtures

## Webhook Testing

Local webhook testing and delivery behavior.

- https://docs.stripe.com/webhooks
- https://docs.stripe.com/webhooks/signatures

## Test Clocks

Official test clocks for time-dependent Billing resources.

- https://docs.stripe.com/billing/testing/test-clocks

## Checkout

- https://docs.stripe.com/payments/checkout

## Billing Portal

- https://docs.stripe.com/customer-management
