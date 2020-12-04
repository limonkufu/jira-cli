package query

import (
	"github.com/ankitpokhrel/jira-cli/pkg/jql"
)

type issueParams struct {
	latest     bool
	watching   bool
	resolution string
	issueType  string
	status     string
	priority   string
	reporter   string
	assignee   string
	labels     []string
	reverse    bool
}

// Issue is a query type for issue command.
type Issue struct {
	Project string
	Flags   FlagParser
	params  *issueParams
}

// NewIssue creates and initialize new issue type.
func NewIssue(project string, flags FlagParser) (*Issue, error) {
	issue := Issue{
		Project: project,
		Flags:   flags,
	}

	ip := issueParams{}

	err := initParams(&ip, flags)
	if err != nil {
		return nil, err
	}

	issue.params = &ip

	return &issue, nil
}

func initParams(params *issueParams, flags FlagParser) error {
	latest, err := flags.GetBool("history")
	if err != nil {
		return err
	}

	watching, err := flags.GetBool("watching")
	if err != nil {
		return err
	}

	resolution, err := flags.GetString("resolution")
	if err != nil {
		return err
	}

	issueType, err := flags.GetString("type")
	if err != nil {
		return err
	}

	status, err := flags.GetString("status")
	if err != nil {
		return err
	}

	priority, err := flags.GetString("priority")
	if err != nil {
		return err
	}

	reporter, err := flags.GetString("reporter")
	if err != nil {
		return err
	}

	assignee, err := flags.GetString("assignee")
	if err != nil {
		return err
	}

	labels, err := flags.GetStringArray("label")
	if err != nil {
		return err
	}

	reverse, err := flags.GetBool("reverse")
	if err != nil {
		return err
	}

	params.latest = latest
	params.watching = watching
	params.resolution = resolution
	params.issueType = issueType
	params.status = status
	params.priority = priority
	params.reporter = reporter
	params.assignee = assignee
	params.labels = labels
	params.reverse = reverse

	return nil
}

// Get returns constructed jql query.
func (i *Issue) Get() string {
	obf := "created"

	q := jql.NewJQL(i.Project)

	q.And(func() {
		if i.params.latest {
			q.History()
			obf = "lastViewed"
		}

		if i.params.watching {
			q.Watching()
		}

		q.FilterBy("type", i.params.issueType).
			FilterBy("resolution", i.params.resolution).
			FilterBy("status", i.params.status).
			FilterBy("priority", i.params.priority).
			FilterBy("reporter", i.params.reporter).
			FilterBy("assignee", i.params.assignee)

		if len(i.params.labels) > 0 {
			q.In("labels", i.params.labels...)
		}
	})

	if i.params.reverse {
		q.OrderBy(obf, jql.DirectionAscending)
	} else {
		q.OrderBy(obf, jql.DirectionDescending)
	}

	return q.String()
}