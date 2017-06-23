package task

import (
	"regexp"
	"strconv"
)

var (
	// (0.000414s) UPDATE "vms" SET "instance_id" = 61, "agent_id" = '014d2291-5117-417
	dbStatement = regexp.MustCompile("^\\(([\\d\\.]+)s\\) (SELECT|UPDATE|COMMIT|BEGIN|INSERT|DELETE|UPDATE)")
)

type DBStatement struct {
	Type string
	Time float64 // todo duration
}

var _ Action = DBStatement{}

func NewDBStatement(str string) *DBStatement {
	if m := dbStatement.FindStringSubmatch(str); len(m) > 0 {
		d1, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			return nil // todo
		}

		return &DBStatement{Type: m[2], Time: d1}
	}

	return nil
}

func (r DBStatement) Relation() Relation       { return ExactRelation{"db"} }
func (r DBStatement) ShortDescription() string { return "[db]" }
