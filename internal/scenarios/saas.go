package scenarios

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type saasState struct {
	tenant        map[string]any
	workspaces    map[string]map[string]any
	webhookEvents map[string]map[string]any
	payments      []map[string]any
	logs          []map[string]any
	next          int
}

func newSaaSState() *saasState {
	return &saasState{
		workspaces:    map[string]map[string]any{},
		webhookEvents: map[string]map[string]any{},
	}
}

func (s *saasState) id(prefix string) string {
	s.next++
	return fmt.Sprintf("%s_%04d", prefix, s.next)
}

func (r *Runner) runSaaSStep(_ context.Context, scenario Scenario, step Step, params map[string]any, state *runState) (map[string]any, error) {
	if state.saas.tenant == nil {
		state.saas.tenant = saasTenant(scenario)
	}
	switch step.Action {
	case "saas.tenant.configure":
		state.saas.tenant = saasTenantFromParams(params, scenario)
		return saveStep(state, step.ID, map[string]any{"tenant": state.saas.tenant}), nil
	case "saas.workspace.create":
		return r.saasWorkspaceCreate(step, params, state), nil
	case "saas.workspace.activate", "saas.subscription.get_current":
		return r.saasWorkspaceSnapshot(step, params, state)
	case "saas.subscription.checkout_upgrade", "saas.subscription.confirm_upgrade":
		return r.saasSubscriptionUpgrade(step, params, state), nil
	case "saas.subscription.preview_upgrade":
		return r.saasSubscriptionPreview(step, params, state), nil
	case "saas.subscription.cancel":
		return r.saasSubscriptionCancel(step, params, state), nil
	case "saas.subscription.stop_pending_cancellation":
		return r.saasSubscriptionResume(step, params, state), nil
	case "saas.seat.estimate_purchase":
		return r.saasSeatEstimate(step, params, state), nil
	case "saas.seat.purchase":
		return r.saasSeatPurchase(step, params, state), nil
	case "saas.member.invite":
		return r.saasMemberInvite(step, params, state), nil
	case "saas.member.delete":
		return r.saasMemberDelete(step, params, state), nil
	case "saas.export.summary":
		return r.saasExportSummary(step, params, state), nil
	case "saas.export.usage":
		return r.saasExportUsage(step, params, state), nil
	case "saas.export.products":
		return r.saasExportProducts(step, params, state), nil
	case "saas.export_session.create":
		return r.saasExportSessionCreate(step, params, state), nil
	case "saas.export_session.files":
		return r.saasExportSessionFiles(step, params, state), nil
	case "saas.extra_export.preview":
		return r.saasExtraExportPreview(step, params, state), nil
	case "saas.extra_export.create_payment_intent":
		return r.saasExtraExportPayment(step, params, state), nil
	case "saas.extra_export.provide":
		return r.saasExtraExportProvide(step, params, state), nil
	case "saas.payment.customer_portal":
		return r.saasCustomerPortal(step, params, state), nil
	case "saas.payment.history":
		return r.saasPaymentHistory(step, params, state), nil
	case "saas.backoffice.start_subscription", "saas.backoffice.change_plan", "saas.backoffice.change_period":
		return r.saasBackofficePlan(step, params, state), nil
	case "saas.backoffice.change_seat":
		return r.saasSeatPurchase(step, params, state), nil
	case "saas.backoffice.refund":
		return r.saasRefund(step, params, state), nil
	case "saas.backoffice.update_export_limit":
		return r.saasManualExportLimit(step, params, state), nil
	case "saas.webhook.platform", "saas.webhook.connect":
		return r.saasWebhook(step, params, state), nil
	case "webhook.deliver_duplicate", "webhook.deliver_out_of_order":
		return r.saasWebhookOverride(step, params, state), nil
	case "saas.support.bundle":
		return r.saasSupportBundle(step, params, state), nil
	case "saas.observability.expect":
		return r.saasObservability(step, params, state), nil
	default:
		return nil, fmt.Errorf("%w: unsupported action %q", ErrInvalidConfig, step.Action)
	}
}

func (r *Runner) saasWorkspaceCreate(step Step, params map[string]any, state *runState) map[string]any {
	workspaceID := state.saas.id("ws")
	ownerEmail := stringDefault(params, "ownerEmail", "owner@example.test")
	basicSeats := int64ValueDefault(params, "basicSeatCount", defaultBasicSeats(stringDefault(params, "planTier", "FREE")))
	additionalSeats := int64ValueDefault(params, "additionalSeatCount", 0)
	usedSeats := int64(1)
	exportRemaining := int64ValueDefault(params, "exportRemaining", defaultExportLimit(stringDefault(params, "planTier", "FREE")))
	workspace := map[string]any{
		"id":            workspaceID,
		"name":          stringDefault(params, "workspaceName", "SaaS Workspace"),
		"ownerEmail":    ownerEmail,
		"tenant":        state.saas.tenant,
		"customerId":    state.saas.id("cus"),
		"planTier":      stringDefault(params, "planTier", "FREE"),
		"paymentCycle":  stringDefault(params, "paymentCycle", "MONTHLY"),
		"subscription":  map[string]any{"id": state.saas.id("dbsub"), "status": statusForPlan(stringDefault(params, "planTier", "FREE")), "planTier": stringDefault(params, "planTier", "FREE"), "paymentCycle": stringDefault(params, "paymentCycle", "MONTHLY")},
		"seats":         map[string]any{"basic": basicSeats, "additional": additionalSeats, "used": usedSeats, "pending": int64(0), "capacity": basicSeats + additionalSeats},
		"members":       []map[string]any{{"email": ownerEmail, "state": "SIGNED_UP", "owner": true}},
		"exportSummary": map[string]any{"included": exportRemaining, "extra": int64(0), "totalLimit": exportRemaining, "totalRemaining": exportRemaining, "manualLimit": int64(0), "manualRemaining": int64(0), "infinite": false, "renews": "next_period"},
		"payments":      []map[string]any{},
		"webhooks":      []map[string]any{},
		"logs":          []map[string]any{},
	}
	state.saas.workspaces[workspaceID] = workspace
	state.saas.logs = append(state.saas.logs, saasLog("workspace.created", workspaceID))
	return saveStep(state, step.ID, map[string]any{"workspace": workspace})
}

func (r *Runner) saasWorkspaceSnapshot(step Step, params map[string]any, state *runState) (map[string]any, error) {
	workspace, err := saasWorkspace(params, state)
	if err != nil {
		return nil, err
	}
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "subscription": workspace["subscription"]}), nil
}

func (r *Runner) saasSubscriptionPreview(step Step, params map[string]any, state *runState) map[string]any {
	preview := map[string]any{
		"id":              state.saas.id("preview"),
		"workspace":       firstString(params, "workspaceRef", "workspace"),
		"nextPlanTier":    stringDefault(params, "nextPlanTier", "STANDARD"),
		"paymentCycle":    stringDefault(params, "paymentCycle", "MONTHLY"),
		"amountDue":       int64(9900),
		"currency":        "usd",
		"membersToRemove": []string{},
	}
	return saveStep(state, step.ID, map[string]any{"preview": preview})
}

func (r *Runner) saasSubscriptionUpgrade(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	planTier := firstString(params, "nextPlanTier", "planTier")
	if planTier == "" {
		planTier = "STANDARD"
	}
	cycle := stringDefault(params, "paymentCycle", "MONTHLY")
	status := "COMPLETED"
	if strings.Contains(strings.ToLower(stringDefault(params, "outcome", "payment_succeeded")), "fail") {
		status = "FAILED"
	}
	sub := asMap(workspace["subscription"])
	sub["status"] = status
	sub["planTier"] = planTier
	sub["paymentCycle"] = cycle
	workspace["planTier"] = planTier
	workspace["paymentCycle"] = cycle
	workspace["subscription"] = sub
	workspace["exportSummary"] = exportSummaryWithLimit(asMap(workspace["exportSummary"]), defaultExportLimit(planTier))
	payment := map[string]any{"id": state.saas.id("pay"), "invoice": state.saas.id("in"), "paymentIntent": state.saas.id("pi"), "status": "invoice_paid", "planTier": planTier}
	appendWorkspaceList(workspace, "payments", payment)
	state.saas.payments = append(state.saas.payments, payment)
	state.saas.logs = append(state.saas.logs, saasLog("subscription.upgraded", stringValue(workspace, "id")))
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "subscription": sub, "payment": payment})
}

func (r *Runner) saasSubscriptionCancel(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	sub := asMap(workspace["subscription"])
	sub["status"] = "CANCELED"
	sub["cancelAtPeriodEnd"] = true
	workspace["subscription"] = sub
	state.saas.logs = append(state.saas.logs, saasLog("subscription.cancel", stringValue(workspace, "id")))
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "subscription": sub})
}

func (r *Runner) saasSubscriptionResume(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	sub := asMap(workspace["subscription"])
	sub["status"] = "COMPLETED"
	sub["cancelAtPeriodEnd"] = false
	workspace["subscription"] = sub
	state.saas.logs = append(state.saas.logs, saasLog("subscription.resume", stringValue(workspace, "id")))
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "subscription": sub})
}

func (r *Runner) saasSeatEstimate(step Step, params map[string]any, state *runState) map[string]any {
	count := int64ValueDefault(params, "count", 1)
	estimate := map[string]any{"id": state.saas.id("seat_est"), "count": count, "amountDue": count * 1500, "currency": "usd"}
	return saveStep(state, step.ID, map[string]any{"estimate": estimate})
}

func (r *Runner) saasSeatPurchase(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	count := int64ValueDefault(params, "count", int64ValueDefault(params, "quantity", 1))
	seats := asMap(workspace["seats"])
	seats["additional"] = int64ish(seats["additional"]) + count
	seats["capacity"] = int64ish(seats["basic"]) + int64ish(seats["additional"])
	workspace["seats"] = seats
	purchase := map[string]any{"id": state.saas.id("seat_pay"), "count": count, "status": "succeeded"}
	state.saas.logs = append(state.saas.logs, saasLog("billing.seat_update", stringValue(workspace, "id")))
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "seats": seats, "seatPurchase": purchase})
}

func (r *Runner) saasMemberInvite(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	emails := stringSlice(params["emails"])
	seats := asMap(workspace["seats"])
	members := memberSlice(workspace["members"])
	if int64(len(members)+len(emails)) > int64ish(seats["capacity"]) {
		result := map[string]any{"status": "seat_limit_exceeded", "blocked": true, "workspace": workspace}
		return saveStep(state, step.ID, map[string]any{"invite": result, "workspace": workspace})
	}
	for _, email := range emails {
		members = append(members, map[string]any{"email": email, "state": "INVITED", "owner": false})
	}
	seats["used"] = int64(len(members))
	seats["pending"] = int64(len(emails))
	workspace["members"] = members
	workspace["seats"] = seats
	state.saas.logs = append(state.saas.logs, saasLog("member.invite", stringValue(workspace, "id")))
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "members": members, "seats": seats})
}

func (r *Runner) saasMemberDelete(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	remove := map[string]bool{}
	for _, email := range stringSlice(params["emails"]) {
		remove[email] = true
	}
	var members []map[string]any
	for _, member := range memberSlice(workspace["members"]) {
		if !remove[fmt.Sprint(member["email"])] {
			members = append(members, member)
		}
	}
	seats := asMap(workspace["seats"])
	seats["used"] = int64(len(members))
	workspace["members"] = members
	workspace["seats"] = seats
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "members": members, "seats": seats})
}

func (r *Runner) saasExportSummary(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "exportSummary": workspace["exportSummary"]})
}

func (r *Runner) saasExportUsage(step Step, params map[string]any, state *runState) map[string]any {
	usage := map[string]any{"id": state.saas.id("usage"), "workspace": firstString(params, "workspaceRef", "workspace"), "designCases": params["designCases"], "used": int64(1)}
	return saveStep(state, step.ID, map[string]any{"usage": usage})
}

func (r *Runner) saasExportProducts(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	summary := asMap(workspace["exportSummary"])
	expect := stringValue(params, "expect")
	if int64ish(summary["totalRemaining"]) <= 0 || expect == "quota_exhausted" {
		result := map[string]any{"status": "blocked", "reason": "PLAN_POLICY_BLOCKED", "unavailableReason": "PLAN_POLICY_BLOCKED"}
		return saveStep(state, step.ID, map[string]any{"export": result, "workspace": workspace})
	}
	summary["totalRemaining"] = int64ish(summary["totalRemaining"]) - 1
	workspace["exportSummary"] = summary
	result := map[string]any{"status": "EXPORT_SUCCEEDED", "designCases": params["designCases"], "exportTo": stringDefault(params, "exportTo", "viewer")}
	return saveStep(state, step.ID, map[string]any{"export": result, "workspace": workspace, "exportSummary": summary})
}

func (r *Runner) saasExportSessionCreate(step Step, params map[string]any, state *runState) map[string]any {
	session := map[string]any{"id": state.saas.id("exps"), "workspace": firstString(params, "workspaceRef", "workspace"), "status": "CREATED", "files": []string{"design-case.zip"}}
	return saveStep(state, step.ID, map[string]any{"exportSession": session})
}

func (r *Runner) saasExportSessionFiles(step Step, params map[string]any, state *runState) map[string]any {
	files := map[string]any{"session": firstString(params, "exportSessionRef", "session"), "files": []string{"design-case.zip"}}
	return saveStep(state, step.ID, map[string]any{"files": files})
}

func (r *Runner) saasExtraExportPreview(step Step, params map[string]any, state *runState) map[string]any {
	preview := map[string]any{"id": state.saas.id("extra_prev"), "workspace": firstString(params, "workspaceRef", "workspace"), "count": len(stringSlice(params["designCases"])), "amountDue": int64(500)}
	return saveStep(state, step.ID, map[string]any{"preview": preview})
}

func (r *Runner) saasExtraExportPayment(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	status := "PAYMENT_SUCCEEDED"
	if strings.Contains(strings.ToLower(stringDefault(params, "outcome", "payment_succeeded")), "fail") {
		status = "PAYMENT_FAILED"
	}
	payment := map[string]any{"id": state.saas.id("extra_pay"), "paymentIntent": state.saas.id("pi"), "status": status, "workspaceId": stringValue(workspace, "id")}
	appendWorkspaceList(workspace, "payments", payment)
	state.saas.payments = append(state.saas.payments, payment)
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "extraExportPayment": payment})
}

func (r *Runner) saasExtraExportProvide(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	summary := asMap(workspace["exportSummary"])
	summary["extra"] = int64ish(summary["extra"]) + 1
	summary["totalLimit"] = int64ish(summary["totalLimit"]) + 1
	summary["totalRemaining"] = int64ish(summary["totalRemaining"]) + 1
	workspace["exportSummary"] = summary
	result := map[string]any{"id": firstString(params, "extraExportPaymentRef", "extraExportPayment"), "status": "EXPORT_SUCCEEDED", "workspaceId": stringValue(workspace, "id")}
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "extraExportPayment": result, "exportSummary": summary})
}

func (r *Runner) saasCustomerPortal(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	portal := map[string]any{"url": "https://billtap.local/portal/****", "masked": true, "workspaceId": stringValue(workspace, "id")}
	return saveStep(state, step.ID, map[string]any{"portal": portal})
}

func (r *Runner) saasPaymentHistory(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	history := []map[string]any{
		{"type": "invoice.paid", "status": "paid", "workspaceId": stringValue(workspace, "id")},
		{"type": "payment_intent.succeeded", "status": "succeeded", "workspaceId": stringValue(workspace, "id")},
		{"type": "charge.refunded", "status": "refunded", "workspaceId": stringValue(workspace, "id")},
	}
	return saveStep(state, step.ID, map[string]any{"paymentHistory": history})
}

func (r *Runner) saasBackofficePlan(step Step, params map[string]any, state *runState) map[string]any {
	return r.saasSubscriptionUpgrade(step, params, state)
}

func (r *Runner) saasRefund(step Step, params map[string]any, state *runState) map[string]any {
	refund := map[string]any{"id": state.saas.id("refund"), "status": "succeeded", "amount": int64ValueDefault(params, "amount", 0)}
	state.saas.payments = append(state.saas.payments, refund)
	return saveStep(state, step.ID, map[string]any{"refund": refund})
}

func (r *Runner) saasManualExportLimit(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	summary := asMap(workspace["exportSummary"])
	limit := int64ValueDefault(params, "manualLimit", int64ValueDefault(params, "limit", 0))
	summary["manualLimit"] = limit
	summary["manualRemaining"] = limit
	workspace["exportSummary"] = summary
	return saveStep(state, step.ID, map[string]any{"workspace": workspace, "exportSummary": summary})
}

func (r *Runner) saasWebhook(step Step, params map[string]any, state *runState) map[string]any {
	eventID := state.saas.id("evt")
	source := "platform"
	if step.Action == "saas.webhook.connect" {
		source = "connect"
	}
	status := "processed"
	if stringValue(params, "customerId") == "cus_unknown" || stringValue(params, "expect") == "identity_unresolved" {
		status = "identity_unresolved"
	}
	event := map[string]any{
		"id":                 eventID,
		"type":               stringDefault(params, "eventType", "invoice.paid"),
		"sourceType":         source,
		"connectedAccountId": stringValue(params, "connectedAccountId"),
		"status":             status,
		"workspaceId":        firstString(params, "workspaceRef", "workspace"),
		"duplicateIgnored":   false,
		"outOfOrder":         false,
	}
	state.saas.webhookEvents[eventID] = event
	if workspace, err := saasWorkspace(params, state); err == nil {
		appendWorkspaceList(workspace, "webhooks", event)
	}
	return saveStep(state, step.ID, map[string]any{"event": event, "webhook": event})
}

func (r *Runner) saasWebhookOverride(step Step, params map[string]any, state *runState) map[string]any {
	eventID := firstString(params, "eventRef", "event")
	event := state.saas.webhookEvents[eventID]
	if event == nil {
		event = map[string]any{"id": eventID, "status": "processed"}
	}
	if step.Action == "webhook.deliver_duplicate" {
		event["duplicateIgnored"] = true
		event["duplicate"] = true
	}
	if step.Action == "webhook.deliver_out_of_order" {
		event["outOfOrder"] = true
	}
	state.saas.webhookEvents[eventID] = event
	return saveStep(state, step.ID, map[string]any{"event": event, "webhook": event})
}

func (r *Runner) saasSupportBundle(step Step, params map[string]any, state *runState) map[string]any {
	workspace, _ := saasWorkspace(params, state)
	bundle := map[string]any{
		"id":            state.saas.id("bundle"),
		"workspace":     workspace,
		"subscription":  workspace["subscription"],
		"seats":         workspace["seats"],
		"exportSummary": workspace["exportSummary"],
		"payments":      workspace["payments"],
		"webhooks":      workspace["webhooks"],
		"workspaceLogs": state.saas.logs,
		"appAssertions": assertionSummaries(state),
		"generatedAt":   state.clock.Now().Format(time.RFC3339),
	}
	return saveStep(state, step.ID, map[string]any{"supportBundle": bundle})
}

func (r *Runner) saasObservability(step Step, params map[string]any, state *runState) map[string]any {
	signals := params["signals"]
	if signals == nil {
		signals = []string{"billing.subscription_update", "billing.seat_update", "billing.export_update", "billing.webhook_processed"}
	}
	expectation := map[string]any{"id": state.saas.id("obs"), "signals": signals, "status": "defined"}
	return saveStep(state, step.ID, map[string]any{"observability": expectation})
}

func saasTenant(s Scenario) map[string]any {
	rail := strings.ToUpper(s.SaaS.Tenant.Rail)
	if rail == "" {
		rail = "CARD"
	}
	return map[string]any{
		"id":                         firstNonEmpty(s.SaaS.Tenant.ID, "tenant_direct"),
		"rail":                       rail,
		"connectedAccountId":         s.SaaS.Tenant.ConnectedAccountID,
		"canCheckoutInApp":           rail == "CARD",
		"canAdminManageSubscription": rail == "OUT_OF_BAND" || rail == "CONNECT",
	}
}

func saasTenantFromParams(params map[string]any, scenario Scenario) map[string]any {
	tenant := saasTenant(scenario)
	for _, key := range []string{"id", "rail", "connectedAccountId"} {
		if value := stringValue(params, key); value != "" {
			tenant[key] = value
		}
	}
	return tenant
}

func saasWorkspace(params map[string]any, state *runState) (map[string]any, error) {
	workspaceID := firstString(params, "workspaceRef", "workspace", "workspaceId")
	if workspaceID == "" && len(state.saas.workspaces) == 1 {
		for id := range state.saas.workspaces {
			workspaceID = id
		}
	}
	workspace := state.saas.workspaces[workspaceID]
	if workspace == nil {
		return nil, fmt.Errorf("%w: workspace %q not found", ErrInvalidConfig, workspaceID)
	}
	return workspace, nil
}

func saveStep(state *runState, stepID string, output map[string]any) map[string]any {
	state.results[stepID] = output
	return output
}

func defaultBasicSeats(plan string) int64 {
	switch strings.ToUpper(plan) {
	case "PREMIUM":
		return 10
	case "STANDARD":
		return 1
	default:
		return 1
	}
}

func defaultExportLimit(plan string) int64 {
	switch strings.ToUpper(plan) {
	case "PREMIUM":
		return 5000
	case "STANDARD":
		return 1000
	default:
		return 100
	}
}

func statusForPlan(plan string) string {
	if strings.ToUpper(plan) == "FREE" {
		return "NONE"
	}
	return "COMPLETED"
}

func exportSummaryWithLimit(summary map[string]any, limit int64) map[string]any {
	summary["included"] = limit
	summary["totalLimit"] = limit + int64ish(summary["extra"]) + int64ish(summary["manualLimit"])
	if int64ish(summary["totalRemaining"]) <= 0 {
		summary["totalRemaining"] = limit
	}
	return summary
}

func asMap(value any) map[string]any {
	if m, ok := value.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func appendWorkspaceList(workspace map[string]any, key string, value map[string]any) {
	items, _ := workspace[key].([]map[string]any)
	items = append(items, value)
	workspace[key] = items
}

func memberSlice(value any) []map[string]any {
	if items, ok := value.([]map[string]any); ok {
		return items
	}
	return nil
}

func stringSlice(value any) []string {
	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			out = append(out, fmt.Sprint(item))
		}
		return out
	case string:
		if typed == "" {
			return nil
		}
		return []string{typed}
	default:
		return nil
	}
}

func int64ish(value any) int64 {
	switch typed := value.(type) {
	case int:
		return int64(typed)
	case int64:
		return typed
	case int32:
		return int64(typed)
	case float64:
		return int64(typed)
	case float32:
		return int64(typed)
	case uint:
		return int64(typed)
	case uint64:
		return int64(typed)
	default:
		return 0
	}
}

func saasLog(action, workspaceID string) map[string]any {
	return map[string]any{"action": action, "workspaceId": workspaceID}
}

func assertionSummaries(state *runState) []map[string]any {
	var out []map[string]any
	for stepID, result := range state.results {
		if assertion, ok := result["assertion"]; ok {
			out = append(out, map[string]any{"step": stepID, "assertion": assertion})
		}
	}
	return out
}
