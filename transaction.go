package mongoose

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *Model[M]) Transaction(fnc func(session mongo.SessionContext) error, opts ...*options.TransactionOptions) error {
	session, err := m.connect.Client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(m.Ctx)
	_, err = session.WithTransaction(m.Ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		m.SetContext(sessionContext)
		return nil, fnc(sessionContext)
	}, opts...)

	return err
}
