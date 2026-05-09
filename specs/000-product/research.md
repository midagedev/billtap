# Research Notes

## Existing Tools

### stripe-mock

Official Stripe mock HTTP server. Useful reference for API shapes, but it is not a full stateful billing lab.

Reference:

- https://github.com/stripe/stripe-mock

### Stripe CLI

Useful for local webhook testing and event triggering.

### Stripe Test Clocks

Useful for testmode billing time travel. Billtap should learn from this but provide local deterministic scenario clock behavior.

### Webhook relay/debugging tools

Tools like Hookdeck and WebhookStash solve parts of webhook capture and replay. Billtap's differentiator is billing state plus app assertion, not generic webhook relay.

## Product Insight

The crowded space is "Stripe mock." The more valuable space is "SaaS billing integration lab."

Billtap should compete on:

- hosted sandbox UI
- stateful subscription lifecycle
- deterministic scenarios
- webhook reliability testing
- app entitlement assertions
- developer dashboard

