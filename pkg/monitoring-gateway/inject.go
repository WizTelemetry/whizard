package monitoringgateway

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/prometheus-community/prom-label-proxy/injectproxy"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

const (
	queryParam          = "query"
	matchersParam       = "match[]"
	targetMatchersParam = "match_target[]"
)

func enforceQueryValues(e *injectproxy.Enforcer, v url.Values) (values string, hasQuery bool, err error) {
	// If no values were given or no query is present,
	// e.g. because the query came in the POST body
	// but the URL query string was passed, then finish early.
	if v.Get(queryParam) == "" {
		return v.Encode(), false, nil
	}
	expr, err := parser.ParseExpr(v.Get(queryParam))
	if err != nil {
		queryParseError := newQueryParseError(err)
		return "", true, queryParseError
	}

	if err := e.EnforceNode(expr); err != nil {
		if _, ok := err.(injectproxy.IllegalLabelMatcherError); ok {
			return "", true, err
		}
		enforceLabelError := newEnforceLabelError(err)
		return "", true, enforceLabelError
	}

	v.Set(queryParam, expr.String())
	return v.Encode(), true, nil
}

type queryParseError struct {
	msg string
}

func (e queryParseError) Error() string {
	return e.msg
}

func newQueryParseError(err error) queryParseError {
	return queryParseError{msg: fmt.Sprintf("error parsing query string %q", err.Error())}
}

type enforceLabelError struct {
	msg string
}

func (e enforceLabelError) Error() string {
	return e.msg
}

func newEnforceLabelError(err error) enforceLabelError {
	return enforceLabelError{msg: fmt.Sprintf("error enforcing label %q", err.Error())}
}

func injectMatcher(q url.Values, matcher *labels.Matcher, matchersParam string) error {
	matchers := q[matchersParam]
	if len(matchers) == 0 {
		q.Set(matchersParam, matchersToString(matcher))
	} else {
		// Inject label to existing matchers.
		for i, m := range matchers {
			ms, err := parser.ParseMetricSelector(m)
			if err != nil {
				return err
			}
			matchers[i] = matchersToString(append(ms, matcher)...)
		}
		q[matchersParam] = matchers
	}
	return nil
}

func matchersToString(ms ...*labels.Matcher) string {
	var el []string
	for _, m := range ms {
		el = append(el, m.String())
	}
	return fmt.Sprintf("{%v}", strings.Join(el, ","))
}
