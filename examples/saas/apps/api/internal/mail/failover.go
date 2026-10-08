package mail

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
)

// FailoverTransport tries its transports in order and stops at the first that
// delivers (MAIL_MAILER=failover, with MAIL_FAILOVER naming them). Each failure
// is logged, so a provider that is quietly down still shows up.
type FailoverTransport struct {
	Transports []Transport
}

// NewFailover returns a transport that tries each of transports in turn.
func NewFailover(transports ...Transport) *FailoverTransport {
	return &FailoverTransport{Transports: transports}
}

// Name implements Transport: "failover(smtp,log)".
func (t *FailoverTransport) Name() string {
	names := make([]string, 0, len(t.Transports))
	for _, tr := range t.Transports {
		names = append(names, tr.Name())
	}
	return "failover(" + strings.Join(names, ",") + ")"
}

// Send implements Transport.
func (t *FailoverTransport) Send(ctx context.Context, m *Message) error {
	if len(t.Transports) == 0 {
		return errors.New("failover: no transports to try")
	}
	var errs []error
	for i, tr := range t.Transports {
		err := tr.Send(ctx, m)
		if err == nil {
			if i > 0 {
				log.Printf("mail: %s delivered %q after %d failed", tr.Name(), m.Subject, i)
			}
			return nil
		}
		log.Printf("mail: %s could not send %q: %v", tr.Name(), m.Subject, err)
		errs = append(errs, fmt.Errorf("%s: %w", tr.Name(), err))
		if ctx.Err() != nil {
			break
		}
	}
	return fmt.Errorf("failover: every transport failed: %w", errors.Join(errs...))
}
