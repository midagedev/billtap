# Data Model

## Customer

- id
- email
- name
- metadata
- createdAt

## Product

- id
- name
- description
- active
- metadata

## Price

- id
- productId
- currency
- unitAmount
- recurringInterval
- recurringIntervalCount
- active
- metadata

## CheckoutSession

- id
- customerId
- mode
- lineItems
- successUrl
- cancelUrl
- status
- paymentStatus
- subscriptionId
- invoiceId
- paymentIntentId
- createdAt
- completedAt

## Subscription

- id
- customerId
- status
- items
- currentPeriodStart
- currentPeriodEnd
- cancelAtPeriodEnd
- canceledAt
- latestInvoiceId
- metadata

## Invoice

- id
- customerId
- subscriptionId
- status
- currency
- subtotal
- total
- amountDue
- amountPaid
- attemptCount
- nextPaymentAttempt
- paymentIntentId
- createdAt

## PaymentIntent

- id
- customerId
- invoiceId
- amount
- currency
- status
- failureCode
- failureMessage
- createdAt

## WebhookEndpoint

- id
- url
- secret
- enabledEvents
- active
- createdAt

## Event

- id
- type
- objectType
- objectId
- payload
- createdAt
- scenarioRunId

## DeliveryAttempt

- id
- eventId
- endpointId
- attemptNumber
- status
- requestHeaders
- requestBody
- responseStatus
- responseBody
- error
- scheduledAt
- deliveredAt

## ScenarioRun

- id
- name
- status
- startedAt
- finishedAt
- currentStep
- report

## AssertionResult

- id
- scenarioRunId
- stepId
- target
- status
- expected
- actual
- message

## TenantPaymentRoute

- id
- tenantId
- rail
- connectedAccountId
- checkoutEnabled
- adminManagedEnabled
- metadata

## Workspace

- id
- ownerAccountId
- name
- customerId
- activeSubscriptionId
- tenantPaymentRouteId
- createdAt

## WorkspaceMember

- id
- workspaceId
- accountId
- email
- role
- hasSignedUp
- createdAt

## EntitlementSnapshot

- id
- workspaceId
- subscriptionId
- planTier
- paymentCycle
- subscriptionStatus
- basicSeatCount
- additionalSeatCount
- usedSeatCount
- freeTrialUsed
- isTrialing
- currentPeriodStart
- currentPeriodEnd
- canCheckoutInApp
- canAdminManageSubscription
- exportDefaultLimit
- exportTotalLimit
- exportTotalRemaining
- exportManualLimit
- exportManualRemaining
- exportIsInfinite
- exportRenewalDate
- sourceEventId
- createdAt

## ExportSession

- id
- workspaceId
- accountId
- designCases
- exportTo
- extraExportPaymentId
- status
- createdAt

## ExtraExportPayment

- id
- workspaceId
- accountId
- customerId
- paymentIntentId
- exportSessionId
- exportRequestHash
- provisionAttemptId
- includedCount
- extraCount
- unitAmount
- amount
- currency
- status
- failureCode
- metadata
- createdAt

## PaymentHistoryItem

- id
- accountId
- workspaceId
- invoiceId
- paymentIntentId
- eventType
- invoiceType
- amount
- currency
- status
- hostedInvoiceUrl
- createdAt

## SupportAction

- id
- workspaceId
- accountId
- actorId
- actionType
- request
- result
- relatedEventId
- createdAt

## ProfileEvent

- id
- profile
- eventType
- sourceObjectType
- sourceObjectId
- workspaceId
- accountId
- evidence
- createdAt

## ObservabilityExpectation

- id
- scenarioRunId
- profile
- signalName
- expectedAttributes
- collector
- status
- actual
- message
- createdAt
